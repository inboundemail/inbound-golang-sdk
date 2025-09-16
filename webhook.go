package inboundgo

import (
	"encoding/json"
	"fmt"
	"io"
)

// ParseWebhookPayload parses an incoming webhook payload into the WebhookPayload struct
func ParseWebhookPayload(reader io.Reader) (*WebhookPayload, error) {
	var payload WebhookPayload
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}
	return &payload, nil
}

// GetFromAddress extracts the properly formatted from address from the webhook
func (w *WebhookPayload) GetFromAddress() string {
	if len(w.Email.From.Addresses) > 0 {
		addr := w.Email.From.Addresses[0]
		if addr.Name != nil && *addr.Name != "" {
			return fmt.Sprintf("%s <%s>", *addr.Name, addr.Address)
		}
		return addr.Address
	}
	return ""
}

// GetToAddress extracts the properly formatted to address from the webhook
func (w *WebhookPayload) GetToAddress() string {
	if len(w.Email.To.Addresses) > 0 {
		addr := w.Email.To.Addresses[0]
		if addr.Name != nil && *addr.Name != "" {
			return fmt.Sprintf("%s <%s>", *addr.Name, addr.Address)
		}
		return addr.Address
	}
	return ""
}

// GetHeaders converts the headers from the webhook format to a standard map[string][]string format
func (w *WebhookPayload) GetHeaders() map[string][]string {
	headers := make(map[string][]string)
	for k, v := range w.Email.ParsedData.Headers {
		switch val := v.(type) {
		case string:
			headers[k] = []string{val}
		case []string:
			headers[k] = val
		case []any:
			var strSlice []string
			for _, item := range val {
				if str, ok := item.(string); ok {
					strSlice = append(strSlice, str)
				}
			}
			if len(strSlice) > 0 {
				headers[k] = strSlice
			}
		case map[string]any:
			// Handle complex header structures like dkim-signature
			if text, ok := val["text"].(string); ok {
				headers[k] = []string{text}
			} else if value, ok := val["value"].(string); ok {
				headers[k] = []string{value}
			}
		}
	}
	return headers
}
