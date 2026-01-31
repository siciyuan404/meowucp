package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/service"
	"github.com/meowucp/internal/ucp/model"
)

type CheckoutHandler struct {
	services    *service.Services
	idGenerator func() string
	config      CheckoutHandlerConfig
}

type CheckoutHandlerConfig struct {
	Links           []model.Link
	ContinueURLBase string
}

func NewCheckoutHandler(services *service.Services) *CheckoutHandler {
	return NewCheckoutHandlerWithConfig(services, CheckoutHandlerConfig{})
}

func NewCheckoutHandlerWithConfig(services *service.Services, config CheckoutHandlerConfig) *CheckoutHandler {
	return &CheckoutHandler{
		services:    services,
		idGenerator: defaultCheckoutID,
		config:      config,
	}
}

func (h *CheckoutHandler) Create(c *gin.Context) {
	var req model.CheckoutCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	recoverableMessages := make([]model.Message, 0, 2)
	if req.Currency == "" {
		recoverableMessages = append(recoverableMessages, model.Message{
			Type:     "error",
			Code:     "missing_field",
			Content:  "Currency is required",
			Severity: "recoverable",
		})
	}
	if len(req.LineItems) == 0 {
		recoverableMessages = append(recoverableMessages, model.Message{
			Type:     "error",
			Code:     "missing_field",
			Content:  "Line items are required",
			Severity: "recoverable",
		})
	}

	lineItemsJSON, err := json.Marshal(req.LineItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	totals := computeTotals(req.LineItems)
	totalsJSON, err := json.Marshal(totals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	links := resolvedLinks(resolveBaseURL(c), h.config.Links)
	linksJSON, err := json.Marshal(links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	continueURL := buildContinueURL(resolveBaseURL(c), h.config.ContinueURLBase, "", h.idGenerator)
	paymentHandlers := loadPaymentHandlers(h.services)
	status, messages := resolveMessagesAndStatus(len(paymentHandlers) > 0, recoverableMessages, nil)
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	checkoutID := extractCheckoutID(continueURL, h.idGenerator)
	session := &domain.CheckoutSession{
		ID:          checkoutID,
		Status:      status,
		Currency:    req.Currency,
		LineItems:   string(lineItemsJSON),
		Totals:      string(totalsJSON),
		Links:       string(linksJSON),
		Messages:    string(messagesJSON),
		ContinueURL: continueURL,
	}

	if h.services == nil || h.services.Checkout == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service_unavailable"})
		return
	}

	if err := h.services.Checkout.Create(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create_failed"})
		return
	}

	response := model.CheckoutSession{
		ID:          checkoutID,
		LineItems:   req.LineItems,
		Status:      status,
		Currency:    req.Currency,
		Totals:      totals,
		Messages:    messages,
		Links:       links,
		ContinueURL: continueURL,
		Payment: model.Payment{
			Handlers: paymentHandlers,
		},
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CheckoutHandler) Get(c *gin.Context) {
	if h.services == nil || h.services.Checkout == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service_unavailable"})
		return
	}

	checkoutID := c.Param("id")
	if checkoutID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_id"})
		return
	}

	session, err := h.services.Checkout.GetByID(checkoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}

	var lineItems []model.LineItem
	if err := json.Unmarshal([]byte(session.LineItems), &lineItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decode_failed"})
		return
	}

	var totals []model.Total
	if err := json.Unmarshal([]byte(session.Totals), &totals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decode_failed"})
		return
	}

	var links []model.Link
	if session.Links != "" {
		if err := json.Unmarshal([]byte(session.Links), &links); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decode_failed"})
			return
		}
	}

	var messages []model.Message
	if session.Messages != "" {
		if err := json.Unmarshal([]byte(session.Messages), &messages); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decode_failed"})
			return
		}
	}

	response := model.CheckoutSession{
		ID:          session.ID,
		LineItems:   lineItems,
		Status:      session.Status,
		Currency:    session.Currency,
		Totals:      totals,
		Messages:    messages,
		Links:       links,
		ContinueURL: session.ContinueURL,
		Payment: model.Payment{
			Handlers: loadPaymentHandlers(h.services),
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *CheckoutHandler) Update(c *gin.Context) {
	if h.services == nil || h.services.Checkout == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service_unavailable"})
		return
	}

	checkoutID := c.Param("id")
	if checkoutID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_id"})
		return
	}

	var req model.CheckoutUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	recoverableMessages := make([]model.Message, 0, 2)
	if req.Currency == "" {
		recoverableMessages = append(recoverableMessages, model.Message{
			Type:     "error",
			Code:     "missing_field",
			Content:  "Currency is required",
			Severity: "recoverable",
		})
	}
	if len(req.LineItems) == 0 {
		recoverableMessages = append(recoverableMessages, model.Message{
			Type:     "error",
			Code:     "missing_field",
			Content:  "Line items are required",
			Severity: "recoverable",
		})
	}

	buyerInputMessages := make([]model.Message, 0, 1)
	if req.RequiresSignIn {
		buyerInputMessages = append(buyerInputMessages, model.Message{
			Type:     "error",
			Code:     "requires_sign_in",
			Content:  "Sign-in required",
			Severity: "requires_buyer_input",
		})
	}

	lineItemsJSON, err := json.Marshal(req.LineItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	totals := computeTotals(req.LineItems)
	totalsJSON, err := json.Marshal(totals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	links := resolvedLinks(resolveBaseURL(c), h.config.Links)
	linksJSON, err := json.Marshal(links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	continueURL := buildContinueURL(resolveBaseURL(c), h.config.ContinueURLBase, checkoutID, h.idGenerator)
	paymentHandlers := loadPaymentHandlers(h.services)
	status, messages := resolveMessagesAndStatus(len(paymentHandlers) > 0, recoverableMessages, buyerInputMessages)
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	session := &domain.CheckoutSession{
		ID:          checkoutID,
		Status:      status,
		Currency:    req.Currency,
		LineItems:   string(lineItemsJSON),
		Totals:      string(totalsJSON),
		Links:       string(linksJSON),
		Messages:    string(messagesJSON),
		ContinueURL: continueURL,
	}

	if err := h.services.Checkout.Update(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
		return
	}

	response := model.CheckoutSession{
		ID:          checkoutID,
		LineItems:   req.LineItems,
		Status:      status,
		Currency:    req.Currency,
		Totals:      totals,
		Messages:    messages,
		Links:       links,
		ContinueURL: continueURL,
		Payment: model.Payment{
			Handlers: paymentHandlers,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *CheckoutHandler) Complete(c *gin.Context) {
	if h.services == nil || h.services.Checkout == nil || h.services.UCPOrder == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service_unavailable"})
		return
	}

	checkoutID := c.Param("id")
	if checkoutID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_id"})
		return
	}

	var req model.CheckoutCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	if req.PaymentData.HandlerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_payment_data"})
		return
	}

	session, err := h.services.Checkout.GetByID(checkoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}

	var lineItems []model.LineItem
	if err := json.Unmarshal([]byte(session.LineItems), &lineItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decode_failed"})
		return
	}

	totals := computeTotals(lineItems)

	orderItems := make([]*domain.OrderItem, 0, len(lineItems))
	var subtotalMinor int64
	for _, item := range lineItems {
		subtotalMinor += item.Item.Price * int64(item.Quantity)
		unitPrice := float64(item.Item.Price) / 100
		orderItems = append(orderItems, &domain.OrderItem{
			ProductName: item.Item.Title,
			SKU:         item.Item.ID,
			Quantity:    item.Quantity,
			UnitPrice:   unitPrice,
			TotalPrice:  unitPrice * float64(item.Quantity),
		})
	}

	order := &domain.Order{
		OrderNo:       buildOrderNo(checkoutID),
		Status:        "paid",
		Subtotal:      float64(subtotalMinor) / 100,
		Total:         float64(subtotalMinor) / 100,
		Currency:      session.Currency,
		PaymentMethod: req.PaymentData.HandlerID,
		PaymentStatus: "paid",
	}

	paymentPayload, err := json.Marshal(req.PaymentData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode_failed"})
		return
	}

	payment := &domain.Payment{
		Amount:         float64(subtotalMinor) / 100,
		PaymentMethod:  req.PaymentData.HandlerID,
		Status:         "paid",
		PaymentPayload: string(paymentPayload),
	}

	createdOrder, _, err := h.services.UCPOrder.CreateFromCheckout(order, orderItems, payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "complete_failed"})
		return
	}

	session.Status = "completed"
	sessionTotals, err := json.Marshal(totals)
	if err == nil {
		session.Totals = string(sessionTotals)
	}
	_ = h.services.Checkout.Update(session)

	response := model.CheckoutSession{
		ID:          session.ID,
		LineItems:   lineItems,
		Status:      "completed",
		Currency:    session.Currency,
		Totals:      totals,
		Links:       resolvedLinks(resolveBaseURL(c), h.config.Links),
		ContinueURL: session.ContinueURL,
		Payment: model.Payment{
			Handlers: loadPaymentHandlers(h.services),
		},
		Order: &model.OrderRef{
			ID: strconv.FormatInt(createdOrder.ID, 10),
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *CheckoutHandler) Cancel(c *gin.Context) {
	if h.services == nil || h.services.Checkout == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service_unavailable"})
		return
	}

	checkoutID := c.Param("id")
	if checkoutID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_id"})
		return
	}

	session, err := h.services.Checkout.GetByID(checkoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}

	session.Status = "canceled"
	if err := h.services.Checkout.Update(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cancel_failed"})
		return
	}

	var lineItems []model.LineItem
	_ = json.Unmarshal([]byte(session.LineItems), &lineItems)
	var totals []model.Total
	_ = json.Unmarshal([]byte(session.Totals), &totals)
	var links []model.Link
	_ = json.Unmarshal([]byte(session.Links), &links)
	var messages []model.Message
	_ = json.Unmarshal([]byte(session.Messages), &messages)

	response := model.CheckoutSession{
		ID:          session.ID,
		LineItems:   lineItems,
		Status:      "canceled",
		Currency:    session.Currency,
		Totals:      totals,
		Messages:    messages,
		Links:       links,
		ContinueURL: session.ContinueURL,
		Payment: model.Payment{
			Handlers: loadPaymentHandlers(h.services),
		},
	}

	c.JSON(http.StatusOK, response)
}

