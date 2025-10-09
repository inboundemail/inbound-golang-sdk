package inboundgo

import "time"

// Base configuration
type InboundEmailConfig struct {
	ApiKey  string `json:"apiKey"`
	BaseUrl string `json:"baseUrl,omitempty"`
}

// Standard response pattern - { data, error }
type ApiResponse[T any] struct {
	Data  *T     `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

// Pagination interface
type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Total   int  `json:"total"`
	HasMore bool `json:"hasMore,omitempty"`
}

// Idempotency options for email sending
type IdempotencyOptions struct {
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
}

// Mail API Types
type EmailItem struct {
	ID              string     `json:"id"`
	EmailID         string     `json:"emailId"`
	MessageID       *string    `json:"messageId"`
	Subject         string     `json:"subject"`
	From            string     `json:"from"`
	FromName        *string    `json:"fromName"`
	Recipient       string     `json:"recipient"`
	Preview         string     `json:"preview"`
	ReceivedAt      time.Time  `json:"receivedAt"`
	IsRead          bool       `json:"isRead"`
	ReadAt          *time.Time `json:"readAt"`
	IsArchived      bool       `json:"isArchived"`
	ArchivedAt      *time.Time `json:"archivedAt"`
	HasAttachments  bool       `json:"hasAttachments"`
	AttachmentCount int        `json:"attachmentCount"`
	ParseSuccess    *bool      `json:"parseSuccess"`
	ParseError      *string    `json:"parseError"`
	CreatedAt       time.Time  `json:"createdAt"`
}

type GetMailRequest struct {
	Limit           *int   `json:"limit,omitempty"`
	Offset          *int   `json:"offset,omitempty"`
	Search          string `json:"search,omitempty"`
	Status          string `json:"status,omitempty"` // 'all' | 'processed' | 'failed'
	Domain          string `json:"domain,omitempty"`
	TimeRange       string `json:"timeRange,omitempty"` // '24h' | '7d' | '30d' | '90d'
	IncludeArchived *bool  `json:"includeArchived,omitempty"`
	EmailAddress    string `json:"emailAddress,omitempty"`
	EmailID         string `json:"emailId,omitempty"`
}

type GetMailResponse struct {
	Emails     []EmailItem `json:"emails"`
	Pagination Pagination  `json:"pagination"`
}

type PostMailRequest struct {
	EmailID  string  `json:"emailId"`
	To       string  `json:"to"`
	Subject  string  `json:"subject"`
	TextBody string  `json:"textBody"`
	HTMLBody *string `json:"htmlBody,omitempty"`
}

type PostMailResponse struct {
	Message string `json:"message"`
}

type GetMailByIDResponse struct {
	ID          string    `json:"id"`
	EmailID     string    `json:"emailId"`
	Subject     string    `json:"subject"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	TextBody    string    `json:"textBody"`
	HTMLBody    string    `json:"htmlBody"`
	ReceivedAt  time.Time `json:"receivedAt"`
	Attachments []any     `json:"attachments"`
}

// Endpoints API Types
type WebhookConfig struct {
	URL           string            `json:"url"`
	Timeout       int               `json:"timeout"`
	RetryAttempts int               `json:"retryAttempts"`
	Headers       map[string]string `json:"headers,omitempty"`
}

type EmailConfig struct {
	Email string `json:"email"`
}

type EmailGroupConfig struct {
	Emails []string `json:"emails"`
}

type DeliveryStats struct {
	Total        int     `json:"total"`
	Successful   int     `json:"successful"`
	Failed       int     `json:"failed"`
	LastDelivery *string `json:"lastDelivery"`
}

type EndpointWithStats struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`   // 'webhook' | 'email' | 'email_group'
	Config        any           `json:"config"` // WebhookConfig | EmailConfig | EmailGroupConfig
	IsActive      bool          `json:"isActive"`
	Description   *string       `json:"description"`
	UserID        string        `json:"userId"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	GroupEmails   []string      `json:"groupEmails"`
	DeliveryStats DeliveryStats `json:"deliveryStats"`
}

type GetEndpointsRequest struct {
	Limit  *int   `json:"limit,omitempty"`
	Offset *int   `json:"offset,omitempty"`
	Type   string `json:"type,omitempty"`   // 'webhook' | 'email' | 'email_group'
	Active string `json:"active,omitempty"` // 'true' | 'false'
}

