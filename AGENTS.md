# AGENTS.md - Inbound Go SDK

This file provides context and guidance for AI coding assistants working on the Inbound Go SDK.

## Project Overview

This is the official Go SDK for [Inbound](https://inbound.new) - an email infrastructure platform. The SDK provides a type-safe, idiomatic Go interface for sending and receiving emails, managing domains, configuring webhooks, and handling email attachments.

**Repository**: https://github.com/inboundemail/inbound-golang-sdk  
**Documentation**: https://docs.inbound.new/api-reference  
**Package**: github.com/inboundemail/inbound-golang-sdk

## Project Structure

```
.
├── inbound.go              # Main client and service implementations
├── types.go                # All type definitions and API request/response structs
├── webhook.go              # Webhook signature verification utilities
├── *_test.go               # Test files (one per feature area)
├── examples/               # Example usage code
├── go.mod                  # Go module definition (requires Go 1.21+)
└── README.md              # User-facing documentation
```

## Architecture

### Main Client
- `Inbound` struct is the main client (in `inbound.go`)
- Services are accessed via methods: `client.Email()`, `client.Mail()`, `client.Domain()`, etc.
- All services hold a reference to the main client for making HTTP requests

### Service Pattern
Each service follows this pattern:
```go
type ServiceName struct {
    client *Inbound
}

func NewServiceName(client *Inbound) *ServiceName {
    return &ServiceName{client: client}
}
```

Services include:
- **MailService**: Inbound email operations (list, get, mark read, reply)
- **EmailService**: Outbound email operations (send, schedule, reply)
- **DomainService**: Domain management and verification
- **EndpointService**: Webhook endpoint configuration
- **EmailAddressService**: Email address creation and management (nested under EmailService)
- **ThreadService**: Email thread/conversation management
- **AttachmentService**: Download email attachments

### Response Pattern
All API methods return `*ApiResponse[T]` which contains:
```go
type ApiResponse[T any] struct {
    Data  *T     `json:"data,omitempty"`
    Error string `json:"error,omitempty"`
}
```

This allows handling both network errors (returned as Go errors) and API errors (in the `Error` field).

## Code Conventions

### Naming
- Use `PascalCase` for exported types, methods, and functions
- Use `camelCase` for unexported variables and functions
- Prefix request types with HTTP method: `PostEmailsRequest`, `GetMailRequest`, etc.
- Suffix response types with `Response`: `PostEmailsResponse`, etc.

### Comments
- All exported functions have godoc comments
- Include `API Reference:` links to official docs where applicable
- Keep comments concise and focused on usage, not implementation

### Error Handling
- Return Go errors for network/client failures
- Return API errors in the `ApiResponse.Error` field
- Use `fmt.Errorf` with `%w` for error wrapping when appropriate
- For direct data returns (like `Download`), return standard Go errors

### Type Design
- Use pointers for optional fields: `*string`, `*int`, `*bool`
- Use `any` type for fields that can be multiple types (e.g., `To` can be string or []string)
- Add JSON tags to all struct fields
- Use `omitempty` for optional fields in requests

### Helper Functions
Provide pointer helpers at the bottom of `inbound.go`:
- `String(v string) *string`
- `Int(v int) *int`
- `Bool(v bool) *bool`

## Commonly Used Commands

### Testing
```bash
# Run all tests
go test ./...

# Run specific test
go test -v -run TestName

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Building
```bash
# Build all packages
go build ./...

# Verify module dependencies
go mod verify

# Tidy dependencies
go mod tidy
```

### Linting
```bash
# Run staticcheck (if installed)
staticcheck ./...

# Run go vet
go vet ./...
```

## Testing Guidelines

### Test File Organization
- One test file per service: `send_email_test.go`, `webhook_test.go`, etc.
- Use `httptest.NewServer` to mock API responses
- Test both success and error cases
- Verify HTTP method, headers (especially Authorization), and request body

### Test Structure
```go
func TestServiceMethod(t *testing.T) {
    tests := []struct {
        name           string
        // input params
        serverResponse interface{}
        serverStatus   int
        expectError    bool
        // assertions
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // setup mock server
            // create client
            // execute method
            // verify results
        })
    }
}
```

## Adding New Features

### Adding a New API Method
1. Add request/response types to `types.go`
2. Add method to appropriate service in `inbound.go`
3. Add godoc comment with API reference link
4. Create test file `feature_test.go` with comprehensive tests
5. Update README.md with usage example if it's a major feature

### Adding a New Service
1. Define service struct in `inbound.go`:
   ```go
   type NewService struct {
       client *Inbound
   }
   ```
2. Add constructor: `func NewNewService(client *Inbound) *NewService`
3. Add accessor method to main client: `func (c *Inbound) New() *NewService`
4. Implement service methods using `makeRequest[T]` helper
5. Create `new_service_test.go` with tests

## API Documentation References

When implementing new features, always reference the official API docs:

- **Base API**: https://docs.inbound.new/api-reference
- **Send Email**: https://docs.inbound.new/api-reference/emails/send-email
- **Schedule Email**: https://docs.inbound.new/api-reference/emails/schedule-email
- **Reply to Email**: https://docs.inbound.new/api-reference/emails/reply-to-email
- **List Mail**: https://docs.inbound.new/api-reference/mail/list-emails
- **Get Email**: https://docs.inbound.new/api-reference/mail/get-email
- **Domains**: https://docs.inbound.new/api-reference/domains
- **Endpoints**: https://docs.inbound.new/api-reference/endpoints
- **Email Addresses**: https://docs.inbound.new/api-reference/email-addresses
- **Threads**: https://docs.inbound.new/api-reference/threads
- **Attachments**: https://docs.inbound.new/api-reference/attachments/download-attachment
- **Webhook Structure**: https://docs.inbound.new/webhook

## Important Implementation Details

### Request Building
- Use `buildQueryString()` helper for GET request parameters
- Query params are built from struct fields with `json` tags
- URL path parameters should be escaped with `url.PathEscape()`

### HTTP Client
- Default timeout: 30 seconds
- Can be customized via `WithHTTPClient()`
- All requests include `Authorization: Bearer <api-key>` header
- All requests set `Content-Type: application/json` (except attachment downloads)

### Context Support
- All API methods accept `context.Context` as first parameter
- Use `http.NewRequestWithContext()` for cancellation support
- Respect context timeouts and cancellations

### Idempotency
- Supported on email send, reply, and schedule operations
- Pass `IdempotencyOptions` with `IdempotencyKey` field
- Key is sent as `Idempotency-Key` HTTP header

## Common Patterns

### Making API Requests
```go
func (s *Service) MethodName(ctx context.Context, params *RequestType) (*ApiResponse[ResponseType], error) {
    endpoint := "/path/to/endpoint" + buildQueryString(params)
    return makeRequest[ResponseType](s.client, ctx, "GET", endpoint, nil, nil)
}
```

### Handling Optional Headers
```go
headers := make(map[string]string)
if options != nil && options.IdempotencyKey != "" {
    headers["Idempotency-Key"] = options.IdempotencyKey
}
return makeRequest[T](s.client, ctx, "POST", endpoint, body, headers)
```

### URL Construction with Path Parameters
```go
endpoint := fmt.Sprintf("/resource/%s", id)
// or with URL encoding
endpoint := fmt.Sprintf("/resource/%s/%s", id, url.PathEscape(filename))
```

## Dependencies

- **Standard library only** - no external dependencies
- Go 1.21+ required (uses generics)

## License

MIT License - This is an open-source SDK.
