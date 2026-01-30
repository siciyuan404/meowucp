package model

type CheckoutCreateRequest struct {
	LineItems []LineItem `json:"line_items"`
	Currency  string     `json:"currency"`
}

type CheckoutUpdateRequest struct {
	ID        string     `json:"id"`
	LineItems []LineItem `json:"line_items"`
	Currency  string     `json:"currency"`
}

type CheckoutCompleteRequest struct {
	PaymentData PaymentInstrument `json:"payment_data"`
}

type CheckoutSession struct {
	UCP         *UCPMeta   `json:"ucp,omitempty"`
	ID          string     `json:"id"`
	LineItems   []LineItem `json:"line_items"`
	Status      string     `json:"status"`
	Currency    string     `json:"currency"`
	Totals      []Total    `json:"totals"`
	Messages    []Message  `json:"messages,omitempty"`
	Links       []Link     `json:"links"`
	ContinueURL string     `json:"continue_url,omitempty"`
	Payment     Payment    `json:"payment"`
	Order       *OrderRef  `json:"order,omitempty"`
}

type UCPMeta struct {
	Version      string           `json:"version"`
	Capabilities []CapabilityMeta `json:"capabilities,omitempty"`
}

type CapabilityMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type LineItem struct {
	ID       string `json:"id,omitempty"`
	Item     Item   `json:"item"`
	Quantity int    `json:"quantity"`
}

type Item struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Price    int64  `json:"price"`
	ImageURL string `json:"image_url,omitempty"`
}

type Total struct {
	Type   string `json:"type"`
	Amount int64  `json:"amount"`
}

type Message struct {
	Type     string `json:"type"`
	Code     string `json:"code,omitempty"`
	Content  string `json:"content"`
	Severity string `json:"severity,omitempty"`
}

type Link struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

type Payment struct {
	Handlers []PaymentHandler `json:"handlers"`
}

type PaymentInstrument struct {
	ID         string `json:"id,omitempty"`
	HandlerID  string `json:"handler_id"`
	Type       string `json:"type"`
	Credential any    `json:"credential,omitempty"`
}

type OrderRef struct {
	ID           string `json:"id"`
	PermalinkURL string `json:"permalink_url,omitempty"`
}
