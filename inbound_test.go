package inboundgo

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test successful client creation
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
	
	if client.apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", client.apiKey)
	}
	
	if client.baseURL != "https://inbound.new/api/v2" {
		t.Errorf("Expected default base URL 'https://inbound.new/api/v2', got '%s'", client.baseURL)
	}
}

func TestNewClientWithCustomBaseURL(t *testing.T) {
	customURL := "https://custom-api.example.com"
	client, err := NewClient("test-api-key", customURL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if client.baseURL != customURL {
		t.Errorf("Expected base URL '%s', got '%s'", customURL, client.baseURL)
	}
}

func TestNewClientEmptyAPIKey(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Fatal("Expected error for empty API key, got nil")
	}
	
	expectedError := "API key is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test String helper
	s := "test"
	ptr := String(s)
	if ptr == nil {
		t.Fatal("String() returned nil")
	}
	if *ptr != s {
		t.Errorf("Expected '%s', got '%s'", s, *ptr)
	}
	
	// Test Int helper
	i := 42
	intPtr := Int(i)
	if intPtr == nil {
		t.Fatal("Int() returned nil")
	}
	if *intPtr != i {
		t.Errorf("Expected %d, got %d", i, *intPtr)
	}
	
	// Test Bool helper
	b := true
	boolPtr := Bool(b)
	if boolPtr == nil {
		t.Fatal("Bool() returned nil")
	}
	if *boolPtr != b {
		t.Errorf("Expected %v, got %v", b, *boolPtr)
	}
}

func TestServiceInitialization(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Test that all services can be accessed
	if client.Mail() == nil {
		t.Error("Mail service should not be nil")
	}
	
	if client.Email() == nil {
		t.Error("Email service should not be nil")
	}
	
	if client.Domain() == nil {
		t.Error("Domain service should not be nil")
	}
	
	if client.Endpoint() == nil {
		t.Error("Endpoint service should not be nil")
	}
	
	// Test nested email address service
	emailService := client.Email()
	if emailService.Address == nil {
		t.Error("Email address service should not be nil")
	}
}

func TestBuildQueryString(t *testing.T) {
	// Test with nil params
	result := buildQueryString(nil)
	if result != "" {
		t.Errorf("Expected empty string for nil params, got '%s'", result)
	}
	
	// Test with struct
	params := struct {
		Limit  *int    `json:"limit,omitempty"`
		Status string  `json:"status,omitempty"`
		Active *bool   `json:"active,omitempty"`
	}{
		Limit:  Int(10),
		Status: "verified",
		Active: Bool(true),
	}
	
	result = buildQueryString(params)
	// Should contain all parameters
	if result == "" {
		t.Error("Expected non-empty query string")
	}
	
	// Should start with ?
	if result[0] != '?' {
		t.Errorf("Expected query string to start with '?', got '%s'", result)
	}
}

func TestWithHTTPClient(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Test method chaining
	result := client.WithHTTPClient(nil)
	if result != client {
		t.Error("WithHTTPClient should return the same client instance")
	}
}
