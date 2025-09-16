# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2025-01-16

### Added
- **Webhook Support**: Complete webhook payload parsing for inbound emails
  - `WebhookPayload` type matching official Inbound documentation structure
  - Support for `email.received` webhook events with nested email data
  - `WebhookEmailData`, `WebhookParsedData`, and `WebhookCleanedContent` types
  - `WebhookAddressGroup` and `WebhookAddress` for email address handling
  - `WebhookAttachment` type for email attachments in webhooks
  - `ParseWebhookPayload()` function for parsing incoming webhook requests
  - Helper methods: `GetFromAddress()`, `GetToAddress()`, `GetHeaders()`
  - Support for both `parsedData` and `cleanedContent` from webhook payloads
  - Complex header parsing (strings, arrays, objects like DKIM signatures)
  - Comprehensive webhook parsing tests with edge case coverage

### Technical Details
- Added `webhook.go` with webhook parsing utilities
- Added `webhook_test.go` with comprehensive test coverage
- Updated type definitions to match official Inbound webhook structure
- Support for flexible date handling (string or Date object)
- Proper handling of optional fields and null values

## [0.1.0] - 2024-01-XX

### Added
- Initial release of the Inbound Go SDK
- Complete API coverage for all Inbound Email endpoints
- Support for sending emails with attachments
- Email scheduling functionality
- Inbound email management (list, get, mark read/unread, archive)
- Domain management (create, list, get, update, delete, verify)
- Endpoint management (webhooks and email forwarding)
- Email address management
- Context support for all operations
- Idempotency key support
- Comprehensive documentation with API reference links
- Helper functions for pointer creation (String, Int, Bool)
- Convenience methods (QuickReply, SetupDomain, CreateForwarder, ScheduleReminder)
- Complete TypeScript to Go type conversions
- Examples and usage documentation

### Features
- **Mail Service**: Handle inbound emails
  - List emails with filtering and pagination
  - Get specific emails by ID
  - Mark emails as read/unread
  - Archive/unarchive emails
  - Reply to emails
  - Bulk operations
  
- **Email Service**: Send emails
  - Send immediate emails
  - Schedule emails for future delivery
  - Support for attachments (both remote URLs and base64 content)
  - Rich HTML and text content
  - CC, BCC, reply-to support
  - Email tags for tracking
  - Idempotency support
  
- **Domain Service**: Manage domains
  - Add new domains
  - List and filter domains
  - Get domain details
  - Update domain settings (catch-all configuration)
  - Domain verification
  - DNS record management
  
- **Endpoint Service**: Manage webhooks and email endpoints
  - Create webhook endpoints with custom headers and retry logic
  - Create email forwarding endpoints
  - List, get, update, delete endpoints
  - Test endpoint connectivity
  
- **Email Address Service**: Manage email addresses
  - Create email addresses linked to domains
  - Configure routing to endpoints
  - List, get, update, delete email addresses
  - Receipt rule management

### Technical Details
- Built with Go 1.21+
- Uses generics for type-safe API responses
- Comprehensive error handling
- HTTP client customization support
- Query parameter building from struct tags
- JSON marshaling/unmarshaling for all API types
