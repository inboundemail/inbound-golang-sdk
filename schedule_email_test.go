package inboundgo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/inboundemail/inbound-golang-sdk"
)

func TestEmailScheduling(t *testing.T) {
	t.Run("should schedule email with natural language date", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "scheduled-email-123",
				"scheduled_at": "2024-01-01T12:00:00Z",
				"status": "scheduled",
				"timezone": "America/New_York"
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Schedule(ctx, &inboundgo.PostScheduleEmailRequest{
			From:        "test@example.com",
			To:          "recipient@example.com",
			Subject:     "Scheduled Email Test",
			Text:        inboundgo.String("This email is scheduled for later"),
			HTML:        inboundgo.String("<p>This email is scheduled for later</p>"),
			ScheduledAt: "in 2 hours",
			Timezone:    inboundgo.String("America/New_York"),
		}, nil)

		if err != nil {
			t.Fatalf("Failed to schedule email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if response.Data.ID != "scheduled-email-123" {
			t.Errorf("Expected ID 'scheduled-email-123', got '%s'", response.Data.ID)
		}

		if response.Data.Status != "scheduled" {
			t.Errorf("Expected status 'scheduled', got '%s'", response.Data.Status)
		}

		if response.Data.Timezone != "America/New_York" {
			t.Errorf("Expected timezone 'America/New_York', got '%s'", response.Data.Timezone)
		}
	})

	t.Run("should schedule email with ISO 8601 date", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "scheduled-iso-456",
				"scheduled_at": "2024-12-25T09:00:00Z",
				"status": "scheduled",
				"timezone": "UTC"
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		futureDate := time.Now().Add(3 * time.Hour).Format(time.RFC3339)

		response, err := client.Email().Schedule(ctx, &inboundgo.PostScheduleEmailRequest{
			From:        "test@example.com",
			To:          "recipient@example.com",
			Subject:     "ISO Scheduled Email",
			Text:        inboundgo.String("This email uses ISO 8601 formatting"),
			ScheduledAt: futureDate,
		}, nil)

		if err != nil {
			t.Fatalf("Failed to schedule email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if response.Data.Status != "scheduled" {
			t.Errorf("Expected status 'scheduled', got '%s'", response.Data.Status)
		}
	})

	t.Run("should schedule email with attachments and CID images", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "scheduled-with-attachments",
				"scheduled_at": "2024-01-01T15:30:00Z",
				"status": "scheduled",
				"timezone": "UTC"
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		// Sample base64 content
		testImageBase64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=="

		response, err := client.Email().Schedule(ctx, &inboundgo.PostScheduleEmailRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Scheduled Email with Attachments",
			HTML: inboundgo.String(`
				<div>
					<h1>Scheduled Newsletter</h1>
					<img src="cid:newsletter-logo" alt="Logo" style="width: 200px;" />
					<p>This email was scheduled in advance!</p>
				</div>
			`),
			Text:        inboundgo.String("This scheduled email contains attachments"),
			ScheduledAt: "tomorrow at 10am",
			Attachments: []inboundgo.AttachmentData{
				{
					Content:     inboundgo.String(testImageBase64),
					Filename:    "newsletter-logo.png",
					ContentType: inboundgo.String("image/png"),
					ContentID:   inboundgo.String("newsletter-logo"),
				},
				{
					Content:     inboundgo.String("U2FtcGxlIFBERiBjb250ZW50"), // base64 for "Sample PDF content"
					Filename:    "newsletter.pdf",
					ContentType: inboundgo.String("application/pdf"),
				},
			},
		}, nil)

		if err != nil {
			t.Fatalf("Failed to schedule email with attachments: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if response.Data.Status != "scheduled" {
			t.Errorf("Expected status 'scheduled', got '%s'", response.Data.Status)
		}
	})

	t.Run("should schedule email with idempotency key", func(t *testing.T) {
		var capturedHeaders http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedHeaders = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "scheduled-idempotent",
				"scheduled_at": "2024-01-01T16:00:00Z",
				"status": "scheduled",
				"timezone": "UTC"
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Schedule(ctx, &inboundgo.PostScheduleEmailRequest{
			From:        "test@example.com",
			To:          "recipient@example.com",
			Subject:     "Idempotent Scheduled Email",
			Text:        inboundgo.String("This scheduled email has an idempotency key"),
			ScheduledAt: "in 4 hours",
		}, &inboundgo.IdempotencyOptions{
			IdempotencyKey: "unique-schedule-key-123",
		})

		if err != nil {
			t.Fatalf("Failed to schedule email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		// Check idempotency header was included
		idempotencyKey := capturedHeaders.Get("Idempotency-Key")
		if idempotencyKey != "unique-schedule-key-123" {
			t.Errorf("Expected Idempotency-Key 'unique-schedule-key-123', got '%s'", idempotencyKey)
		}
	})
}

func TestScheduledEmailManagement(t *testing.T) {
	t.Run("should list scheduled emails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"data": [
					{
						"id": "scheduled-1",
						"from": "test@example.com",
						"to": ["recipient@example.com"],
						"subject": "First Scheduled Email",
						"scheduled_at": "2024-01-01T10:00:00Z",
						"status": "scheduled",
						"timezone": "UTC",
						"created_at": "2024-01-01T08:00:00Z",
						"attempts": 0
					},
					{
						"id": "scheduled-2",
						"from": "test@example.com",
						"to": ["recipient@example.com"],
						"subject": "Second Scheduled Email",
						"scheduled_at": "2024-01-01T11:00:00Z",
						"status": "scheduled",
						"timezone": "America/New_York",
						"created_at": "2024-01-01T08:30:00Z",
						"attempts": 0
					}
				],
				"pagination": {
					"limit": 10,
					"offset": 0,
					"total": 2
				}
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().ListScheduled(ctx, &inboundgo.GetScheduledEmailsRequest{
			Limit:  inboundgo.Int(10),
			Status: "scheduled",
		})

		if err != nil {
			t.Fatalf("Failed to list scheduled emails: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if len(response.Data.Data) != 2 {
			t.Errorf("Expected 2 scheduled emails, got %d", len(response.Data.Data))
		}

		if response.Data.Data[0].ID != "scheduled-1" {
			t.Errorf("Expected first email ID 'scheduled-1', got '%s'", response.Data.Data[0].ID)
		}

		if response.Data.Data[0].Status != "scheduled" {
			t.Errorf("Expected first email status 'scheduled', got '%s'", response.Data.Data[0].Status)
		}
	})

	t.Run("should get specific scheduled email details", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "scheduled-details-123",
				"from": "sender@example.com",
				"to": ["recipient1@example.com", "recipient2@example.com"],
				"cc": ["cc@example.com"],
				"subject": "Detailed Scheduled Email",
				"text": "This is the text content",
				"html": "<p>This is the HTML content</p>",
				"headers": {"X-Custom": "value"},
				"attachments": [
					{
						"filename": "document.pdf",
						"contentType": "application/pdf"
					}
				],
				"tags": [{"name": "campaign", "value": "newsletter"}],
				"scheduled_at": "2024-01-01T14:00:00Z",
				"timezone": "America/New_York",
				"status": "scheduled",
				"attempts": 0,
				"max_attempts": 3,
				"created_at": "2024-01-01T10:00:00Z",
				"updated_at": "2024-01-01T10:00:00Z"
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().GetScheduled(ctx, "scheduled-details-123")
		if err != nil {
			t.Fatalf("Failed to get scheduled email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if response.Data.ID != "scheduled-details-123" {
			t.Errorf("Expected ID 'scheduled-details-123', got '%s'", response.Data.ID)
		}

		if response.Data.Status != "scheduled" {
			t.Errorf("Expected status 'scheduled', got '%s'", response.Data.Status)
		}

		if len(response.Data.To) != 2 {
			t.Errorf("Expected 2 recipients, got %d", len(response.Data.To))
		}

		if response.Data.Timezone != "America/New_York" {
			t.Errorf("Expected timezone 'America/New_York', got '%s'", response.Data.Timezone)
		}
	})

	t.Run("should cancel a scheduled email", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("Expected DELETE method, got %s", r.Method)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "cancelled-email-789",
				"status": "cancelled",
				"cancelled_at": "2024-01-01T12:30:00Z"
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Cancel(ctx, "cancelled-email-789")
		if err != nil {
			t.Fatalf("Failed to cancel scheduled email: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if response.Data.ID != "cancelled-email-789" {
			t.Errorf("Expected ID 'cancelled-email-789', got '%s'", response.Data.ID)
		}

		if response.Data.Status != "cancelled" {
			t.Errorf("Expected status 'cancelled', got '%s'", response.Data.Status)
		}

		if response.Data.CancelledAt == "" {
			t.Error("Expected cancelled_at timestamp, got empty string")
		}
	})
}