type GetEndpointsResponse struct {
	Data       []EndpointWithStats `json:"data"`
	Pagination Pagination          `json:"pagination"`
}

type PostEndpointsRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // 'webhook' | 'email' | 'email_group'
	Description *string `json:"description,omitempty"`
	Config      any     `json:"config"` // WebhookConfig | EmailConfig | EmailGroupConfig
}

type PostEndpointsResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Config      any       `json:"config"`
	IsActive    bool      `json:"isActive"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GetEndpointByIDResponse struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Type             string        `json:"type"`
	Config           any           `json:"config"`
	IsActive         bool          `json:"isActive"`
	Description      *string       `json:"description"`
	DeliveryStats    DeliveryStats `json:"deliveryStats"`
	RecentDeliveries []any         `json:"recentDeliveries"`
	AssociatedEmails []any         `json:"associatedEmails"`
	CatchAllDomains  []any         `json:"catchAllDomains"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
}

type PutEndpointByIDRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"isActive,omitempty"`
	Config      any     `json:"config,omitempty"`
}

type PutEndpointByIDResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"isActive"`
	Config      any       `json:"config"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DeleteEndpointByIDResponse struct {
	Message string `json:"message"`
	Cleanup struct {
		EmailAddressesUpdated int   `json:"emailAddressesUpdated"`
		DomainsUpdated        int   `json:"domainsUpdated"`
		GroupEmailsDeleted    int   `json:"groupEmailsDeleted"`
		DeliveriesDeleted     int   `json:"deliveriesDeleted"`
		EmailAddresses        []any `json:"emailAddresses"`
		Domains               []any `json:"domains"`
	} `json:"cleanup"`
}

// Domains API Types
type CatchAllEndpoint struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive bool   `json:"isActive"`
}

type DNSRecord struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Value      string  `json:"value"`
	IsVerified bool    `json:"isVerified,omitempty"`
	Status     string  `json:"status,omitempty"`
	Error      *string `json:"error,omitempty"`
}

type VerificationCheck struct {
	DNSRecords      []DNSRecord `json:"dnsRecords,omitempty"`
	SESStatus       string      `json:"sesStatus,omitempty"`
	IsFullyVerified bool        `json:"isFullyVerified,omitempty"`
	LastChecked     string      `json:"lastChecked,omitempty"`
}

type DomainStats struct {
	TotalEmailAddresses  int  `json:"totalEmailAddresses"`
	ActiveEmailAddresses int  `json:"activeEmailAddresses"`
	HasCatchAll          bool `json:"hasCatchAll"`
}

type DomainWithStats struct {
	ID                 string             `json:"id"`
	Domain             string             `json:"domain"`
	Status             string             `json:"status"`
	CanReceiveEmails   bool               `json:"canReceiveEmails"`
	HasMXRecords       bool               `json:"hasMxRecords"`
	DomainProvider     *string            `json:"domainProvider"`
	ProviderConfidence *string            `json:"providerConfidence"`
	LastDNSCheck       *time.Time         `json:"lastDnsCheck"`
	LastSESCheck       *time.Time         `json:"lastSesCheck"`
	IsCatchAllEnabled  bool               `json:"isCatchAllEnabled"`
	CatchAllEndpointID *string            `json:"catchAllEndpointId"`
	ReceiveDMARCEmails bool               `json:"receiveDmarcEmails"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	UserID             string             `json:"userId"`
	Stats              DomainStats        `json:"stats"`
	CatchAllEndpoint   *CatchAllEndpoint  `json:"catchAllEndpoint,omitempty"`
	VerificationCheck  *VerificationCheck `json:"verificationCheck,omitempty"`
}

type GetDomainsRequest struct {
	Limit      *int   `json:"limit,omitempty"`
	Offset     *int   `json:"offset,omitempty"`
	Status     string `json:"status,omitempty"`     // 'pending' | 'verified' | 'failed'
	CanReceive string `json:"canReceive,omitempty"` // 'true' | 'false'
	Check      string `json:"check,omitempty"`      // 'true' | 'false'
}

