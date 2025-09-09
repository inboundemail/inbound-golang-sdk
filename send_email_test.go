package inboundgo_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/R44VC0RP/inbound-golang-sdk"
)

func TestEmailSending(t *testing.T) {
	// Sample base64 content for testing (small PNG image)
	sampleBase64PNG := base64.StdEncoding.EncodeToString([]byte("fake-png-data"))
	
	t.Run("should send email with basic content", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "email-123", "messageId": "msg-456", "status": "sent"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Test Email",
			Text:    inboundgo.String("This is a test email"),
			HTML:    inboundgo.String("<p>This is a test email</p>"),
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if response.Data.ID != "email-123" {
			t.Errorf("Expected ID 'email-123', got '%s'", response.Data.ID)
		}
	})

	t.Run("should send email with single base64 attachment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the request contains attachment data
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "email-with-attachment", "messageId": "msg-789"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Test Email with Attachment",
			Text:    inboundgo.String("This email has an attachment"),
			Attachments: []inboundgo.AttachmentData{
				{
					Content:     inboundgo.String(sampleBase64PNG),
					Filename:    "test-image.png",
					ContentType: inboundgo.String("image/png"),
				},
			},
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}
	})

	t.Run("should send email with multiple attachments", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "email-multi-attachment", "messageId": "msg-multi"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Test Email with Multiple Attachments",
			Text:    inboundgo.String("This email has multiple attachments"),
			Attachments: []inboundgo.AttachmentData{
				{
					Content:     inboundgo.String(sampleBase64PNG),
					Filename:    "image1.png",
					ContentType: inboundgo.String("image/png"),
				},
				{
					Content:     inboundgo.String(base64.StdEncoding.EncodeToString([]byte("Sample document content"))),
					Filename:    "document.txt",
					ContentType: inboundgo.String("text/plain"),
				},
			},
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}
	})

	t.Run("should send email with remote file attachment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "email-remote-attachment", "messageId": "msg-remote"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Test Email with Remote Attachment",
			Text:    inboundgo.String("This email has a remote attachment"),
			Attachments: []inboundgo.AttachmentData{
				{
					Path:     inboundgo.String("https://httpbin.org/image/png"),
					Filename: "remote-image.png",
				},
			},
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}
	})

	t.Run("should send email with CID image embedding", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "email-cid", "messageId": "msg-cid"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Test Email with CID Image",
			HTML:    inboundgo.String(`<p>Check out our logo: <img src="cid:company-logo" alt="Logo" /></p>`),
			Text:    inboundgo.String("This email has an embedded image"),
			Attachments: []inboundgo.AttachmentData{
				{
					Content:     inboundgo.String(sampleBase64PNG),
					Filename:    "logo.png",
					ContentType: inboundgo.String("image/png"),
					ContentID:   inboundgo.String("company-logo"),
				},
			},
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}
	})

	t.Run("should send email with custom headers and tags", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "email-headers-tags", "messageId": "msg-headers"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Test Email with Headers and Tags",
			Text:    inboundgo.String("This email has custom headers and tags"),
			Headers: map[string]string{
				"X-Custom-Header": "test-value",
				"X-Priority":      "high",
			},
			Tags: []inboundgo.EmailTag{
				{Name: "category", Value: "test"},
				{Name: "source", Value: "go-sdk"},
			},
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}
	})

	t.Run("should send email to multiple recipients", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "email-multi-recipients", "messageId": "msg-multi-to"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      []string{"recipient1@example.com", "recipient2@example.com"},
			CC:      []string{"cc@example.com"},
			BCC:     []string{"bcc@example.com"},
			Subject: "Test Email to Multiple Recipients",
			Text:    inboundgo.String("This email goes to multiple people"),
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}
	})

	t.Run("should handle scheduled email sending", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this is a scheduled email (endpoint will be /emails/schedule)
			if strings.Contains(r.URL.Path, "schedule") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": "scheduled-email", "scheduled_at": "2024-01-01T10:00:00Z", "status": "scheduled"}`))
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": "immediate-email", "messageId": "msg-immediate", "status": "sent"}`))
			}
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		// Test immediate send
		response1, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Immediate Email",
			Text:    inboundgo.String("This email is sent immediately"),
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send immediate email: %v", err)
		}

		if response1.Error != "" {
			t.Errorf("Expected no error, got: %s", response1.Error)
		}

		// Test scheduled send
		response2, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:        "test@example.com",
			To:          "recipient@example.com",
			Subject:     "Scheduled Email",
			Text:        inboundgo.String("This email is scheduled"),
			ScheduledAt: inboundgo.String("in 1 hour"),
		}, nil)

		if err != nil {
			t.Fatalf("Failed to send scheduled email: %v", err)
		}

		if response2.Error != "" {
			t.Errorf("Expected no error, got: %s", response2.Error)
		}
	})
}