func TestSchedulingErrors(t *testing.T) {
	t.Run("should handle invalid scheduled_at date", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Invalid date format: 'invalid date format'"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Schedule(ctx, &inboundgo.PostScheduleEmailRequest{
			From:        "test@example.com",
			To:          "recipient@example.com",
			Subject:     "Invalid Schedule Test",
			Text:        inboundgo.String("This should fail"),
			ScheduledAt: "invalid date format",
		}, nil)

		if err != nil {
			t.Fatalf("Expected API response, got error: %v", err)
		}

		if response.Error == "" {
			t.Error("Expected error in response")
		}

		if !strings.Contains(response.Error, "Invalid date format") {
			t.Errorf("Expected 'Invalid date format' error, got: %s", response.Error)
		}
	})

	t.Run("should handle past date scheduling", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Cannot schedule emails in the past"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		pastDate := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)

		response, err := client.Email().Schedule(ctx, &inboundgo.PostScheduleEmailRequest{
			From:        "test@example.com",
			To:          "recipient@example.com",
			Subject:     "Past Date Test",
			Text:        inboundgo.String("This should fail"),
			ScheduledAt: pastDate,
		}, nil)

		if err != nil {
			t.Fatalf("Expected API response, got error: %v", err)
		}

		if response.Error == "" {
			t.Error("Expected error in response")
		}

		if !strings.Contains(response.Error, "Cannot schedule emails in the past") {
			t.Errorf("Expected past date error, got: %s", response.Error)
		}
	})

	t.Run("should handle cancelling non-existent email", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Scheduled email not found"}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.Email().Cancel(ctx, "non-existent-email-id")
		if err != nil {
			t.Fatalf("Expected API response, got error: %v", err)
		}

		if response.Error == "" {
			t.Error("Expected error in response")
		}

		if !strings.Contains(response.Error, "not found") {
			t.Errorf("Expected 'not found' error, got: %s", response.Error)
		}
	})
}

