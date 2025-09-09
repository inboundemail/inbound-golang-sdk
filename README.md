# Inbound Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/R44VC0RP/inbound-go.svg)](https://pkg.go.dev/github.com/R44VC0RP/inbound-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/R44VC0RP/inbound-go)](https://goreportcard.com/report/github.com/R44VC0RP/inbound-go)

A Go SDK for [Inbound](https://inbound.new) - Email infrastructure made simple for Go developers.

## 🚀 Quick Start

### Installation

```bash
go get github.com/R44VC0RP/inbound-go
```

### Send your first email

```go
package main

import (
    "context"
    "fmt"
    "log"

    inbound "github.com/R44VC0RP/inbound-go"
)

func main() {
    client, err := inbound.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    
    // Send an email
    resp, err := client.Email().Send(ctx, &inbound.PostEmailsRequest{
        From:    "hello@yourdomain.com",
        To:      "recipient@example.com",
        Subject: "Hello from Go!",
        Text:    inbound.String("This email was sent using the Inbound Go SDK!"),
        HTML:    inbound.String("<h1>Hello from Go!</h1><p>This email was sent using the Inbound Go SDK!</p>"),
    }, nil)
    
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.Error != "" {
        log.Fatal(resp.Error)
    }
    
    fmt.Printf("Email sent! ID: %s\n", resp.Data.ID)
}
```

## 📧 Features

- **Send Emails**: Send transactional emails with attachments, scheduling, and rich content
- **Receive Emails**: Handle inbound emails with webhook processing  
- **Domain Management**: Add and verify your domains
- **Email Address Management**: Create and manage email addresses
- **Endpoint Management**: Configure webhook and email endpoints
- **Scheduling**: Schedule emails for future delivery
- **Attachments**: Support for file attachments and embedded images
- **Idempotency**: Built-in support for idempotent operations
- **Context Support**: All operations support Go's context for timeouts and cancellation

## 📚 Usage Examples

### Initialize the client

```go
// Basic initialization
client, err := inbound.NewClient("your-api-key")

// With custom base URL
client, err := inbound.NewClient("your-api-key", "https://custom-api-url.com")

// With custom HTTP client
httpClient := &http.Client{Timeout: 10 * time.Second}
client, err := inbound.NewClient("your-api-key")
client.WithHTTPClient(httpClient)
```

### Send emails with attachments

```go
resp, err := client.Email().Send(ctx, &inbound.PostEmailsRequest{
    From:    "sender@yourdomain.com",
    To:      []string{"user1@example.com", "user2@example.com"},
    CC:      "manager@example.com",
    Subject: "Monthly Report",
    HTML:    inbound.String("<h1>Monthly Report</h1><p>Please find the report attached.</p>"),
    Attachments: []inbound.AttachmentData{
        {
            Filename:    "report.pdf",
            Content:     inbound.String("base64-encoded-content"),
            ContentType: inbound.String("application/pdf"),
        },
        {
            Path:     inbound.String("https://example.com/logo.png"),
            Filename: "logo.png",
            ContentID: inbound.String("company-logo"), // For embedding in HTML
        },
    },
    Tags: []inbound.EmailTag{
        {Name: "category", Value: "reports"},
        {Name: "priority", Value: "high"},
    },
}, nil)
```

### Schedule emails

```go
// Schedule with natural language
resp, err := client.Email().Schedule(ctx, &inbound.PostScheduleEmailRequest{
    From:        "notifications@yourdomain.com",
    To:          "user@example.com",
    Subject:     "Reminder: Meeting Tomorrow",
    Text:        inbound.String("Don't forget about our meeting tomorrow at 2 PM!"),
    ScheduledAt: "tomorrow at 1 PM",
    Timezone:    inbound.String("America/New_York"),
}, nil)

// Schedule with ISO 8601 timestamp
resp, err := client.Email().Schedule(ctx, &inbound.PostScheduleEmailRequest{
    From:        "notifications@yourdomain.com", 
    To:          "user@example.com",
    Subject:     "Scheduled Notification",
    HTML:        inbound.String("<p>This is a scheduled email!</p>"),
    ScheduledAt: "2024-12-25T09:00:00Z",
}, nil)
```

### Manage inbound emails

```go
// List received emails
emails, err := client.Mail().List(ctx, &inbound.GetMailRequest{
    Limit:     inbound.Int(10),
    TimeRange: "7d",
    Status:    "all",
})

// Get specific email
email, err := client.Mail().Get(ctx, "email-id")

// Mark as read
_, err = client.Mail().MarkRead(ctx, "email-id")

// Reply to an email
_, err = client.Mail().Reply(ctx, &inbound.PostMailRequest{
    EmailID:  "original-email-id",
    To:       "sender@example.com",
    Subject:  "Re: Your inquiry",
    TextBody: "Thank you for your message. We'll get back to you soon!",
})
```

### Domain management

```go
// Add a domain
domain, err := client.Domain().Create(ctx, &inbound.PostDomainsRequest{
    Domain: "yourdomain.com",
})

// List domains
domains, err := client.Domain().List(ctx, &inbound.GetDomainsRequest{
    Limit:  inbound.Int(20),
    Status: "verified",
})

// Get DNS records for verification
records, err := client.Domain().GetDNSRecords(ctx, "domain-id")

// Verify domain
_, err = client.Domain().Verify(ctx, "domain-id")
```

### Email address management

```go
// Create an email address
emailAddr, err := client.Email().Address.Create(ctx, &inbound.PostEmailAddressesRequest{
    Address:    "support@yourdomain.com",
    DomainID:   "domain-id",
    EndpointID: inbound.String("webhook-endpoint-id"),
    IsActive:   inbound.Bool(true),
})

// List email addresses
addresses, err := client.Email().Address.List(ctx, &inbound.GetEmailAddressesRequest{
    DomainID: "domain-id",
    IsActive: "true",
})
```

### Webhook endpoints

```go
// Create a webhook endpoint
endpoint, err := client.Endpoint().Create(ctx, &inbound.PostEndpointsRequest{
    Name: "Main Webhook",
    Type: "webhook",
    Config: &inbound.WebhookConfig{
        URL:           "https://yourdomain.com/webhook/inbound",
        Timeout:       30000,
        RetryAttempts: 3,
        Headers: map[string]string{
            "X-Custom-Header": "value",
        },
    },
})

// Test endpoint connectivity  
_, err = client.Endpoint().Test(ctx, "endpoint-id")
```

### Convenience methods

```go
// Quick text reply
_, err = client.QuickReply(ctx, "email-id", "Thanks for your message!", "support@yourdomain.com", nil)

// One-step domain setup with webhook
webhookURL := "https://yourdomain.com/webhook"
_, err = client.SetupDomain(ctx, "yourdomain.com", &webhookURL)

// Create email forwarder
_, err = client.CreateForwarder(ctx, "info@yourdomain.com", "support@yourdomain.com")

// Schedule a reminder
_, err = client.ScheduleReminder(ctx, "user@example.com", "Meeting Reminder", "in 1 hour", "notifications@yourdomain.com", nil)
```

## 🔧 Error Handling

All methods return an `ApiResponse[T]` struct with either data or an error:

```go
resp, err := client.Email().Send(ctx, emailParams, nil)
if err != nil {
    // Handle network/client errors
    log.Fatal(err)
}

if resp.Error != "" {
    // Handle API errors
    log.Printf("API Error: %s", resp.Error)
    return
}

// Use the data
fmt.Printf("Email sent with ID: %s\n", resp.Data.ID)
```

## 🌐 API Reference

All methods are thoroughly documented with links to the official API documentation:

- **Mail Service**: [Inbound Emails API](https://docs.inbound.new/api-reference/mail)
- **Email Service**: [Send Emails API](https://docs.inbound.new/api-reference/emails) 
- **Domain Service**: [Domains API](https://docs.inbound.new/api-reference/domains)
- **Endpoint Service**: [Endpoints API](https://docs.inbound.new/api-reference/endpoints)
- **Email Address Service**: [Email Addresses API](https://docs.inbound.new/api-reference/email-addresses)

For complete API documentation, visit [docs.inbound.new](https://docs.inbound.new/api-reference).

## ⚡ Advanced Features

### Custom HTTP Client

```go
client, err := inbound.NewClient("your-api-key")
if err != nil {
    log.Fatal(err)
}

// Use custom HTTP client with timeout and retry logic
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    90 * time.Second,
        DisableCompression: true,
    },
}

client.WithHTTPClient(httpClient)
```

### Idempotency

```go
// Use idempotency key to prevent duplicate sends
options := &inbound.IdempotencyOptions{
    IdempotencyKey: "unique-key-12345",
}

resp, err := client.Email().Send(ctx, emailParams, options)
```

### Context with timeout

```go
// Set timeout for operations
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := client.Email().Send(ctx, emailParams, nil)
```

## 🛠 Development

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

### Running examples

```bash
cd examples
go run send-email/main.go
```

## 🤝 Contributing

We welcome contributions! Please see the [contributing guidelines](https://github.com/R44VC0RP/inbound/blob/main/CONTRIBUTING.md) in the main repository.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Commit your changes (`git commit -am 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Main Inbound Repository](https://github.com/R44VC0RP/inbound)
- [Inbound Website](https://inbound.new)
- [API Documentation](https://docs.inbound.new)
- [Go Package Documentation](https://pkg.go.dev/github.com/R44VC0RP/inbound-go)

---

**Stop juggling email providers. Start building.**
