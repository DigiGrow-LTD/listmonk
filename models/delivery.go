package models

import (
	"time"

	null "gopkg.in/volatiletech/null.v6"
)

const (
	// DeliveryStatusSent indicates the message was successfully sent to the mail server.
	DeliveryStatusSent = "sent"
	// DeliveryStatusFailed indicates the message failed to send.
	DeliveryStatusFailed = "failed"
	// DeliveryStatusBounced indicates the message bounced after sending.
	DeliveryStatusBounced = "bounced"

	// Consent types for tracking how subscribers were added.
	ConsentTypeExplicitOptin      = "explicit_optin"
	ConsentTypeLegitimateInterest = "legitimate_interest"
	ConsentTypeContractual        = "contractual"
	ConsentTypeImported           = "imported"
)

// DeliveryLog represents a delivery confirmation record for proving email delivery.
type DeliveryLog struct {
	ID           int64     `db:"id" json:"id"`
	CampaignID   null.Int  `db:"campaign_id" json:"campaign_id"`
	SubscriberID null.Int  `db:"subscriber_id" json:"subscriber_id"`
	ListID       null.Int  `db:"list_id" json:"list_id"`

	// Email details.
	FromEmail string `db:"from_email" json:"from_email"`
	ToEmail   string `db:"to_email" json:"to_email"`
	Subject   string `db:"subject" json:"subject"`

	// SMTP response data.
	MessageID    string `db:"message_id" json:"message_id"`
	SMTPResponse string `db:"smtp_response" json:"smtp_response"`
	SMTPCode     int    `db:"smtp_code" json:"smtp_code"`

	// Status and error.
	Status string      `db:"status" json:"status"`
	Error  null.String `db:"error" json:"error,omitempty"`

	// Timestamps.
	SentAt    time.Time `db:"sent_at" json:"sent_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// Joined fields for queries.
	CampaignName   null.String `db:"campaign_name" json:"campaign_name,omitempty"`
	SubscriberName null.String `db:"subscriber_name" json:"subscriber_name,omitempty"`
	ListName       null.String `db:"list_name" json:"list_name,omitempty"`
}

// DeliveryLogExport represents a delivery log for CSV export.
type DeliveryLogExport struct {
	ID             int64     `json:"id"`
	CampaignID     int64     `json:"campaign_id"`
	CampaignName   string    `json:"campaign_name"`
	SubscriberID   int64     `json:"subscriber_id"`
	SubscriberName string    `json:"subscriber_name"`
	ListID         int64     `json:"list_id"`
	ListName       string    `json:"list_name"`
	FromEmail      string    `json:"from_email"`
	ToEmail        string    `json:"to_email"`
	Subject        string    `json:"subject"`
	MessageID      string    `json:"message_id"`
	SMTPResponse   string    `json:"smtp_response"`
	SMTPCode       int       `json:"smtp_code"`
	Status         string    `json:"status"`
	Error          string    `json:"error"`
	SentAt         time.Time `json:"sent_at"`
}

// ConsentRecord represents consent metadata for a subscription.
type ConsentRecord struct {
	Type      null.String `db:"consent_type" json:"consent_type,omitempty"`
	Source    null.String `db:"consent_source" json:"consent_source,omitempty"`
	IP        null.String `db:"consent_ip" json:"consent_ip,omitempty"`
	UserAgent null.String `db:"consent_user_agent" json:"consent_user_agent,omitempty"`
	AdminID   null.Int    `db:"consent_admin_id" json:"consent_admin_id,omitempty"`
}

// ValidConsentTypes lists the valid consent types.
var ValidConsentTypes = []string{
	ConsentTypeExplicitOptin,
	ConsentTypeLegitimateInterest,
	ConsentTypeContractual,
	ConsentTypeImported,
}

// IsValidConsentType checks if the given consent type is valid.
func IsValidConsentType(t string) bool {
	for _, v := range ValidConsentTypes {
		if v == t {
			return true
		}
	}
	return false
}

// RequiresContractualConsent checks if a list category requires contractual or legitimate_interest consent.
func RequiresContractualConsent(category string) bool {
	return category == ListCategoryLegal || category == ListCategoryTransactional
}