func computeTotals(items []model.LineItem) []model.Total {
	var subtotal int64
	for _, item := range items {
		subtotal += item.Item.Price * int64(item.Quantity)
	}

	return []model.Total{
		{Type: "subtotal", Amount: subtotal},
		{Type: "total", Amount: subtotal},
	}
}

func defaultCheckoutID() string {
	return "chk_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func loadPaymentHandlers(services *service.Services) []model.PaymentHandler {
	if services == nil || services.Handler == nil {
		return []model.PaymentHandler{}
	}
	handlers, err := services.Handler.List()
	if err != nil {
		return []model.PaymentHandler{}
	}

	result := make([]model.PaymentHandler, 0, len(handlers))
	for _, handler := range handlers {
		instrumentSchemas := []string{}
		if handler.Name == nowPaymentsName {
			instrumentSchemas = []string{cardInstrumentSchema}
		}
		configValue := parseHandlerConfig(handler.Config)
		result = append(result, model.PaymentHandler{
			ID:                strconv.FormatInt(handler.ID, 10),
			Name:              handler.Name,
			Version:           handler.Version,
			Spec:              handler.Spec,
			ConfigSchema:      handler.ConfigSchema,
			InstrumentSchemas: instrumentSchemas,
			Config:            configValue,
		})
	}
	return result
}