type GetDomainsResponse struct {
	Data       []DomainWithStats `json:"data"`
	Pagination Pagination        `json:"pagination"`
	Meta       struct {
		TotalCount        int            `json:"totalCount"`
		VerifiedCount     int            `json:"verifiedCount"`
		WithCatchAllCount int            `json:"withCatchAllCount"`
		StatusBreakdown   map[string]int `json:"statusBreakdown"`
	} `json:"meta"`
}

type PostDomainsRequest struct {
	Domain string `json:"domain"`
}

type PostDomainsResponse struct {
	ID         string      `json:"id"`
	Domain     string      `json:"domain"`
	Status     string      `json:"status"`
	DNSRecords []DNSRecord `json:"dnsRecords"`
	CreatedAt  time.Time   `json:"createdAt"`
}

type GetDomainByIDResponse struct {
	ID                 string            `json:"id"`
	Domain             string            `json:"domain"`
	Status             string            `json:"status"`
	CanReceiveEmails   bool              `json:"canReceiveEmails"`
	IsCatchAllEnabled  bool              `json:"isCatchAllEnabled"`
	CatchAllEndpointID *string           `json:"catchAllEndpointId"`
	Stats              DomainStats       `json:"stats"`
	CatchAllEndpoint   *CatchAllEndpoint `json:"catchAllEndpoint,omitempty"`
	CreatedAt          time.Time         `json:"createdAt"`
	UpdatedAt          time.Time         `json:"updatedAt"`
}

type PutDomainByIDRequest struct {
	IsCatchAllEnabled  bool    `json:"isCatchAllEnabled"`
	CatchAllEndpointID *string `json:"catchAllEndpointId"`
}

type PutDomainByIDResponse struct {
	ID                 string            `json:"id"`
	Domain             string            `json:"domain"`
	IsCatchAllEnabled  bool              `json:"isCatchAllEnabled"`
	CatchAllEndpointID *string           `json:"catchAllEndpointId"`
	CatchAllEndpoint   *CatchAllEndpoint `json:"catchAllEndpoint,omitempty"`
	UpdatedAt          time.Time         `json:"updatedAt"`
}

// Email Addresses API Types
type DomainInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type RoutingInfo struct {
	Type     string  `json:"type"` // 'webhook' | 'endpoint' | 'none'
	ID       *string `json:"id"`
	Name     *string `json:"name"`
	Config   any     `json:"config,omitempty"`
	IsActive bool    `json:"isActive"`
}

type EmailAddressWithDomain struct {
	ID                      string      `json:"id"`
	Address                 string      `json:"address"`
	DomainID                string      `json:"domainId"`
	WebhookID               *string     `json:"webhookId"`
	EndpointID              *string     `json:"endpointId"`
	IsActive                bool        `json:"isActive"`
	IsReceiptRuleConfigured bool        `json:"isReceiptRuleConfigured"`
	ReceiptRuleName         *string     `json:"receiptRuleName"`
	CreatedAt               time.Time   `json:"createdAt"`
	UpdatedAt               time.Time   `json:"updatedAt"`
	UserID                  string      `json:"userId"`
	Domain                  DomainInfo  `json:"domain"`
	Routing                 RoutingInfo `json:"routing"`
}

type GetEmailAddressesRequest struct {
	Limit                   *int   `json:"limit,omitempty"`
	Offset                  *int   `json:"offset,omitempty"`
	DomainID                string `json:"domainId,omitempty"`
	IsActive                string `json:"isActive,omitempty"`                // 'true' | 'false'
	IsReceiptRuleConfigured string `json:"isReceiptRuleConfigured,omitempty"` // 'true' | 'false'
}

type GetEmailAddressesResponse struct {
	Data       []EmailAddressWithDomain `json:"data"`
	Pagination Pagination               `json:"pagination"`
}

type PostEmailAddressesRequest struct {
	Address    string  `json:"address"`
	DomainID   string  `json:"domainId"`
	EndpointID *string `json:"endpointId,omitempty"`
	WebhookID  *string `json:"webhookId,omitempty"`
	IsActive   *bool   `json:"isActive,omitempty"`
}

