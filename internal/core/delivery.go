package core

import (
	"net/http"
	"time"

	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
	null "gopkg.in/volatiletech/null.v6"
)

// InsertDeliveryLog inserts a new delivery log entry.
func (c *Core) InsertDeliveryLog(log models.DeliveryLog) (int64, error) {
	var id int64
	if err := c.q.InsertDeliveryLog.Get(&id,
		log.CampaignID,
		log.SubscriberID,
		log.ListID,
		log.FromEmail,
		log.ToEmail,
		log.Subject,
		log.MessageID,
		log.SMTPResponse,
		log.SMTPCode,
		log.Status,
		log.Error,
		log.SentAt,
	); err != nil {
		c.log.Printf("error inserting delivery log: %v", err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorCreating", "name", "delivery log", "error", pqErrMsg(err)))
	}

	return id, nil
}

// GetDeliveryLogs retrieves delivery logs with optional filters.
func (c *Core) GetDeliveryLogs(campaignID, subscriberID, listID int, status, email string, fromDate, toDate *time.Time, offset, limit int) ([]models.DeliveryLog, int, error) {
	var count int
	if err := c.q.GetDeliveryLogsCount.Get(&count, campaignID, subscriberID, listID, status, email, fromDate, toDate); err != nil {
		c.log.Printf("error getting delivery log count: %v", err)
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "delivery logs", "error", pqErrMsg(err)))
	}

	if count == 0 {
		return []models.DeliveryLog{}, 0, nil
	}

	var out []models.DeliveryLog
	if err := c.q.GetDeliveryLogs.Select(&out, campaignID, subscriberID, listID, status, email, fromDate, toDate, offset, limit); err != nil {
		c.log.Printf("error fetching delivery logs: %v", err)
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "delivery logs", "error", pqErrMsg(err)))
	}

	return out, count, nil
}

// GetDeliveryLog retrieves a single delivery log by ID.
func (c *Core) GetDeliveryLog(id int64) (models.DeliveryLog, error) {
	var out models.DeliveryLog
	if err := c.q.GetDeliveryLog.Get(&out, id); err != nil {
		c.log.Printf("error fetching delivery log: %v", err)
		return models.DeliveryLog{}, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "delivery log", "error", pqErrMsg(err)))
	}

	return out, nil
}

// GetDeliveryLogsByMessageID retrieves delivery logs by message ID.
func (c *Core) GetDeliveryLogsByMessageID(messageID string) ([]models.DeliveryLog, error) {
	var out []models.DeliveryLog
	if err := c.q.GetDeliveryLogsByMessageID.Select(&out, messageID); err != nil {
		c.log.Printf("error fetching delivery logs by message ID: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "delivery logs", "error", pqErrMsg(err)))
	}

	return out, nil
}

// ExportDeliveryLogs exports delivery logs for CSV export.
func (c *Core) ExportDeliveryLogs(campaignID, subscriberID, listID int, status string, fromDate, toDate *time.Time) ([]models.DeliveryLogExport, error) {
	var out []models.DeliveryLogExport
	if err := c.q.GetDeliveryLogsForExport.Select(&out, campaignID, subscriberID, listID, status, fromDate, toDate); err != nil {
		c.log.Printf("error exporting delivery logs: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "delivery logs", "error", pqErrMsg(err)))
	}

	return out, nil
}

// UpdateDeliveryLogStatus updates the status of a delivery log.
func (c *Core) UpdateDeliveryLogStatus(id int64, status, errorMsg string) error {
	var err null.String
	if errorMsg != "" {
		err = null.StringFrom(errorMsg)
	}

	if _, e := c.q.UpdateDeliveryLogStatus.Exec(id, status, err); e != nil {
		c.log.Printf("error updating delivery log status: %v", e)
		return echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorUpdating", "name", "delivery log", "error", pqErrMsg(e)))
	}

	return nil
}

// DeleteDeliveryLogsBefore deletes delivery logs older than the specified date.
func (c *Core) DeleteDeliveryLogsBefore(before time.Time) (int, error) {
	res, err := c.q.DeleteDeliveryLogsBefore.Exec(before)
	if err != nil {
		c.log.Printf("error deleting delivery logs: %v", err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorDeleting", "name", "delivery logs", "error", pqErrMsg(err)))
	}

	n, _ := res.RowsAffected()
	return int(n), nil
}
