package inboundgo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAttachmentDownload(t *testing.T) {
	tests := []struct {
		name           string
		emailID        string
		filename       string
		serverResponse []byte
		serverStatus   int
		expectError    bool
		errorContains  string
	}{
		{
			name:           "successful download",
			emailID:        "test-email-id",
			filename:       "document.pdf",
			serverResponse: []byte("PDF file content here"),
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "filename with spaces",
			emailID:        "test-email-id",
			filename:       "my document.pdf",
			serverResponse: []byte("PDF file content"),
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:         "email not found",
			emailID:      "non-existent",
			filename:     "document.pdf",
			serverStatus: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:         "attachment not found",
			emailID:      "test-email-id",
			filename:     "missing.pdf",
			serverStatus: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:         "unauthorized",
			emailID:      "test-email-id",
			filename:     "document.pdf",
			serverStatus: http.StatusUnauthorized,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify method
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				// Verify auth header
				auth := r.Header.Get("Authorization")
				if auth != "Bearer test-api-key" {
					t.Errorf("Expected auth header 'Bearer test-api-key', got '%s'", auth)
				}

				// Send response
				if tt.serverStatus >= 400 {
					w.WriteHeader(tt.serverStatus)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "Error occurred",
					})
					return
				}

				w.Header().Set("Content-Type", "application/pdf")
				w.Header().Set("Content-Disposition", `attachment; filename="`+tt.filename+`"`)
				w.WriteHeader(tt.serverStatus)
				w.Write(tt.serverResponse)
			}))
			defer server.Close()

			client, err := NewClient("test-api-key", server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			data, err := client.Attachment().Download(context.Background(), tt.emailID, tt.filename)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if string(data) != string(tt.serverResponse) {
				t.Errorf("Expected data '%s', got '%s'", string(tt.serverResponse), string(data))
			}
		})
	}
}