type PostEmailAddressesResponse struct {
	ID         string      `json:"id"`
	Address    string      `json:"address"`
	DomainID   string      `json:"domainId"`
	EndpointID *string     `json:"endpointId"`
	IsActive   bool        `json:"isActive"`
	Domain     DomainInfo  `json:"domain"`
	Routing    RoutingInfo `json:"routing"`
	CreatedAt  time.Time   `json:"createdAt"`
}

type GetEmailAddressByIDResponse struct {
	ID                      string      `json:"id"`
	Address                 string      `json:"address"`
	DomainID                string      `json:"domainId"`
	EndpointID              *string     `json:"endpointId"`
	IsActive                bool        `json:"isActive"`
	IsReceiptRuleConfigured bool        `json:"isReceiptRuleConfigured"`
	Domain                  DomainInfo  `json:"domain"`
	Routing                 RoutingInfo `json:"routing"`
	CreatedAt               time.Time   `json:"createdAt"`
	UpdatedAt               time.Time   `json:"updatedAt"`
}

type PutEmailAddressByIDRequest struct {
	IsActive   *bool   `json:"isActive,omitempty"`
	EndpointID *string `json:"endpointId,omitempty"`
	WebhookID  *string `json:"webhookId,omitempty"`
}

type PutEmailAddressByIDResponse struct {
	ID        string      `json:"id"`
	Address   string      `json:"address"`
	IsActive  bool        `json:"isActive"`
	Domain    DomainInfo  `json:"domain"`
	Routing   RoutingInfo `json:"routing"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type DeleteEmailAddressByIDResponse struct {
	Message string `json:"message"`
	Cleanup struct {
		EmailAddress   string `json:"emailAddress"`
		Domain         string `json:"domain"`
		SESRuleUpdated bool   `json:"sesRuleUpdated"`
	} `json:"cleanup"`
}

// Enhanced attachment interface supporting both remote and base64 content
type AttachmentData struct {
	Path        *string `json:"path,omitempty"`        // Remote file URL
	Content     *string `json:"content,omitempty"`     // Base64 encoded content
	Filename    string  `json:"filename"`              // Required display name
	ContentType *string `json:"contentType,omitempty"` // Optional MIME type
	ContentID   *string `json:"content_id,omitempty"`  // Content ID for embedding images in HTML (max 128 chars)
}

type EmailTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Emails API Types (for sending)
type PostEmailsRequest struct {
	From        string            `json:"from"`
	To          any               `json:"to"` // string or []string
	Subject     string            `json:"subject"`
	BCC         any               `json:"bcc,omitempty"`     // string or []string
	CC          any               `json:"cc,omitempty"`      // string or []string
	ReplyTo     any               `json:"replyTo,omitempty"` // string or []string
	HTML        *string           `json:"html,omitempty"`
	Text        *string           `json:"text,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []AttachmentData  `json:"attachments,omitempty"`
	Tags        []EmailTag        `json:"tags,omitempty"`
	ScheduledAt *string           `json:"scheduled_at,omitempty"` // Schedule email to be sent later
	Timezone    *string           `json:"timezone,omitempty"`     // User's timezone for natural language parsing
}

type PostEmailsResponse struct {
	ID          string  `json:"id"`
	MessageID   *string `json:"messageId,omitempty"`    // AWS SES Message ID
	ScheduledAt *string `json:"scheduled_at,omitempty"` // ISO 8601 timestamp
	Status      *string `json:"status,omitempty"`       // 'sent' | 'scheduled'
	Timezone    *string `json:"timezone,omitempty"`     // Timezone used for scheduling
}

type GetEmailByIDResponse struct {
	Object    string    `json:"object"`
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	CC        []string  `json:"cc"`
	BCC       []string  `json:"bcc"`
	ReplyTo   []string  `json:"reply_to"`
	Subject   string    `json:"subject"`
	Text      string    `json:"text"`
	HTML      string    `json:"html"`
	CreatedAt time.Time `json:"created_at"`
	LastEvent string    `json:"last_event"` // 'pending' | 'delivered' | 'failed'
}

