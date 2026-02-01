package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/meowucp/internal/domain"
)

type IdempotencyStore interface {
	Create(record *domain.IdempotencyKey) error
	FindByUserIDAndKey(userID int64, key string) (*domain.IdempotencyKey, error)
	Update(record *domain.IdempotencyKey) error
}

type idempotencySnapshot struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

type captureWriter struct {
	gin.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *captureWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *captureWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func Idempotency(store IdempotencyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("Idempotency-Key")
		if key == "" || store == nil {
			c.Next()
			return
		}

		userID := int64(0)
		if value, ok := c.Get("user_id"); ok {
			if id, ok := value.(int64); ok {
				userID = id
			}
		}

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		hash := hashRequestBody(bodyBytes)

		record, err := store.FindByUserIDAndKey(userID, key)
		if err == nil {
			if record.RequestHash != hash {
				c.JSON(http.StatusConflict, gin.H{"error": "idempotency_key_reuse"})
				c.Abort()
				return
			}
			if record.ResponseSnapshot != nil && record.Status == "completed" {
				var snapshot idempotencySnapshot
				if unmarshalErr := json.Unmarshal([]byte(*record.ResponseSnapshot), &snapshot); unmarshalErr == nil {
					c.Data(snapshot.Status, "application/json", []byte(snapshot.Body))
					c.Abort()
					return
				}
			}
		} else if !gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "idempotency_store_failed"})
			c.Abort()
			return
		}

		if record == nil {
			record = &domain.IdempotencyKey{
				UserID:      userID,
				Key:         key,
				RequestHash: hash,
				Status:      "pending",
			}
			if err := store.Create(record); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "idempotency_store_failed"})
				c.Abort()
				return
			}
		}

		writer := &captureWriter{ResponseWriter: c.Writer}
		c.Writer = writer
		c.Next()

		snapshot := idempotencySnapshot{Status: writer.status, Body: writer.body.String()}
		payload, err := json.Marshal(snapshot)
		if err != nil {
			return
		}
		snapshotText := string(payload)
		record.ResponseSnapshot = &snapshotText
		record.Status = "completed"
		_ = store.Update(record)
	}
}

func hashRequestBody(body []byte) string {
	digest := sha256.Sum256(body)
	return hex.EncodeToString(digest[:])
}