func buildOrderNo(checkoutID string) string {
	return "ORD-" + checkoutID
}

func resolveMessagesAndStatus(hasHandlers bool, recoverable []model.Message, buyerInput []model.Message) (string, []model.Message) {
	buyerInputMessages := append([]model.Message{}, buyerInput...)
	if !hasHandlers {
		buyerInputMessages = append(buyerInputMessages, model.Message{
			Type:     "error",
			Code:     "payment_required",
			Content:  "Payment method required",
			Severity: "requires_buyer_input",
		})
		buyerInputMessages = append(buyerInputMessages, model.Message{
			Type:     "error",
			Code:     "payment_handlers_missing",
			Content:  "Payment handlers are missing",
			Severity: "requires_buyer_input",
		})
	}

	messages := make([]model.Message, 0, len(buyerInputMessages)+len(recoverable))
	messages = append(messages, buyerInputMessages...)
	messages = append(messages, recoverable...)

	status := "ready_for_complete"
	if len(buyerInputMessages) > 0 {
		status = "requires_escalation"
	} else if len(recoverable) > 0 {
		status = "incomplete"
	}

	return status, messages
}

func resolvedLinks(baseURL string, configured []model.Link) []model.Link {
	if len(configured) > 0 {
		return configured
	}

	return []model.Link{
		{Type: "privacy_policy", URL: baseURL + "/privacy"},
		{Type: "terms_of_service", URL: baseURL + "/terms"},
		{Type: "refund_policy", URL: baseURL + "/refund"},
	}
}

func buildContinueURL(baseURL, configuredBase, checkoutID string, fallback func() string) string {
	if checkoutID == "" {
		checkoutID = fallback()
	}
	continueBase := strings.TrimSpace(configuredBase)
	if continueBase == "" {
		continueBase = baseURL + "/checkout-sessions"
	}
	continueBase = strings.TrimRight(continueBase, "/")
	return continueBase + "/" + checkoutID
}

func extractCheckoutID(continueURL string, fallback func() string) string {
	if continueURL == "" {
		return fallback()
	}
	parts := strings.Split(continueURL, "/")
	if len(parts) == 0 {
		return fallback()
	}
	return parts[len(parts)-1]
}