// Reply API Types
type PostEmailReplyRequest struct {
	From            string            `json:"from"`
	FromName        *string           `json:"from_name,omitempty"`
	To              any               `json:"to,omitempty"`  // string or []string
	CC              any               `json:"cc,omitempty"`  // string or []string
	BCC             any               `json:"bcc,omitempty"` // string or []string
	Subject         *string           `json:"subject,omitempty"`
	Text            *string           `json:"text,omitempty"`
	HTML            *string           `json:"html,omitempty"`
	ReplyTo         any               `json:"replyTo,omitempty"` // string or []string
	Headers         map[string]string `json:"headers,omitempty"`
	Attachments     []AttachmentData  `json:"attachments,omitempty"`
	Tags            []EmailTag        `json:"tags,omitempty"`
	IncludeOriginal *bool             `json:"includeOriginal,omitempty"`
	ReplyAll        *bool             `json:"replyAll,omitempty"`
	Simple          *bool             `json:"simple,omitempty"`
}

type PostEmailReplyResponse struct {
	ID                string  `json:"id"`
	MessageID         string  `json:"messageId"`
	AWSMessageID      *string `json:"awsMessageId,omitempty"`
	RepliedToEmailID  string  `json:"repliedToEmailId"`
	RepliedToThreadID *string `json:"repliedToThreadId,omitempty"`
	IsThreadReply     bool    `json:"isThreadReply"`
}

// Email Scheduling API Types
type PostScheduleEmailRequest struct {
	From        string            `json:"from"`
	To          any               `json:"to"` // string or []string
	Subject     string            `json:"subject"`
	BCC         any               `json:"bcc,omitempty"`     // string or []string
	CC          any               `json:"cc,omitempty"`      // string or []string
	ReplyTo     any               `json:"replyTo,omitempty"` // string or []string
	HTML        *string           `json:"html,omitempty"`
	Text        *string           `json:"text,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []AttachmentData  `json:"attachments,omitempty"`
	Tags        []EmailTag        `json:"tags,omitempty"`
	ScheduledAt string            `json:"scheduled_at"`       // ISO 8601 or natural language
	Timezone    *string           `json:"timezone,omitempty"` // User's timezone for natural language parsing
}

type PostScheduleEmailResponse struct {
	ID          string `json:"id"`
	ScheduledAt string `json:"scheduled_at"` // Normalized ISO 8601 timestamp
	Status      string `json:"status"`       // 'scheduled'
	Timezone    string `json:"timezone"`
}

type GetScheduledEmailsRequest struct {
	Limit  *int   `json:"limit,omitempty"`
	Offset *int   `json:"offset,omitempty"`
	Status string `json:"status,omitempty"` // Filter by status
}

type ScheduledEmailItem struct {
	ID          string   `json:"id"`
	From        string   `json:"from"`
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	ScheduledAt string   `json:"scheduled_at"`
	Status      string   `json:"status"`
	Timezone    string   `json:"timezone"`
	CreatedAt   string   `json:"created_at"`
	Attempts    int      `json:"attempts"`
	LastError   *string  `json:"last_error,omitempty"`
}

type GetScheduledEmailsResponse struct {
	Data       []ScheduledEmailItem `json:"data"`
	Pagination Pagination           `json:"pagination"`
}

type GetScheduledEmailResponse struct {
	ID          string            `json:"id"`
	From        string            `json:"from"`
	To          []string          `json:"to"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	ReplyTo     []string          `json:"replyTo,omitempty"`
	Subject     string            `json:"subject"`
	Text        *string           `json:"text,omitempty"`
	HTML        *string           `json:"html,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []AttachmentData  `json:"attachments,omitempty"`
	Tags        []EmailTag        `json:"tags,omitempty"`
	ScheduledAt string            `json:"scheduled_at"`
	Timezone    string            `json:"timezone"`
	Status      string            `json:"status"`
	Attempts    int               `json:"attempts"`
	MaxAttempts int               `json:"max_attempts"`
	NextRetryAt *string           `json:"next_retry_at,omitempty"`
	LastError   *string           `json:"last_error,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	SentAt      *string           `json:"sent_at,omitempty"`
	SentEmailID *string           `json:"sent_email_id,omitempty"`
}

type DeleteScheduledEmailResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"` // 'cancelled'
	CancelledAt string `json:"cancelled_at"`
}

