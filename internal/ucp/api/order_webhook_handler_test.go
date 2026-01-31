package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/service"
	"github.com/meowucp/internal/ucp/model"
)

type fakeWebhookRepo struct {
	items       map[string]*domain.UCPWebhookEvent
	createCount int
}

func newFakeWebhookRepo() *fakeWebhookRepo {
	return &fakeWebhookRepo{items: map[string]*domain.UCPWebhookEvent{}}
}

type fakeWebhookAuditRepo struct {
	items       []*domain.UCPWebhookAudit
	createCount int
}

func newFakeWebhookAuditRepo() *fakeWebhookAuditRepo {
	return &fakeWebhookAuditRepo{}
}

func (f *fakeWebhookAuditRepo) Create(audit *domain.UCPWebhookAudit) error {
	f.items = append(f.items, audit)
	f.createCount++
	return nil
}

func (f *fakeWebhookAuditRepo) List(offset, limit int) ([]*domain.UCPWebhookAudit, error) {
	if offset >= len(f.items) {
		return []*domain.UCPWebhookAudit{}, nil
	}
	end := offset + limit
	if end > len(f.items) {
		end = len(f.items)
	}
	return f.items[offset:end], nil
}

func (f *fakeWebhookAuditRepo) Count() (int64, error) {
	return int64(len(f.items)), nil
}

type fakeReplayStore struct {
	seen        map[string]bool
	createCount int
}

func newFakeReplayStore() *fakeReplayStore {
	return &fakeReplayStore{seen: map[string]bool{}}
}

func (f *fakeReplayStore) Seen(hash string) (bool, error) {
	return f.seen[hash], nil
}

func (f *fakeReplayStore) Mark(hash string, _ int) error {
	f.seen[hash] = true
	f.createCount++
	return nil
}

func (f *fakeReplayStore) FindByHash(hash string) (*domain.UCPWebhookReplay, error) {
	if f.seen[hash] {
		return &domain.UCPWebhookReplay{PayloadHash: hash}, nil
	}
	return nil, errors.New("not found")
}

func (f *fakeReplayStore) Create(replay *domain.UCPWebhookReplay) error {
	f.seen[replay.PayloadHash] = true
	f.createCount++
	return nil
}

type fakeWebhookQueue struct {
	items       []string
	createCount int
}

func newFakeWebhookQueue() *fakeWebhookQueue {
	return &fakeWebhookQueue{}
}

func (f *fakeWebhookQueue) Enqueue(eventID string, payload string) error {
	f.items = append(f.items, eventID+"|"+payload)
	f.createCount++
	return nil
}

func (f *fakeWebhookQueue) Create(job *domain.UCPWebhookJob) error {
	f.items = append(f.items, job.EventID+"|"+job.Payload)
	f.createCount++
	return nil
}

func (f *fakeWebhookQueue) ListDue(_ int) ([]*domain.UCPWebhookJob, error) {
	return []*domain.UCPWebhookJob{}, nil
}

func (f *fakeWebhookQueue) Update(_ *domain.UCPWebhookJob) error {
	return nil
}

func (f *fakeWebhookQueue) List(_ int, _ int) ([]*domain.UCPWebhookJob, error) {
	return []*domain.UCPWebhookJob{}, nil
}

func (f *fakeWebhookQueue) Count() (int64, error) {
	return int64(len(f.items)), nil
}

func (f *fakeWebhookQueue) FindByID(_ int64) (*domain.UCPWebhookJob, error) {
	return nil, errors.New("not implemented")
}

type fakeSignatureVerifier struct {
	err error
}

func (f fakeSignatureVerifier) Verify(_ *http.Request, _ []byte) error {
	return f.err
}

func (f *fakeWebhookRepo) Create(event *domain.UCPWebhookEvent) error {
	f.items[event.EventID] = event
	f.createCount++
	return nil
}

func (f *fakeWebhookRepo) FindByEventID(eventID string) (*domain.UCPWebhookEvent, error) {
	event, ok := f.items[eventID]
	if !ok {
		return nil, errors.New("not found")
	}
	return event, nil
}

func (f *fakeWebhookRepo) UpdateStatus(eventID string, status string) error {
	if event, ok := f.items[eventID]; ok {
		event.Status = status
		return nil
	}
	return errors.New("not found")
}

func (f *fakeWebhookRepo) MarkProcessed(eventID string) error {
	if event, ok := f.items[eventID]; ok {
		event.Status = "processed"
		return nil
	}
	return errors.New("not found")
}