func TestEmailSendingErrors(t *testing.T) {
	t.Run("should handle missing required fields", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Missing required fields: to, subject"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From: "test@example.com",
			// Missing To and Subject
			Text: inboundgo.String("Test email"),
		}, nil)

		if err != nil {
			t.Fatalf("Expected API response, got error: %v", err)
		}

		if response.Error == "" {
			t.Error("Expected error in response")
		}

		if !strings.Contains(response.Error, "Missing required fields") {
			t.Errorf("Expected 'Missing required fields' error, got: %s", response.Error)
		}
	})

	t.Run("should handle unauthorized domain", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "You don't have permission to send from domain unauthorized-domain.com"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@unauthorized-domain.com",
			To:      "recipient@example.com",
			Subject: "Test Subject",
			Text:    inboundgo.String("Test content"),
		}, nil)

		if err != nil {
			t.Fatalf("Expected API response, got error: %v", err)
		}

		if response.Error == "" {
			t.Error("Expected error in response")
		}

		if !strings.Contains(response.Error, "don't have permission") {
			t.Errorf("Expected permission error, got: %s", response.Error)
		}
	})

	t.Run("should handle invalid attachment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Attachment validation failed: filename is required"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Test with Invalid Attachment",
			Text:    inboundgo.String("Test content"),
			Attachments: []inboundgo.AttachmentData{
				{
					Content: inboundgo.String("test-content"),
					// Missing Filename
				},
			},
		}, nil)

		if err != nil {
			t.Fatalf("Expected API response, got error: %v", err)
		}

		if response.Error == "" {
			t.Error("Expected error in response")
		}

		if !strings.Contains(response.Error, "filename is required") {
			t.Errorf("Expected filename error, got: %s", response.Error)
		}
	})
}

func TestGetSentEmail(t *testing.T) {
	t.Run("should retrieve sent email by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"object": "email",
				"id": "test-email-123",
				"from": "sender@example.com",
				"to": ["recipient@example.com"],
				"cc": [],
				"bcc": [],
				"reply_to": [],
				"subject": "Test Email",
				"text": "Test content",
				"html": "<p>Test content</p>",
				"created_at": "2024-01-01T10:00:00Z",
				"last_event": "delivered"
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Get(ctx, "test-email-123")
		if err != nil {
			t.Fatalf("Failed to get email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if response.Data.ID != "test-email-123" {
			t.Errorf("Expected ID 'test-email-123', got '%s'", response.Data.ID)
		}

		if response.Data.Subject != "Test Email" {
			t.Errorf("Expected subject 'Test Email', got '%s'", response.Data.Subject)
		}

		if response.Data.LastEvent != "delivered" {
			t.Errorf("Expected last_event 'delivered', got '%s'", response.Data.LastEvent)
		}
	})

	t.Run("should handle non-existent email", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Email not found"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Get(ctx, "non-existent-id")
		if err != nil {
			t.Fatalf("Expected API response, got error: %v", err)
		}

		if response.Error == "" {
			t.Error("Expected error in response")
		}

		if !strings.Contains(response.Error, "Email not found") {
			t.Errorf("Expected 'Email not found' error, got: %s", response.Error)
		}
	})
}
