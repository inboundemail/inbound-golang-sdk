package inboundgo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inboundemail/inbound-golang-sdk"
)

// Helper function to check if error is a network-related error
func isNetworkError(err error) bool {
	if err != nil && !isNetworkError(err) {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "connect") || 
		   strings.Contains(errStr, "network") || 
		   strings.Contains(errStr, "timeout") || 
		   strings.Contains(errStr, "EOF") ||
		   strings.Contains(errStr, "no such host") ||
		   strings.Contains(errStr, "connection refused")
}

func TestHierarchicalStructure(t *testing.T) {
	client, err := inboundgo.NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("should have hierarchical email.address structure", func(t *testing.T) {
		emailService := client.Email()
		if emailService == nil {
			t.Fatal("Email service should not be nil")
		}

		if emailService.Address == nil {
			t.Fatal("Email address service should not be nil")
		}

		// Test that all methods exist
		ctx := context.Background()

		// These will fail with network errors but we're just testing structure
		_, err := emailService.Address.Create(ctx, &inboundgo.PostEmailAddressesRequest{
			Address:  "test@example.com",
			DomainID: "test-domain",
		})
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.Address.List(ctx, nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.Address.Get(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.Address.Update(ctx, "test-id", &inboundgo.PutEmailAddressByIDRequest{})
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.Address.Delete(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}
	})

	t.Run("should have mail methods for inbound emails", func(t *testing.T) {
		mailService := client.Mail()
		if mailService == nil {
			t.Fatal("Mail service should not be nil")
		}

		ctx := context.Background()

		// Test all methods exist (will fail with network errors)
		_, err := mailService.List(ctx, nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = mailService.Get(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = mailService.Thread(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = mailService.MarkRead(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = mailService.MarkUnread(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = mailService.Archive(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = mailService.Unarchive(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = mailService.Reply(ctx, &inboundgo.PostMailRequest{
			EmailID:  "test-id",
			To:       "test@example.com",
			Subject:  "Test",
			TextBody: "Test",
		})
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = mailService.Bulk(ctx, []string{"id1", "id2"}, map[string]interface{}{
			"isRead": true,
		})
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}
	})

	t.Run("should have email methods for outbound emails", func(t *testing.T) {
		emailService := client.Email()
		if emailService == nil {
			t.Fatal("Email service should not be nil")
		}

		ctx := context.Background()

		// Test all methods exist
		_, err := emailService.Send(ctx, &inboundgo.PostEmailsRequest{
			From:    "test@example.com",
			To:      "recipient@example.com",
			Subject: "Test",
			Text:    inboundgo.String("Test"),
		}, nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.Get(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.Reply(ctx, "test-id", &inboundgo.PostEmailReplyRequest{
			From: "test@example.com",
			Text: inboundgo.String("Test"),
		}, nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.Schedule(ctx, &inboundgo.PostScheduleEmailRequest{
			From:        "test@example.com",
			To:          "recipient@example.com",
			Subject:     "Test",
			ScheduledAt: "in 1 hour",
		}, nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.ListScheduled(ctx, nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.GetScheduled(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = emailService.Cancel(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}
	})

	t.Run("should have domain methods", func(t *testing.T) {
		domainService := client.Domain()
		if domainService == nil {
			t.Fatal("Domain service should not be nil")
		}

		ctx := context.Background()

		// Test all methods exist
		_, err := domainService.Create(ctx, &inboundgo.PostDomainsRequest{
			Domain: "test.com",
		})
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = domainService.List(ctx, nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = domainService.Get(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = domainService.Update(ctx, "test-id", &inboundgo.PutDomainByIDRequest{
			IsCatchAllEnabled:  true,
			CatchAllEndpointID: inboundgo.String("endpoint-id"),
		})
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = domainService.Delete(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = domainService.Verify(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = domainService.GetDNSRecords(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = domainService.CheckStatus(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}
	})

	t.Run("should have endpoint methods", func(t *testing.T) {
		endpointService := client.Endpoint()
		if endpointService == nil {
			t.Fatal("Endpoint service should not be nil")
		}

		ctx := context.Background()

		// Test all methods exist
		_, err := endpointService.Create(ctx, &inboundgo.PostEndpointsRequest{
			Name: "Test Endpoint",
			Type: "webhook",
			Config: &inboundgo.WebhookConfig{
				URL:           "https://example.com/webhook",
				Timeout:       30000,
				RetryAttempts: 3,
			},
		})
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = endpointService.List(ctx, nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = endpointService.Get(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = endpointService.Update(ctx, "test-id", &inboundgo.PutEndpointByIDRequest{})
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = endpointService.Delete(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = endpointService.Test(ctx, "test-id")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}
	})

	t.Run("should have convenience methods", func(t *testing.T) {
		ctx := context.Background()

		// Test all convenience methods exist
		_, err := client.QuickReply(ctx, "email-id", "message", "from@example.com", nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = client.SetupDomain(ctx, "test.com", inboundgo.String("https://webhook.example.com"))
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = client.CreateForwarder(ctx, "from@example.com", "to@example.com")
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}

		_, err = client.ScheduleReminder(ctx, "user@example.com", "Test Subject", "in 1 hour", "from@example.com", nil)
		if err != nil && !isNetworkError(err) {
			t.Errorf("Expected network error or nil, got: %v", err)
		}
	})
}

func TestAPIResponsePattern(t *testing.T) {
	// Create a test server that returns success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": [{"id": "test-123", "domain": "example.com"}], "pagination": {"limit": 10, "offset": 0, "total": 1}}`))
	}))
	defer server.Close()

	client, err := inboundgo.NewClient("test-api-key", server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	t.Run("should return { data, error } pattern on success", func(t *testing.T) {
		result, err := client.Domain().List(ctx, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.Error != "" {
			t.Errorf("Expected no error in response, got %s", result.Error)
		}

		if result.Data == nil {
			t.Error("Expected data in response, got nil")
		}

		if len(result.Data.Data) == 0 {
			t.Error("Expected domains in response data")
		}
	})

	// Create a test server that returns error
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid API key"}`))
	}))
	defer errorServer.Close()

	errorClient, err := inboundgo.NewClient("test-api-key", errorServer.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("should return { data, error } pattern on error", func(t *testing.T) {
		result, err := errorClient.Domain().List(ctx, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.Error == "" {
			t.Error("Expected error in response")
		}

		if result.Data != nil {
			t.Error("Expected no data in error response")
		}

		if result.Error != "Invalid API key" {
			t.Errorf("Expected 'Invalid API key', got %s", result.Error)
		}
	})

	// Test network errors
	t.Run("should handle network errors", func(t *testing.T) {
		// Use invalid URL to trigger network error
		networkClient, err := inboundgo.NewClient("test-api-key", "http://invalid-url-that-does-not-exist.com")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		result, err := networkClient.Domain().List(ctx, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// The response should contain either an error or data, but not both
		if result.Error == "" && result.Data == nil {
			t.Error("Expected either error or data in response")
		}

		if result.Error != "" && result.Data != nil {
			t.Error("Expected either error or data, not both")
		}
	})
}