// Threads API Types
type ThreadLatestMessage struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"` // 'inbound' | 'outbound'
	Subject        *string `json:"subject"`
	FromText       string  `json:"fromText"`
	TextPreview    *string `json:"textPreview"`
	IsRead         bool    `json:"isRead"`
	HasAttachments bool    `json:"hasAttachments"`
	Date           *string `json:"date"`
}

type ThreadSummary struct {
	ID                string                `json:"id"`
	RootMessageID     string                `json:"rootMessageId"`
	NormalizedSubject *string               `json:"normalizedSubject"`
	ParticipantEmails []string              `json:"participantEmails"`
	MessageCount      int                   `json:"messageCount"`
	LastMessageAt     string                `json:"lastMessageAt"`
	CreatedAt         string                `json:"createdAt"`
	HasUnread         bool                  `json:"hasUnread"`
	IsArchived        bool                  `json:"isArchived"`
	LatestMessage     *ThreadLatestMessage  `json:"latestMessage,omitempty"`
}

type GetThreadsRequest struct {
	Limit    *int   `json:"limit,omitempty"`
	Offset   *int   `json:"offset,omitempty"`
	Search   string `json:"search,omitempty"`
	Unread   *bool  `json:"unread,omitempty"`
	Archived *bool  `json:"archived,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Address  string `json:"address,omitempty"`
}

type GetThreadsFilters struct {
	Search       *string `json:"search,omitempty"`
	UnreadOnly   *bool   `json:"unreadOnly,omitempty"`
	ArchivedOnly *bool   `json:"archivedOnly,omitempty"`
	Domain       *string `json:"domain,omitempty"`
	Address      *string `json:"address,omitempty"`
}

type GetThreadsResponse struct {
	Threads    []ThreadSummary   `json:"threads"`
	Pagination Pagination        `json:"pagination"`
	Filters    GetThreadsFilters `json:"filters"`
}

type ThreadAttachment struct {
	Filename            string `json:"filename"`
	ContentType         string `json:"contentType"`
	Size                int    `json:"size"`
	ContentID           string `json:"contentId"`
	ContentDisposition  string `json:"contentDisposition"`
}

type ThreadMessage struct {
	ID             string             `json:"id"`
	MessageID      *string            `json:"messageId"`
	Type           string             `json:"type"` // 'inbound' | 'outbound'
	ThreadPosition int                `json:"threadPosition"`
	Subject        *string            `json:"subject"`
	TextBody       *string            `json:"textBody"`
	HTMLBody       *string            `json:"htmlBody"`
	From           string             `json:"from"`
	FromName       *string            `json:"fromName"`
	FromAddress    *string            `json:"fromAddress"`
	To             []string           `json:"to"`
	CC             []string           `json:"cc"`
	BCC            []string           `json:"bcc"`
	Date           *string            `json:"date"`
	ReceivedAt     *string            `json:"receivedAt"`
	SentAt         *string            `json:"sentAt"`
	IsRead         bool               `json:"isRead"`
	ReadAt         *string            `json:"readAt"`
	HasAttachments bool               `json:"hasAttachments"`
	Attachments    []ThreadAttachment `json:"attachments"`
	InReplyTo      *string            `json:"inReplyTo"`
	References     []string           `json:"references"`
	Headers        map[string]any     `json:"headers"`
	Tags           []EmailTag         `json:"tags,omitempty"`
	Status         *string            `json:"status,omitempty"`
	FailureReason  *string            `json:"failureReason,omitempty"`
}

type ThreadMetadata struct {
	ID                string   `json:"id"`
	RootMessageID     string   `json:"rootMessageId"`
	NormalizedSubject *string  `json:"normalizedSubject"`
	ParticipantEmails []string `json:"participantEmails"`
	MessageCount      int      `json:"messageCount"`
	LastMessageAt     string   `json:"lastMessageAt"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
}

type GetThreadByIDResponse struct {
	Thread     ThreadMetadata  `json:"thread"`
	Messages   []ThreadMessage `json:"messages"`
	TotalCount int             `json:"totalCount"`
}

type PostThreadActionsRequest struct {
	Action string `json:"action"` // 'mark_as_read' | 'mark_as_unread' | 'archive' | 'unarchive'
}