func TestOrderWebhookIdempotent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeWebhookRepo()
	auditRepo := newFakeWebhookAuditRepo()
	replayStore := newFakeReplayStore()
	queue := newFakeWebhookQueue()
	services := &service.Services{
		Webhook:       service.NewWebhookEventService(repo),
		WebhookAudit:  service.NewWebhookAuditService(auditRepo),
		WebhookReplay: service.NewWebhookReplayService(replayStore),
		WebhookQueue:  service.NewWebhookQueueService(queue),
	}
	handler := NewOrderWebhookHandlerWithVerifier(services, fakeSignatureVerifier{})

	r := gin.New()
	r.POST("/ucp/v1/order-webhooks", handler.Receive)

	payload := model.OrderWebhookEvent{
		EventID:   "evt_1",
		EventType: "order.paid",
		Timestamp: "2026-01-29T10:45:32Z",
		Order: model.OrderWebhookOrder{
			ID:     "100001",
			Status: "paid",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	firstReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", bytes.NewReader(body))
	firstReq.Header.Set("Content-Type", "application/json")
	firstResp := httptest.NewRecorder()
	r.ServeHTTP(firstResp, firstReq)

	if firstResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", firstResp.Code)
	}

	secondReq := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", bytes.NewReader(body))
	secondReq.Header.Set("Content-Type", "application/json")
	secondResp := httptest.NewRecorder()
	r.ServeHTTP(secondResp, secondReq)

	if secondResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", secondResp.Code)
	}

	if repo.createCount != 1 {
		t.Fatalf("expected create count 1, got %d", repo.createCount)
	}

	event := repo.items["evt_1"]
	if event == nil {
		t.Fatalf("expected event to be stored")
	}
	if event.PayloadHash != hashBody(body) {
		t.Fatalf("expected payload hash to match")
	}
}

func TestOrderWebhookRejectsInvalidSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeWebhookRepo()
	auditRepo := newFakeWebhookAuditRepo()
	replayStore := newFakeReplayStore()
	queue := newFakeWebhookQueue()
	services := &service.Services{
		Webhook:       service.NewWebhookEventService(repo),
		WebhookAudit:  service.NewWebhookAuditService(auditRepo),
		WebhookReplay: service.NewWebhookReplayService(replayStore),
		WebhookQueue:  service.NewWebhookQueueService(queue),
	}
	handler := NewOrderWebhookHandlerWithVerifier(services, fakeSignatureVerifier{err: errors.New("bad sig")})

	r := gin.New()
	r.POST("/ucp/v1/order-webhooks", handler.Receive)

	payload := model.OrderWebhookEvent{
		EventID:   "evt_1",
		EventType: "order.paid",
		Timestamp: "2026-01-29T10:45:32Z",
		Order: model.OrderWebhookOrder{
			ID:     "100001",
			Status: "paid",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.Code)
	}
	if repo.createCount != 0 {
		t.Fatalf("expected no events to be stored")
	}
	if auditRepo.createCount != 1 {
		t.Fatalf("expected audit log to be created")
	}
	if auditRepo.items[0].Reason == "" {
		t.Fatalf("expected audit reason to be set")
	}
	if auditRepo.items[0].EventID != "evt_1" {
		t.Fatalf("expected event_id to be captured")
	}
}

func TestOrderWebhookRejectsReplay(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeWebhookRepo()
	auditRepo := newFakeWebhookAuditRepo()
	replayStore := newFakeReplayStore()
	queue := newFakeWebhookQueue()
	services := &service.Services{
		Webhook:       service.NewWebhookEventService(repo),
		WebhookAudit:  service.NewWebhookAuditService(auditRepo),
		WebhookReplay: service.NewWebhookReplayService(replayStore),
		WebhookQueue:  service.NewWebhookQueueService(queue),
	}
	handler := NewOrderWebhookHandlerWithVerifier(services, fakeSignatureVerifier{})

	r := gin.New()
	r.POST("/ucp/v1/order-webhooks", handler.Receive)

	payload := model.OrderWebhookEvent{
		EventID:   "evt_1",
		EventType: "order.paid",
		Timestamp: "2026-01-29T10:45:32Z",
		Order: model.OrderWebhookOrder{
			ID:     "100001",
			Status: "paid",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	hash := hashBody(body)
	replayStore.seen[hash] = true

	req := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", resp.Code)
	}
	if auditRepo.createCount != 1 {
		t.Fatalf("expected audit log to be created")
	}
	if queue.createCount != 0 {
		t.Fatalf("expected no queue enqueue")
	}
}

func TestOrderWebhookEnqueuesJob(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeWebhookRepo()
	auditRepo := newFakeWebhookAuditRepo()
	replayStore := newFakeReplayStore()
	queue := newFakeWebhookQueue()
	services := &service.Services{
		Webhook:       service.NewWebhookEventService(repo),
		WebhookAudit:  service.NewWebhookAuditService(auditRepo),
		WebhookReplay: service.NewWebhookReplayService(replayStore),
		WebhookQueue:  service.NewWebhookQueueService(queue),
	}
	handler := NewOrderWebhookHandlerWithVerifier(services, fakeSignatureVerifier{})

	r := gin.New()
	r.POST("/ucp/v1/order-webhooks", handler.Receive)

	payload := model.OrderWebhookEvent{
		EventID:   "evt_1",
		EventType: "order.paid",
		Timestamp: "2026-01-29T10:45:32Z",
		Order: model.OrderWebhookOrder{
			ID:     "100001",
			Status: "paid",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/ucp/v1/order-webhooks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}
	if queue.createCount != 1 {
		t.Fatalf("expected enqueue to be called")
	}
}

func hashBody(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}
