package models

import (
	"github.com/lib/pq"
	null "gopkg.in/volatiletech/null.v6"
)

const (
	ListTypePrivate    = "private"
	ListTypePublic     = "public"
	ListOptinSingle    = "single"
	ListOptinDouble    = "double"
	ListStatusActive   = "active"
	ListStatusArchived = "archived"

	// List category types determine behavior for tracking and unsubscribe.
	ListCategoryMarketing     = "marketing"     // Standard marketing emails - full tracking, unsubscribe allowed.
	ListCategoryTransactional = "transactional" // No unsubscribe, no tracking pixels/link tracking.
	ListCategoryLegal         = "legal"         // No unsubscribe, delivery confirmation logging, requires explicit template.
	ListCategoryService       = "service"       // Admin-configurable unsubscribe behavior.
)

// List represents a mailing list.
type List struct {
	Base

	UUID             string         `db:"uuid" json:"uuid"`
	Name             string         `db:"name" json:"name"`
	Type             string         `db:"type" json:"type"`
	Optin            string         `db:"optin" json:"optin"`
	Status           string         `db:"status" json:"status"`
	Tags             pq.StringArray `db:"tags" json:"tags"`
	Description      string         `db:"description" json:"description"`
	SubscriberCount  int            `db:"subscriber_count" json:"subscriber_count"`
	SubscriberCounts StringIntMap   `db:"subscriber_statuses" json:"subscriber_statuses"`
	SubscriberID     int            `db:"subscriber_id" json:"-"`

	// Category defines the list behavior (marketing, transactional, legal, service).
	Category      string `db:"category" json:"category"`
	NoUnsubscribe bool   `db:"no_unsubscribe" json:"no_unsubscribe"` // For service lists: admin-configurable.
	NoTracking    bool   `db:"no_tracking" json:"no_tracking"`       // Disable tracking pixels and link tracking.

	// This is only relevant when querying the lists of a subscriber.
	SubscriptionStatus    string    `db:"subscription_status" json:"subscription_status,omitempty"`
	SubscriptionCreatedAt null.Time `db:"subscription_created_at" json:"subscription_created_at,omitempty"`
	SubscriptionUpdatedAt null.Time `db:"subscription_updated_at" json:"subscription_updated_at,omitempty"`

	// Pseudofield for getting the total number of subscribers
	// in searches and queries.
	Total int `db:"total" json:"-"`
}

// AllowsUnsubscribe returns true if the list allows subscribers to unsubscribe themselves.
func (l *List) AllowsUnsubscribe() bool {
	switch l.Category {
	case ListCategoryMarketing:
		return true
	case ListCategoryTransactional, ListCategoryLegal:
		return false
	case ListCategoryService:
		return !l.NoUnsubscribe
	default:
		return true
	}
}

// RequiresDeliveryLogging returns true if the list requires delivery confirmation logging.
func (l *List) RequiresDeliveryLogging() bool {
	return l.Category == ListCategoryLegal || l.Category == ListCategoryTransactional
}

// AllowsTracking returns true if the list allows tracking pixels and link tracking.
func (l *List) AllowsTracking() bool {
	if l.NoTracking {
		return false
	}
	// Transactional and legal lists don't allow tracking by default.
	return l.Category != ListCategoryTransactional && l.Category != ListCategoryLegal
}
