package migrations

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
)

// V7_0_0 adds list categories (marketing, transactional, legal, service),
// delivery confirmation logging, and consent tracking.
func V7_0_0(db *sqlx.DB, fs stuffbin.FileSystem, ko *koanf.Koanf, lo *log.Logger) error {
	// Create list_category enum for the new list category types.
	_, err := db.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'list_category') THEN
				CREATE TYPE list_category AS ENUM ('marketing', 'transactional', 'legal', 'service');
			END IF;
		END $$;

		-- Add list_category column to lists table with default 'marketing'.
		ALTER TABLE lists ADD COLUMN IF NOT EXISTS category list_category NOT NULL DEFAULT 'marketing';

		-- Add no_unsubscribe option for service lists to allow admin-configurable unsubscribe behavior.
		ALTER TABLE lists ADD COLUMN IF NOT EXISTS no_unsubscribe BOOLEAN NOT NULL DEFAULT FALSE;

		-- Add no_tracking option to disable tracking pixels and link tracking.
		ALTER TABLE lists ADD COLUMN IF NOT EXISTS no_tracking BOOLEAN NOT NULL DEFAULT FALSE;

		CREATE INDEX IF NOT EXISTS idx_lists_category ON lists(category);
	`)
	if err != nil {
		return err
	}

	// Create consent_type enum for tracking how subscribers were added to lists.
	_, err = db.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'consent_type') THEN
				CREATE TYPE consent_type AS ENUM ('explicit_optin', 'legitimate_interest', 'contractual', 'imported');
			END IF;
		END $$;

		-- Add consent tracking fields to subscriber_lists table.
		ALTER TABLE subscriber_lists ADD COLUMN IF NOT EXISTS consent_type consent_type NULL;
		ALTER TABLE subscriber_lists ADD COLUMN IF NOT EXISTS consent_source TEXT NULL;
		ALTER TABLE subscriber_lists ADD COLUMN IF NOT EXISTS consent_ip TEXT NULL;
		ALTER TABLE subscriber_lists ADD COLUMN IF NOT EXISTS consent_user_agent TEXT NULL;
		ALTER TABLE subscriber_lists ADD COLUMN IF NOT EXISTS consent_admin_id INTEGER NULL REFERENCES users(id) ON DELETE SET NULL;

		CREATE INDEX IF NOT EXISTS idx_sub_lists_consent_type ON subscriber_lists(consent_type);
	`)
	if err != nil {
		return err
	}

	// Create delivery_logs table for delivery confirmation logging.
	// This stores SMTP response and message ID for each send to prove delivery.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS delivery_logs (
			id               BIGSERIAL PRIMARY KEY,
			campaign_id      INTEGER NULL REFERENCES campaigns(id) ON DELETE SET NULL ON UPDATE CASCADE,
			subscriber_id    INTEGER NULL REFERENCES subscribers(id) ON DELETE SET NULL ON UPDATE CASCADE,
			list_id          INTEGER NULL REFERENCES lists(id) ON DELETE SET NULL ON UPDATE CASCADE,

			-- Email details.
			from_email       TEXT NOT NULL,
			to_email         TEXT NOT NULL,
			subject          TEXT NOT NULL,

			-- SMTP response data for proving delivery.
			message_id       TEXT NOT NULL DEFAULT '',
			smtp_response    TEXT NOT NULL DEFAULT '',
			smtp_code        INTEGER NOT NULL DEFAULT 0,

			-- Status: sent, failed, bounced.
			status           TEXT NOT NULL DEFAULT 'sent',
			error            TEXT NULL,

			-- Timestamps.
			sent_at          TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_delivery_logs_campaign_id ON delivery_logs(campaign_id);
		CREATE INDEX IF NOT EXISTS idx_delivery_logs_subscriber_id ON delivery_logs(subscriber_id);
		CREATE INDEX IF NOT EXISTS idx_delivery_logs_list_id ON delivery_logs(list_id);
		CREATE INDEX IF NOT EXISTS idx_delivery_logs_message_id ON delivery_logs(message_id);
		CREATE INDEX IF NOT EXISTS idx_delivery_logs_to_email ON delivery_logs(to_email);
		CREATE INDEX IF NOT EXISTS idx_delivery_logs_sent_at ON delivery_logs(sent_at);
		CREATE INDEX IF NOT EXISTS idx_delivery_logs_status ON delivery_logs(status);
	`)
	if err != nil {
		return err
	}

	lo.Println("v7.0.0: added list categories, consent tracking, and delivery logs")

	return nil
}
