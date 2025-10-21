package inboundgo

import (
	"strings"
	"testing"
)

func TestParseWebhookPayload(t *testing.T) {
	payload := `{
  "event": "email.received",
  "timestamp": "2025-09-16T16:47:50.163Z",
  "email": {
    "id": "7U6TcAy-16qmzu297IVoL",
    "messageId": "<test-yaZbgt70Z4J6XKIiAAEvZ@mail.inbound.new>",
    "from": {
      "text": "Inbound Test <test@example.com>",
      "addresses": [
        {
          "name": "Inbound Test",
          "address": "test@example.com"
        }
      ]
    },
    "to": {
      "text": "Test Recipient <test@yourdomain.com>",
      "addresses": [
        {
          "name": null,
          "address": "test@yourdomain.com"
        }
      ]
    },
    "recipient": "test@yourdomain.com",
    "subject": "Test Email - Inbound Email Service",
    "receivedAt": "2025-09-16T16:47:50.163Z",
    "parsedData": {
      "messageId": "<test-yaZbgt70Z4J6XKIiAAEvZ@mail.inbound.new>",
      "date": "2025-09-16T16:47:50.163Z",
      "subject": "Test Email - Inbound Email Service",
      "from": {
        "text": "Inbound Test <test@example.com>",
        "addresses": [
          {
            "name": "Inbound Test",
            "address": "test@example.com"
          }
        ]
      },
      "to": {
        "text": "Test Recipient <test@yourdomain.com>",
        "addresses": [
          {
            "name": null,
            "address": "test@yourdomain.com"
          }
        ]
      },
      "cc": null,
      "bcc": null,
      "replyTo": null,
      "textBody": "This is a test email.\nRendered for webhook testing.",
      "htmlBody": "<div><p>This is a test email.</p><p><strong>Rendered for webhook testing.</strong></p></div>",
      "attachments": [],
      "headers": {
        "received": [
          "from test-mta.inbound.new",
          "by test-mx.google.com"
        ],
        "received-spf": "pass",
        "dkim-signature": {
          "value": "v=1",
          "params": {
            "a": "rsa-sha256",
            "d": "example.com"
          }
        }
      }
    },
    "cleanedContent": {
      "html": "<div><p>This is a test email.</p><p><strong>Rendered for webhook testing.</strong></p></div>",
      "text": "This is a test email.\nRendered for webhook testing.",
      "hasHtml": true,
      "hasText": true,
      "attachments": [],
      "headers": {}
    }
  },
  "endpoint": {
    "id": "LHbWZ1iEOofDXlViXWsDH",
    "name": "6979012cd152 E",
    "type": "webhook"
  }
}`

	webhook, err := ParseWebhookPayload(strings.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to parse webhook payload: %v", err)
	}

	// Test basic webhook fields
	if webhook.Event != "email.received" {
		t.Errorf("Expected event 'email.received', got '%s'", webhook.Event)
	}

	if webhook.Timestamp != "2025-09-16T16:47:50.163Z" {
		t.Errorf("Expected timestamp '2025-09-16T16:47:50.163Z', got '%s'", webhook.Timestamp)
	}

	// Test email fields
	if webhook.Email.ID != "7U6TcAy-16qmzu297IVoL" {
		t.Errorf("Expected email ID '7U6TcAy-16qmzu297IVoL', got '%s'", webhook.Email.ID)
	}

	if webhook.Email.MessageID == nil || *webhook.Email.MessageID != "<test-yaZbgt70Z4J6XKIiAAEvZ@mail.inbound.new>" {
		if webhook.Email.MessageID == nil {
			t.Error("Expected message ID to be present")
		} else {
			t.Errorf("Expected message ID '<test-yaZbgt70Z4J6XKIiAAEvZ@mail.inbound.new>', got '%s'", *webhook.Email.MessageID)
		}
	}

	if webhook.Email.Subject == nil || *webhook.Email.Subject != "Test Email - Inbound Email Service" {
		if webhook.Email.Subject == nil {
			t.Error("Expected subject to be present")
		} else {
			t.Errorf("Expected subject 'Test Email - Inbound Email Service', got '%s'", *webhook.Email.Subject)
		}
	}

	if webhook.Email.Recipient != "test@yourdomain.com" {
		t.Errorf("Expected recipient 'test@yourdomain.com', got '%s'", webhook.Email.Recipient)
	}

	// Test helper methods
	fromAddr := webhook.GetFromAddress()
	if fromAddr != "Inbound Test <test@example.com>" {
		t.Errorf("Expected from address 'Inbound Test <test@example.com>', got '%s'", fromAddr)
	}

	toAddr := webhook.GetToAddress()
	if toAddr != "test@yourdomain.com" {
		t.Errorf("Expected to address 'test@yourdomain.com', got '%s'", toAddr)
	}

	// Test parsed data
	if webhook.Email.ParsedData.TextBody == nil || *webhook.Email.ParsedData.TextBody != "This is a test email.\nRendered for webhook testing." {
		if webhook.Email.ParsedData.TextBody == nil {
			t.Error("Expected text body to be present")
		} else {
			t.Errorf("Expected text body 'This is a test email.\\nRendered for webhook testing.', got '%s'", *webhook.Email.ParsedData.TextBody)
		}
	}

	if webhook.Email.ParsedData.HTMLBody == nil || *webhook.Email.ParsedData.HTMLBody != "<div><p>This is a test email.</p><p><strong>Rendered for webhook testing.</strong></p></div>" {
		if webhook.Email.ParsedData.HTMLBody == nil {
			t.Error("Expected HTML body to be present")
		} else {
			t.Errorf("Expected HTML body '<div><p>This is a test email.</p><p><strong>Rendered for webhook testing.</strong></p></div>', got '%s'", *webhook.Email.ParsedData.HTMLBody)
		}
	}

	// Test cleaned content
	if webhook.Email.CleanedContent.Text == nil || *webhook.Email.CleanedContent.Text != "This is a test email.\nRendered for webhook testing." {
		if webhook.Email.CleanedContent.Text == nil {
			t.Error("Expected cleaned content text to be present")
		} else {
			t.Errorf("Expected cleaned content text 'This is a test email.\\nRendered for webhook testing.', got '%s'", *webhook.Email.CleanedContent.Text)
		}
	}

	if webhook.Email.CleanedContent.HTML == nil || *webhook.Email.CleanedContent.HTML != "<div><p>This is a test email.</p><p><strong>Rendered for webhook testing.</strong></p></div>" {
		if webhook.Email.CleanedContent.HTML == nil {
			t.Error("Expected cleaned content HTML to be present")
		} else {
			t.Errorf("Expected cleaned content HTML '<div><p>This is a test email.</p><p><strong>Rendered for webhook testing.</strong></p></div>', got '%s'", *webhook.Email.CleanedContent.HTML)
		}
	}

	if !webhook.Email.CleanedContent.HasHTML {
		t.Error("Expected cleaned content to have HTML")
	}

	if !webhook.Email.CleanedContent.HasText {
		t.Error("Expected cleaned content to have text")
	}

	// Test headers parsing
	headers := webhook.GetHeaders()
	if len(headers["received"]) != 2 {
		t.Errorf("Expected 2 received headers, got %d", len(headers["received"]))
	}

	if headers["received-spf"][0] != "pass" {
		t.Errorf("Expected received-spf 'pass', got '%s'", headers["received-spf"][0])
	}

	// Test complex header parsing (dkim-signature as object)
	if _, exists := headers["dkim-signature"]; !exists {
		t.Error("Expected dkim-signature header to be parsed")
	}

	// Test endpoint
	if webhook.Endpoint.ID != "LHbWZ1iEOofDXlViXWsDH" {
		t.Errorf("Expected endpoint ID 'LHbWZ1iEOofDXlViXWsDH', got '%s'", webhook.Endpoint.ID)
	}

	if webhook.Endpoint.Name != "6979012cd152 E" {
		t.Errorf("Expected endpoint name '6979012cd152 E', got '%s'", webhook.Endpoint.Name)
	}

	if webhook.Endpoint.Type != "webhook" {
		t.Errorf("Expected endpoint type 'webhook', got '%s'", webhook.Endpoint.Type)
	}
}

