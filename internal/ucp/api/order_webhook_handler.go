package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/service"
	"github.com/meowucp/internal/ucp/model"
)

type OrderWebhookHandler struct {
	services *service.Services
	verifier SignatureVerifier
}

const (
	webhookSignatureHeader = "UCP-Signature"
	webhookKeyIDHeader     = "UCP-Key-Id"
)

type SignatureVerifier interface {
	Verify(r *http.Request, body []byte) error
}

func NewOrderWebhookHandler(services *service.Services) *OrderWebhookHandler {
	return &OrderWebhookHandler{services: services}
}

func NewOrderWebhookHandlerWithVerifier(services *service.Services, verifier SignatureVerifier) *OrderWebhookHandler {
	return &OrderWebhookHandler{services: services, verifier: verifier}
}

func (h *OrderWebhookHandler) Receive(c *gin.Context) {
	if h.services == nil || h.services.Webhook == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service_unavailable"})
		return
	}
	if h.services.WebhookReplay == nil || h.services.WebhookQueue == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "queue_or_replay_unconfigured"})
		return
	}
	if h.verifier == nil {
		_ = h.auditSignatureFailure(c, nil, "verifier_unconfigured")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "signature_verifier_unconfigured"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}

	if err := h.verifier.Verify(c.Request, body); err != nil {
		auditEventID := extractEventID(body)
		_ = h.auditSignatureFailure(c, body, "invalid_signature", auditEventID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_signature"})
		return
	}

	var payload model.OrderWebhookEvent
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	if payload.EventID == "" || payload.Order.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_required_fields"})
		return
	}

	if existing, err := h.services.Webhook.GetByEventID(payload.EventID); err == nil && existing != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}

	payloadHash := hashWebhookBody(body)
	if seen, err := h.services.WebhookReplay.Seen(payloadHash); err == nil && seen {
		_ = h.auditSignatureFailure(c, body, "replay_detected", payload.EventID)
		c.JSON(http.StatusConflict, gin.H{"error": "replay_detected"})
		return
	}
	_ = h.services.WebhookReplay.Mark(payloadHash, 600)

	event := &domain.UCPWebhookEvent{
		EventID:     payload.EventID,
		EventType:   payload.EventType,
		OrderID:     payload.Order.ID,
		Status:      "processed",
		PayloadHash: payloadHash,
		ReceivedAt:  time.Now(),
		ProcessedAt: timePtr(time.Now()),
	}

	if err := h.services.Webhook.Create(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "persist_failed"})
		return
	}

	if err := h.services.WebhookQueue.Enqueue(payload.EventID, string(body)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "enqueue_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *OrderWebhookHandler) auditSignatureFailure(c *gin.Context, body []byte, reason string, eventID ...string) error {
	if h.services == nil || h.services.WebhookAudit == nil {
		return nil
	}

	payloadHash := ""
	if len(body) > 0 {
		payloadHash = hashWebhookBody(body)
	}

	storedEventID := ""
	if len(eventID) > 0 {
		storedEventID = eventID[0]
	}

	audit := &domain.UCPWebhookAudit{
		EventID:         storedEventID,
		Reason:          reason,
		SignatureHeader: c.GetHeader(webhookSignatureHeader),
		KeyID:           c.GetHeader(webhookKeyIDHeader),
		PayloadHash:     payloadHash,
		CreatedAt:       time.Now(),
	}
	return h.services.WebhookAudit.Create(audit)
}

func extractEventID(body []byte) string {
	var payload model.OrderWebhookEvent
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	return payload.EventID
}

func hashWebhookBody(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func timePtr(value time.Time) *time.Time {
	return &value
}
