package inboundgo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inboundgo "github.com/inboundemail/inbound-golang-sdk"
)

func TestIdempotencyKeySupport(t *testing.T) {
	tests := []struct {
		name           string
		expectedHeader string
		testFunc       func(client *inboundgo.Inbound, ctx context.Context) error
	}{
		{
			name:           "email.Send() with idempotency key",
			expectedHeader: "test-key-123",
			testFunc: func(client *inboundgo.Inbound, ctx context.Context) error {
				_, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
					From:    "test@example.com",
					To:      "user@example.com",
					Subject: "Test Email",
					Text:    inboundgo.String("Test message"),
				}, &inboundgo.IdempotencyOptions{
					IdempotencyKey: "test-key-123",
				})
				return err
			},
		},
		{
			name:           "email.Schedule() with idempotency key",
			expectedHeader: "scheduled-key-456",
			testFunc: func(client *inboundgo.Inbound, ctx context.Context) error {
				_, err := client.Email().Schedule(ctx, &inboundgo.PostScheduleEmailRequest{
					From:        "test@example.com",
					To:          "user@example.com",
					Subject:     "Scheduled Email",
					Text:        inboundgo.String("Scheduled message"),
					ScheduledAt: "tomorrow at 10am",
				}, &inboundgo.IdempotencyOptions{
					IdempotencyKey: "scheduled-key-456",
				})
				return err
			},
		},
		{
			name:           "email.Reply() with idempotency key",
			expectedHeader: "reply-key-789",
			testFunc: func(client *inboundgo.Inbound, ctx context.Context) error {
				_, err := client.Email().Reply(ctx, "original-email-123", &inboundgo.PostEmailReplyRequest{
					From: "support@example.com",
					Text: inboundgo.String("Reply message"),
				}, &inboundgo.IdempotencyOptions{
					IdempotencyKey: "reply-key-789",
				})
				return err
			},
		},
		{
			name:           "QuickReply() convenience method with idempotency key",
			expectedHeader: "quick-reply-key",
			testFunc: func(client *inboundgo.Inbound, ctx context.Context) error {
				_, err := client.QuickReply(ctx, "email-123", "Thanks for your message!", "support@example.com", &inboundgo.IdempotencyOptions{
					IdempotencyKey: "quick-reply-key",
				})
				return err
			},
		},
		{
			name:           "ScheduleReminder() convenience method with idempotency key",
			expectedHeader: "reminder-key-456",
			testFunc: func(client *inboundgo.Inbound, ctx context.Context) error {
				_, err := client.ScheduleReminder(ctx, "user@example.com", "Meeting Reminder", "tomorrow at 9am", "reminders@example.com", &inboundgo.IdempotencyOptions{
					IdempotencyKey: "reminder-key-456",
				})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that captures the request headers
			var capturedHeaders http.Header
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedHeaders = r.Header.Clone()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": "test-123", "messageId": "msg-456"}`))
			}))
			defer server.Close()

			client, err := inboundgo.NewClient("test-api-key", server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			ctx := context.Background()

			// Execute the test function
			err = tt.testFunc(client, ctx)
			if err != nil {
				// We expect network-related errors but not nil pointer errors
				if !strings.Contains(err.Error(), "connect") && !strings.Contains(err.Error(), "EOF") {
					t.Fatalf("Unexpected error: %v", err)
				}
			}

			// Check that the Idempotency-Key header was included
			if capturedHeaders == nil {
				t.Fatal("No headers captured")
			}

			idempotencyKey := capturedHeaders.Get("Idempotency-Key")
			if idempotencyKey != tt.expectedHeader {
				t.Errorf("Expected Idempotency-Key header '%s', got '%s'", tt.expectedHeader, idempotencyKey)
			}

			// Check that Content-Type is set
			contentType := capturedHeaders.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			// Check that Authorization header is set
			auth := capturedHeaders.Get("Authorization")
			if auth != "Bearer test-api-key" {
				t.Errorf("Expected Authorization 'Bearer test-api-key', got '%s'", auth)
			}
		})
	}

	t.Run("should work without idempotency key", func(t *testing.T) {
		var capturedHeaders http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedHeaders = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "test-456", "messageId": "msg-789"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		_, err = client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "user@example.com",
			Subject: "Test Email",
			Text:    inboundgo.String("Test message"),
		}, nil) // No idempotency options

		if err != nil && !strings.Contains(err.Error(), "connect") && !strings.Contains(err.Error(), "EOF") {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Should not include Idempotency-Key header
		idempotencyKey := capturedHeaders.Get("Idempotency-Key")
		if idempotencyKey != "" {
			t.Errorf("Expected no Idempotency-Key header, got '%s'", idempotencyKey)
		}

		// Should still include other headers
		contentType := capturedHeaders.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}
	})

	t.Run("should handle empty idempotency key", func(t *testing.T) {
		var capturedHeaders http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedHeaders = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "test-789", "messageId": "msg-012"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		_, err = client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "user@example.com",
			Subject: "Test Email",
			Text:    inboundgo.String("Test message"),
		}, &inboundgo.IdempotencyOptions{
			IdempotencyKey: "", // Empty key
		})

		if err != nil && !strings.Contains(err.Error(), "connect") && !strings.Contains(err.Error(), "EOF") {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Should not include Idempotency-Key header when empty
		idempotencyKey := capturedHeaders.Get("Idempotency-Key")
		if idempotencyKey != "" {
			t.Errorf("Expected no Idempotency-Key header for empty key, got '%s'", idempotencyKey)
		}
	})
}