func TestGetFromAddressWithoutName(t *testing.T) {
	payload := `{
  "event": "email.received",
  "timestamp": "2025-09-16T16:47:50.163Z",
  "email": {
    "from": {
      "text": "test@example.com",
      "addresses": [
        {
          "name": null,
          "address": "test@example.com"
        }
      ]
    },
    "parsedData": {
      "headers": {}
    }
  }
}`

	webhook, err := ParseWebhookPayload(strings.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to parse webhook payload: %v", err)
	}

	fromAddr := webhook.GetFromAddress()
	if fromAddr != "test@example.com" {
		t.Errorf("Expected from address 'test@example.com' (without name), got '%s'", fromAddr)
	}
}

func TestGetAddressesEmpty(t *testing.T) {
	payload := `{
  "event": "email.received",
  "timestamp": "2025-09-16T16:47:50.163Z",
  "email": {
    "from": {
      "text": "",
      "addresses": []
    },
    "to": {
      "text": "",
      "addresses": []
    },
    "parsedData": {
      "headers": {}
    }
  }
}`

	webhook, err := ParseWebhookPayload(strings.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to parse webhook payload: %v", err)
	}

	fromAddr := webhook.GetFromAddress()
	if fromAddr != "" {
		t.Errorf("Expected empty from address, got '%s'", fromAddr)
	}

	toAddr := webhook.GetToAddress()
	if toAddr != "" {
		t.Errorf("Expected empty to address, got '%s'", toAddr)
	}
}
