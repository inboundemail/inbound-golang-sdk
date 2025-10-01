// Package inboundgo provides a Go client for the Inbound Email API.
//
// The Inbound Email API allows you to send emails, manage domains, endpoints,
// and email addresses. This client provides a simple interface to interact
// with all the available endpoints.
//
// Basic Usage:
//
//	client, err := inboundgo.NewClient("your-api-key")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Send an email
//	resp, err := client.Email().Send(ctx, &inboundgo.PostEmailsRequest{
//		From:    "sender@example.com",
//		To:      "recipient@example.com",
//		Subject: "Hello World",
//		Text:    inboundgo.String("Hello from Go!"),
//	}, nil)
//
//	// List inbound emails
//	emails, err := client.Mail().List(ctx, nil)
//
//	// Manage domains
//	domain, err := client.Domain().Create(ctx, &inboundgo.PostDomainsRequest{
//		Domain: "example.com",
//	})
//
// For detailed API documentation, see: https://docs.inbound.new/api-reference
package inboundgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Inbound is the main client for the Inbound Email SDK
type Inbound struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Inbound Email client
func NewClient(apiKey string, baseURL ...string) (*Inbound, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	url := "https://inbound.new/api/v2"
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = baseURL[0]
	}

	return &Inbound{
		apiKey:     apiKey,
		baseURL:    url,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// WithHTTPClient sets a custom HTTP client
func (c *Inbound) WithHTTPClient(client *http.Client) *Inbound {
	c.httpClient = client
	return c
}

// request makes an authenticated request to the API with { data, error } response pattern
func (c *Inbound) request(ctx context.Context, method, endpoint string, body any, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + endpoint

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Set custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.httpClient.Do(req)
}

// makeRequest is a generic helper that handles the complete request cycle
func makeRequest[T any](c *Inbound, ctx context.Context, method, endpoint string, body any, headers map[string]string) (*ApiResponse[T], error) {
	resp, err := c.request(ctx, method, endpoint, body, headers)
	if err != nil {
		return &ApiResponse[T]{Error: err.Error()}, nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ApiResponse[T]{Error: "Failed to read response body"}, nil
	}

	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(respBody, &errorResp); err == nil && errorResp.Error != "" {
			return &ApiResponse[T]{Error: errorResp.Error}, nil
		}
		return &ApiResponse[T]{Error: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)}, nil
	}

	var result T
	if err := json.Unmarshal(respBody, &result); err != nil {
		return &ApiResponse[T]{Error: "Failed to parse response"}, nil
	}

	return &ApiResponse[T]{Data: &result}, nil
}

// buildQueryString builds a query string from a struct
func buildQueryString(params any) string {
	values := url.Values{}

	if params == nil {
		return ""
	}

	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get JSON tag
		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		// Parse JSON tag
		tagParts := strings.Split(tag, ",")
		key := tagParts[0]

		// Check for omitempty
		omitempty := slices.Contains(tagParts[1:], "omitempty")

		// Handle different field types
		switch field.Kind() {
		case reflect.Ptr:
			if field.IsNil() {
				continue
			}
			field = field.Elem()
			fallthrough
		case reflect.String:
			val := field.String()
			if omitempty && val == "" {
				continue
			}
			values.Add(key, val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val := field.Int()
			if omitempty && val == 0 {
				continue
			}
			values.Add(key, strconv.FormatInt(val, 10))
		case reflect.Bool:
			val := field.Bool()
			if omitempty && !val {
				continue
			}
			values.Add(key, strconv.FormatBool(val))
		}
	}

	if len(values) == 0 {
		return ""
	}
	return "?" + values.Encode()
}

// MailService handles mail operations (inbound emails)
type MailService struct {
	client *Inbound
}

// NewMailService creates a new mail service
func NewMailService(client *Inbound) *MailService {
	return &MailService{client: client}
}