type PostThreadActionsResponse struct {
	Success          bool   `json:"success"`
	Action           string `json:"action"`
	ThreadID         string `json:"threadId"`
	AffectedMessages int    `json:"affectedMessages"`
	Message          string `json:"message"`
}

type ThreadDistribution struct {
	SingleMessageThreads int `json:"singleMessageThreads"`
	ShortThreads         int `json:"shortThreads"`
	MediumThreads        int `json:"mediumThreads"`
	LongThreads          int `json:"longThreads"`
}

type ThreadRecentActivity struct {
	ThreadsToday     int `json:"threadsToday"`
	MessagesToday    int `json:"messagesToday"`
	ThreadsThisWeek  int `json:"threadsThisWeek"`
	MessagesThisWeek int `json:"messagesThisWeek"`
}

type ThreadUnreadStats struct {
	UnreadThreads  int `json:"unreadThreads"`
	UnreadMessages int `json:"unreadMessages"`
}

type MostActiveThread struct {
	ThreadID      string  `json:"threadId"`
	MessageCount  int     `json:"messageCount"`
	Subject       *string `json:"subject"`
	LastMessageAt string  `json:"lastMessageAt"`
}

type GetThreadStatsResponse struct {
	TotalThreads            int                  `json:"totalThreads"`
	TotalMessages           int                  `json:"totalMessages"`
	AverageMessagesPerThread float64             `json:"averageMessagesPerThread"`
	MostActiveThread        *MostActiveThread    `json:"mostActiveThread"`
	RecentActivity          ThreadRecentActivity `json:"recentActivity"`
	Distribution            ThreadDistribution   `json:"distribution"`
	UnreadStats             ThreadUnreadStats    `json:"unreadStats"`
}

// Webhook Payload Types - for incoming email.received webhooks
type WebhookPayload struct {
	Event     string              `json:"event"`
	Timestamp string              `json:"timestamp"`
	Email     WebhookEmailData    `json:"email"`
	Endpoint  *WebhookEndpointRef `json:"endpoint,omitempty"`
}

type WebhookEmailData struct {
	ID             string                 `json:"id"`
	MessageID      string                 `json:"messageId"`
	From           WebhookAddressGroup    `json:"from"`
	To             WebhookAddressGroup    `json:"to"`
	Recipient      string                 `json:"recipient"`
	Subject        string                 `json:"subject"`
	ReceivedAt     string                 `json:"receivedAt"`
	ParsedData     WebhookParsedData      `json:"parsedData"`
	CleanedContent *WebhookCleanedContent `json:"cleanedContent,omitempty"`
}

type WebhookAddressGroup struct {
	Text      string           `json:"text"`
	Addresses []WebhookAddress `json:"addresses"`
}

type WebhookAddress struct {
	Name    *string `json:"name"`
	Address string  `json:"address"`
}

type WebhookParsedData struct {
	MessageID   string               `json:"messageId"`
	Date        any                  `json:"date"` // Can be string or Date object
	Subject     string               `json:"subject"`
	From        WebhookAddressGroup  `json:"from"`
	To          WebhookAddressGroup  `json:"to"`
	Cc          *WebhookAddressGroup `json:"cc"`
	Bcc         *WebhookAddressGroup `json:"bcc"`
	ReplyTo     *WebhookAddressGroup `json:"replyTo"`
	InReplyTo   *string              `json:"inReplyTo,omitempty"`
	References  *string              `json:"references,omitempty"`
	TextBody    string               `json:"textBody"`
	HTMLBody    string               `json:"htmlBody"`
	Attachments []WebhookAttachment  `json:"attachments"`
	Headers     map[string]any       `json:"headers"`
	Priority    *string              `json:"priority,omitempty"`
}

type WebhookCleanedContent struct {
	HTML        string              `json:"html"`
	Text        string              `json:"text"`
	HasHTML     bool                `json:"hasHtml"`
	HasText     bool                `json:"hasText"`
	Attachments []WebhookAttachment `json:"attachments"`
	Headers     map[string]any      `json:"headers"`
}

type WebhookAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	ContentID   string `json:"contentId"`
	URL         string `json:"url"`
	DownloadUrl string `json:"downloadUrl"`
}

type WebhookEndpointRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
