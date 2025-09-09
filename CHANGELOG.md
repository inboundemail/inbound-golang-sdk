# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