func TestConvenienceSchedulingMethods(t *testing.T) {
	t.Run("should use ScheduleReminder convenience method", func(t *testing.T) {
		var capturedHeaders http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedHeaders = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "reminder-email",
				"scheduled_at": "2024-01-02T09:00:00Z",
				"status": "scheduled",
				"timezone": "UTC"
			}`))
		}))
		defer server.Close()

		client, err := inboundgo.NewClient("test-api-key", server.URL)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		ctx := context.Background()

		response, err := client.ScheduleReminder(ctx, "user@example.com", "Meeting Reminder", "tomorrow at 9am", "reminders@example.com", &inboundgo.IdempotencyOptions{
			IdempotencyKey: "reminder-key-456",
		})

		if err != nil {
			t.Fatalf("Failed to schedule reminder: %v", err)
		}

		if response.Error != "" {
			t.Errorf("Expected no error, got: %s", response.Error)
		}

		if response.Data == nil {
			t.Fatal("Expected response data, got nil")
		}

		if response.Data.Status != "scheduled" {
			t.Errorf("Expected status 'scheduled', got '%s'", response.Data.Status)
		}

		// Check idempotency header was included
		idempotencyKey := capturedHeaders.Get("Idempotency-Key")
		if idempotencyKey != "reminder-key-456" {
			t.Errorf("Expected Idempotency-Key 'reminder-key-456', got '%s'", idempotencyKey)
		}
	})
}
