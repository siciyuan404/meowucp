package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/meowucp/internal/domain"
	"github.com/meowucp/internal/service"
	"github.com/meowucp/internal/ucp/model"
)

const (
	ucpVersion           = "2026-01-11"
	ucpServiceName       = "dev.ucp.shopping"
	checkoutName         = "dev.ucp.shopping.checkout"
	orderName            = "dev.ucp.shopping.order"
	ucpSpecOverview      = "https://ucp.dev/specification/overview"
	checkoutSpec         = "https://ucp.dev/specification/checkout"
	orderSpec            = "https://ucp.dev/specification/order"
	shoppingRESTSpec     = "https://ucp.dev/services/shopping/rest.openapi.json"
	checkoutSchema       = "https://ucp.dev/schemas/shopping/checkout.json"
	orderSchema          = "https://ucp.dev/schemas/shopping/order.json"
	cardInstrumentSchema = "https://ucp.dev/schemas/shopping/types/card_payment_instrument.json"
	nowPaymentsName      = "com.nowpayments"
)

type ProfileHandler struct {
	services *service.Services
}

func NewProfileHandler(services *service.Services) *ProfileHandler {
	return &ProfileHandler{services: services}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	baseURL := resolveBaseURL(c)

	profile := model.Profile{
		UCP: model.UCPProfile{
			Version: ucpVersion,
			Services: map[string]model.ServiceDefinition{
				ucpServiceName: {
					Version: ucpVersion,
					Spec:    ucpSpecOverview,
					REST: &model.RESTBinding{
						Schema:   shoppingRESTSpec,
						Endpoint: baseURL + "/ucp/v1",
					},
				},
			},
			Capabilities: []model.Capability{
				{
					Name:    checkoutName,
					Version: ucpVersion,
					Spec:    checkoutSpec,
					Schema:  checkoutSchema,
				},
				{
					Name:    orderName,
					Version: ucpVersion,
					Spec:    orderSpec,
					Schema:  orderSchema,
				},
			},
		},
	}

	profile.Payment = &model.PaymentConfig{Handlers: []model.PaymentHandler{}}
	if h.services != nil && h.services.Handler != nil {
		if handlers, err := h.services.Handler.List(); err == nil {
			profile.Payment.Handlers = mapHandlers(handlers)
		}
	}

	c.JSON(http.StatusOK, profile)
}

func resolveBaseURL(c *gin.Context) string {
	proto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto"))
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	host := c.Request.Host
	return proto + "://" + host
}

func mapHandlers(handlers []*domain.PaymentHandler) []model.PaymentHandler {
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

func parseHandlerConfig(value string) any {
	if strings.TrimSpace(value) == "" {
		return map[string]interface{}{}
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(value), &decoded); err == nil {
		return decoded
	}

	return value
}
