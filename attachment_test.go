package inboundgo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
			name:          "email not found",
			emailID:       "non-existent",
			filename:      "document.pdf",
			serverStatus:  http.StatusNotFound,
			expectError:   true,
			errorContains: "404",
		},
		{
			name:          "attachment not found",
			emailID:       "test-email-id",
			filename:      "missing.pdf",
			serverStatus:  http.StatusNotFound,
			expectError:   true,
			errorContains: "404",
		},
		{
			name:          "unauthorized",
			emailID:       "test-email-id",
			filename:      "document.pdf",
			serverStatus:  http.StatusUnauthorized,
			expectError:   true,
			errorContains: "401",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify method
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				// Verify URL path
				expectedPath := "/attachments/" + tt.emailID + "/" + url.PathEscape(tt.filename)
				if r.URL.EscapedPath() != expectedPath {
					t.Errorf("Expected path '%s', got '%s'", expectedPath, r.URL.EscapedPath())
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

			result, err := client.Attachment().Download(context.Background(), tt.emailID, tt.filename)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if string(result.Data) != string(tt.serverResponse) {
				t.Errorf("Expected data '%s', got '%s'", string(tt.serverResponse), string(result.Data))
			}

			// Verify headers are present
			if result.Headers == nil {
				t.Error("Expected headers to be present")
			}
			if result.Headers.Get("Content-Type") != "application/pdf" {
				t.Errorf("Expected Content-Type 'application/pdf', got '%s'", result.Headers.Get("Content-Type"))
			}
		})
	}
}