// List retrieves all emails in the mailbox
//
// API Reference: https://docs.inbound.new/api-reference/mail/list-emails
func (s *MailService) List(ctx context.Context, params *GetMailRequest) (*ApiResponse[GetMailResponse], error) {
	endpoint := "/mail" + buildQueryString(params)
	return makeRequest[GetMailResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Get retrieves a specific email by ID
//
// API Reference: https://docs.inbound.new/api-reference/mail/get-email
func (s *MailService) Get(ctx context.Context, id string) (*ApiResponse[GetMailByIDResponse], error) {
	endpoint := fmt.Sprintf("/mail/%s", id)
	return makeRequest[GetMailByIDResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Thread retrieves email thread/conversation by email ID
func (s *MailService) Thread(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/mail/%s/thread", id)
	return makeRequest[any](s.client, ctx, "GET", endpoint, nil, nil)
}

// MarkRead marks an email as read
func (s *MailService) MarkRead(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/mail/%s", id)
	body := map[string]bool{"isRead": true}
	return makeRequest[any](s.client, ctx, "PATCH", endpoint, body, nil)
}

// MarkUnread marks an email as unread
func (s *MailService) MarkUnread(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/mail/%s", id)
	body := map[string]bool{"isRead": false}
	return makeRequest[any](s.client, ctx, "PATCH", endpoint, body, nil)
}

// Archive archives an email
func (s *MailService) Archive(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/mail/%s", id)
	body := map[string]bool{"isArchived": true}
	return makeRequest[any](s.client, ctx, "PATCH", endpoint, body, nil)
}

// Unarchive unarchives an email
func (s *MailService) Unarchive(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/mail/%s", id)
	body := map[string]bool{"isArchived": false}
	return makeRequest[any](s.client, ctx, "PATCH", endpoint, body, nil)
}

// Reply replies to an email
func (s *MailService) Reply(ctx context.Context, params *PostMailRequest) (*ApiResponse[PostMailResponse], error) {
	return makeRequest[PostMailResponse](s.client, ctx, "POST", "/mail", params, nil)
}

// Bulk performs bulk operations on multiple emails
func (s *MailService) Bulk(ctx context.Context, emailIDs []string, updates map[string]any) (*ApiResponse[any], error) {
	body := map[string]any{
		"emailIds": emailIDs,
		"updates":  updates,
	}
	return makeRequest[any](s.client, ctx, "POST", "/mail/bulk", body, nil)
}

// EmailService handles email operations (sending emails)
type EmailService struct {
	client  *Inbound
	Address *EmailAddressService
}

// NewEmailService creates a new email service
func NewEmailService(client *Inbound) *EmailService {
	return &EmailService{
		client:  client,
		Address: NewEmailAddressService(client),
	}
}

// Send sends an email with optional attachments and idempotency options
// 
// This method supports both immediate sending and scheduled delivery.
// If params.ScheduledAt is set, the email will be scheduled for future delivery.
//
// API Reference: https://docs.inbound.new/api-reference/emails/send-email
func (s *EmailService) Send(ctx context.Context, params *PostEmailsRequest, options *IdempotencyOptions) (*ApiResponse[PostEmailsResponse], error) {
	var endpoint string
	if params.ScheduledAt != nil {
		endpoint = "/emails/schedule"
	} else {
		endpoint = "/emails"
	}

	headers := make(map[string]string)
	if options != nil && options.IdempotencyKey != "" {
		headers["Idempotency-Key"] = options.IdempotencyKey
	}

	return makeRequest[PostEmailsResponse](s.client, ctx, "POST", endpoint, params, headers)
}

// Get retrieves a sent email by ID
//
// API Reference: https://docs.inbound.new/api-reference/emails/get-email
func (s *EmailService) Get(ctx context.Context, id string) (*ApiResponse[GetEmailByIDResponse], error) {
	endpoint := fmt.Sprintf("/emails/%s", id)
	return makeRequest[GetEmailByIDResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Reply replies to an email by ID with optional attachments
//
// API Reference: https://docs.inbound.new/api-reference/emails/reply-to-email
func (s *EmailService) Reply(ctx context.Context, id string, params *PostEmailReplyRequest, options *IdempotencyOptions) (*ApiResponse[PostEmailReplyResponse], error) {
	endpoint := fmt.Sprintf("/emails/%s/reply", id)

	headers := make(map[string]string)
	if options != nil && options.IdempotencyKey != "" {
		headers["Idempotency-Key"] = options.IdempotencyKey
	}

	return makeRequest[PostEmailReplyResponse](s.client, ctx, "POST", endpoint, params, headers)
}

// Schedule schedules an email to be sent at a future time
// 
// Supports both ISO 8601 dates and natural language (e.g., "in 1 hour", "tomorrow at 9am").
//
// API Reference: https://docs.inbound.new/api-reference/emails/schedule-email
func (s *EmailService) Schedule(ctx context.Context, params *PostScheduleEmailRequest, options *IdempotencyOptions) (*ApiResponse[PostScheduleEmailResponse], error) {
	headers := make(map[string]string)
	if options != nil && options.IdempotencyKey != "" {
		headers["Idempotency-Key"] = options.IdempotencyKey
	}

	return makeRequest[PostScheduleEmailResponse](s.client, ctx, "POST", "/emails/schedule", params, headers)
}

// ListScheduled lists scheduled emails with filtering and pagination
//
// API Reference: https://docs.inbound.new/api-reference/emails/list-scheduled-emails
func (s *EmailService) ListScheduled(ctx context.Context, params *GetScheduledEmailsRequest) (*ApiResponse[GetScheduledEmailsResponse], error) {
	endpoint := "/emails/schedule" + buildQueryString(params)
	return makeRequest[GetScheduledEmailsResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// GetScheduled gets details of a specific scheduled email
func (s *EmailService) GetScheduled(ctx context.Context, id string) (*ApiResponse[GetScheduledEmailResponse], error) {
	endpoint := fmt.Sprintf("/emails/schedule/%s", id)
	return makeRequest[GetScheduledEmailResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Cancel cancels a scheduled email (only works if status is 'scheduled')
func (s *EmailService) Cancel(ctx context.Context, id string) (*ApiResponse[DeleteScheduledEmailResponse], error) {
	endpoint := fmt.Sprintf("/emails/schedule/%s", id)
	return makeRequest[DeleteScheduledEmailResponse](s.client, ctx, "DELETE", endpoint, nil, nil)
}

// EmailAddressService handles email address management
type EmailAddressService struct {
	client *Inbound
}

// NewEmailAddressService creates a new email address service
func NewEmailAddressService(client *Inbound) *EmailAddressService {
	return &EmailAddressService{client: client}
}

// Create creates a new email address
//
// API Reference: https://docs.inbound.new/api-reference/email-addresses/create-email-address
func (s *EmailAddressService) Create(ctx context.Context, params *PostEmailAddressesRequest) (*ApiResponse[PostEmailAddressesResponse], error) {
	return makeRequest[PostEmailAddressesResponse](s.client, ctx, "POST", "/email-addresses", params, nil)
}

// List lists all email addresses
//
// API Reference: https://docs.inbound.new/api-reference/email-addresses/list-email-addresses
func (s *EmailAddressService) List(ctx context.Context, params *GetEmailAddressesRequest) (*ApiResponse[GetEmailAddressesResponse], error) {
	endpoint := "/email-addresses" + buildQueryString(params)
	return makeRequest[GetEmailAddressesResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Get gets a specific email address by ID
//
// API Reference: https://docs.inbound.new/api-reference/email-addresses/get-email-address
func (s *EmailAddressService) Get(ctx context.Context, id string) (*ApiResponse[GetEmailAddressByIDResponse], error) {
	endpoint := fmt.Sprintf("/email-addresses/%s", id)
	return makeRequest[GetEmailAddressByIDResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Update updates an email address
//
// API Reference: https://docs.inbound.new/api-reference/email-addresses/update-email-address
func (s *EmailAddressService) Update(ctx context.Context, id string, params *PutEmailAddressByIDRequest) (*ApiResponse[PutEmailAddressByIDResponse], error) {
	endpoint := fmt.Sprintf("/email-addresses/%s", id)
	return makeRequest[PutEmailAddressByIDResponse](s.client, ctx, "PUT", endpoint, params, nil)
}

// Delete deletes an email address
//
// API Reference: https://docs.inbound.new/api-reference/email-addresses/delete-email-address
func (s *EmailAddressService) Delete(ctx context.Context, id string) (*ApiResponse[DeleteEmailAddressByIDResponse], error) {
	endpoint := fmt.Sprintf("/email-addresses/%s", id)
	return makeRequest[DeleteEmailAddressByIDResponse](s.client, ctx, "DELETE", endpoint, nil, nil)
}

// DomainService handles domain management
type DomainService struct {
	client *Inbound
}

// NewDomainService creates a new domain service
func NewDomainService(client *Inbound) *DomainService {
	return &DomainService{client: client}
}

// Create creates a new domain
//
// API Reference: https://docs.inbound.new/api-reference/domains/create-domain
func (s *DomainService) Create(ctx context.Context, params *PostDomainsRequest) (*ApiResponse[PostDomainsResponse], error) {
	return makeRequest[PostDomainsResponse](s.client, ctx, "POST", "/domains", params, nil)
}

// List lists all domains
//
// API Reference: https://docs.inbound.new/api-reference/domains/list-domains
func (s *DomainService) List(ctx context.Context, params *GetDomainsRequest) (*ApiResponse[GetDomainsResponse], error) {
	endpoint := "/domains" + buildQueryString(params)
	return makeRequest[GetDomainsResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Get gets a specific domain by ID
//
// API Reference: https://docs.inbound.new/api-reference/domains/get-domain
func (s *DomainService) Get(ctx context.Context, id string) (*ApiResponse[GetDomainByIDResponse], error) {
	endpoint := fmt.Sprintf("/domains/%s", id)
	return makeRequest[GetDomainByIDResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Update updates domain settings (catch-all configuration)
//
// API Reference: https://docs.inbound.new/api-reference/domains/update-domain
func (s *DomainService) Update(ctx context.Context, id string, params *PutDomainByIDRequest) (*ApiResponse[PutDomainByIDResponse], error) {
	endpoint := fmt.Sprintf("/domains/%s", id)
	return makeRequest[PutDomainByIDResponse](s.client, ctx, "PUT", endpoint, params, nil)
}

// Delete deletes a domain
//
// API Reference: https://docs.inbound.new/api-reference/domains/delete-domain
func (s *DomainService) Delete(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/domains/%s", id)
	return makeRequest[any](s.client, ctx, "DELETE", endpoint, nil, nil)
}

// Verify initiates domain verification
func (s *DomainService) Verify(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/domains/%s/auth", id)
	return makeRequest[any](s.client, ctx, "POST", endpoint, nil, nil)
}

// GetDNSRecords gets DNS records required for domain verification
//
// API Reference: https://docs.inbound.new/api-reference/domains/get-dns-records
func (s *DomainService) GetDNSRecords(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/domains/%s/dns-records", id)
	return makeRequest[any](s.client, ctx, "GET", endpoint, nil, nil)
}

// CheckStatus checks domain verification status
func (s *DomainService) CheckStatus(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/domains/%s/auth", id)
	return makeRequest[any](s.client, ctx, "PATCH", endpoint, nil, nil)
}

// EndpointService handles endpoint management
type EndpointService struct {
	client *Inbound
}

// NewEndpointService creates a new endpoint service
func NewEndpointService(client *Inbound) *EndpointService {
	return &EndpointService{client: client}
}

// Create creates a new endpoint
//
// API Reference: https://docs.inbound.new/api-reference/endpoints/create-endpoint
func (s *EndpointService) Create(ctx context.Context, params *PostEndpointsRequest) (*ApiResponse[PostEndpointsResponse], error) {
	return makeRequest[PostEndpointsResponse](s.client, ctx, "POST", "/endpoints", params, nil)
}

// List lists all endpoints
//
// API Reference: https://docs.inbound.new/api-reference/endpoints/list-endpoints
func (s *EndpointService) List(ctx context.Context, params *GetEndpointsRequest) (*ApiResponse[GetEndpointsResponse], error) {
	endpoint := "/endpoints" + buildQueryString(params)
	return makeRequest[GetEndpointsResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Get gets a specific endpoint by ID
//
// API Reference: https://docs.inbound.new/api-reference/endpoints/get-endpoint
func (s *EndpointService) Get(ctx context.Context, id string) (*ApiResponse[GetEndpointByIDResponse], error) {
	endpoint := fmt.Sprintf("/endpoints/%s", id)
	return makeRequest[GetEndpointByIDResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Update updates an endpoint
//
// API Reference: https://docs.inbound.new/api-reference/endpoints/update-endpoint
func (s *EndpointService) Update(ctx context.Context, id string, params *PutEndpointByIDRequest) (*ApiResponse[PutEndpointByIDResponse], error) {
	endpoint := fmt.Sprintf("/endpoints/%s", id)
	return makeRequest[PutEndpointByIDResponse](s.client, ctx, "PUT", endpoint, params, nil)
}

// Delete deletes an endpoint
//
// API Reference: https://docs.inbound.new/api-reference/endpoints/delete-endpoint
func (s *EndpointService) Delete(ctx context.Context, id string) (*ApiResponse[DeleteEndpointByIDResponse], error) {
	endpoint := fmt.Sprintf("/endpoints/%s", id)
	return makeRequest[DeleteEndpointByIDResponse](s.client, ctx, "DELETE", endpoint, nil, nil)
}

// Test tests endpoint connectivity
func (s *EndpointService) Test(ctx context.Context, id string) (*ApiResponse[any], error) {
	endpoint := fmt.Sprintf("/endpoints/%s/test", id)
	return makeRequest[any](s.client, ctx, "POST", endpoint, nil, nil)
}

// ThreadService handles thread management
type ThreadService struct {
	client *Inbound
}

// NewThreadService creates a new thread service
func NewThreadService(client *Inbound) *ThreadService {
	return &ThreadService{client: client}
}

// List retrieves all email threads with optional filtering
//
// API Reference: https://docs.inbound.new/api-reference/threads/list-threads
func (s *ThreadService) List(ctx context.Context, params *GetThreadsRequest) (*ApiResponse[GetThreadsResponse], error) {
	endpoint := "/threads" + buildQueryString(params)
	return makeRequest[GetThreadsResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// Get retrieves a specific thread by ID with all messages
//
// API Reference: https://docs.inbound.new/api-reference/threads/get-thread
func (s *ThreadService) Get(ctx context.Context, id string) (*ApiResponse[GetThreadByIDResponse], error) {
	endpoint := fmt.Sprintf("/threads/%s", id)
	return makeRequest[GetThreadByIDResponse](s.client, ctx, "GET", endpoint, nil, nil)
}

// PerformAction performs an action on a thread (mark as read, archive, etc.)
//
// API Reference: https://docs.inbound.new/api-reference/threads/thread-actions
func (s *ThreadService) PerformAction(ctx context.Context, id string, params *PostThreadActionsRequest) (*ApiResponse[PostThreadActionsResponse], error) {
	endpoint := fmt.Sprintf("/threads/%s/actions", id)
	return makeRequest[PostThreadActionsResponse](s.client, ctx, "POST", endpoint, params, nil)
}

// Stats retrieves statistics about all threads
//
// API Reference: https://docs.inbound.new/api-reference/threads/thread-stats
func (s *ThreadService) Stats(ctx context.Context) (*ApiResponse[GetThreadStatsResponse], error) {
	return makeRequest[GetThreadStatsResponse](s.client, ctx, "GET", "/threads/stats", nil, nil)
}

// MarkAsRead marks all messages in a thread as read
func (s *ThreadService) MarkAsRead(ctx context.Context, id string) (*ApiResponse[PostThreadActionsResponse], error) {
	return s.PerformAction(ctx, id, &PostThreadActionsRequest{Action: "mark_as_read"})
}

// MarkAsUnread marks all messages in a thread as unread
func (s *ThreadService) MarkAsUnread(ctx context.Context, id string) (*ApiResponse[PostThreadActionsResponse], error) {
	return s.PerformAction(ctx, id, &PostThreadActionsRequest{Action: "mark_as_unread"})
}

// Archive archives a thread
func (s *ThreadService) Archive(ctx context.Context, id string) (*ApiResponse[PostThreadActionsResponse], error) {
	return s.PerformAction(ctx, id, &PostThreadActionsRequest{Action: "archive"})
}

// Unarchive unarchives a thread
func (s *ThreadService) Unarchive(ctx context.Context, id string) (*ApiResponse[PostThreadActionsResponse], error) {
	return s.PerformAction(ctx, id, &PostThreadActionsRequest{Action: "unarchive"})
}

// Add service properties to the main client
func (c *Inbound) Mail() *MailService {
	return NewMailService(c)
}

func (c *Inbound) Email() *EmailService {
	return NewEmailService(c)
}

func (c *Inbound) Domain() *DomainService {
	return NewDomainService(c)
}

func (c *Inbound) Endpoint() *EndpointService {
	return NewEndpointService(c)
}

func (c *Inbound) Thread() *ThreadService {
	return NewThreadService(c)
}

// Convenience Methods

// QuickReply provides a quick text reply to an email
func (c *Inbound) QuickReply(ctx context.Context, emailID, message, from string, options *IdempotencyOptions) (*ApiResponse[PostEmailReplyResponse], error) {
	params := &PostEmailReplyRequest{
		From: from,
		Text: &message,
	}
	return c.Email().Reply(ctx, emailID, params, options)
}

// SetupDomain provides one-step domain setup with optional webhook
func (c *Inbound) SetupDomain(ctx context.Context, domain string, webhookURL *string) (*ApiResponse[any], error) {
	// First create the domain
	domainResult, err := c.Domain().Create(ctx, &PostDomainsRequest{Domain: domain})
	if err != nil {
		return &ApiResponse[any]{Error: err.Error()}, nil
	}
	if domainResult.Error != "" {
		return &ApiResponse[any]{Error: domainResult.Error}, nil
	}

	// If webhook URL provided, create an endpoint
	if webhookURL != nil && *webhookURL != "" {
		endpointResult, err := c.Endpoint().Create(ctx, &PostEndpointsRequest{
			Name: domain + " Webhook",
			Type: "webhook",
			Config: &WebhookConfig{
				URL:           *webhookURL,
				Timeout:       30000,
				RetryAttempts: 3,
			},
		})
		if err != nil {
			return &ApiResponse[any]{Error: err.Error()}, nil
		}

		result := map[string]any{
			"domain":   domainResult.Data,
			"endpoint": endpointResult.Data,
		}
		var interfaceResult any = result
		return &ApiResponse[any]{Data: &interfaceResult}, nil
	}

	// Convert domain result to any
	var domainData any = domainResult.Data
	return &ApiResponse[any]{Data: &domainData}, nil
}

// CreateForwarder creates a simple email forwarding setup
func (c *Inbound) CreateForwarder(ctx context.Context, from, to string) (*ApiResponse[PostEndpointsResponse], error) {
	params := &PostEndpointsRequest{
		Name: fmt.Sprintf("Forward %s to %s", from, to),
		Type: "email",
		Config: &EmailConfig{
			Email: to,
		},
	}
	return c.Endpoint().Create(ctx, params)
}

// ScheduleReminder creates a quick scheduled email reminder
func (c *Inbound) ScheduleReminder(ctx context.Context, to, subject, when, from string, options *IdempotencyOptions) (*ApiResponse[PostScheduleEmailResponse], error) {
	text := fmt.Sprintf("Reminder: %s", subject)
	params := &PostScheduleEmailRequest{
		From:        from,
		To:          to,
		Subject:     subject,
		Text:        &text,
		ScheduledAt: when,
	}
	return c.Email().Schedule(ctx, params, options)
}

// Helper functions for creating pointers to basic types

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Int returns a pointer to the int value passed in.
func Int(v int) *int {
	return &v
}

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}
