package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
)

// GetDeliveryLogs returns paginated delivery logs with optional filters.
func (a *App) GetDeliveryLogs(c echo.Context) error {
	var (
		campaignID, _   = strconv.Atoi(c.QueryParam("campaign_id"))
		subscriberID, _ = strconv.Atoi(c.QueryParam("subscriber_id"))
		listID, _       = strconv.Atoi(c.QueryParam("list_id"))
		status          = c.QueryParam("status")
		email           = c.QueryParam("email")
		page, _         = strconv.Atoi(c.QueryParam("page"))
		perPage, _      = strconv.Atoi(c.QueryParam("per_page"))
	)

	// Parse date filters.
	var fromDate, toDate *time.Time
	if from := c.QueryParam("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			fromDate = &t
		}
	}
	if to := c.QueryParam("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			toDate = &t
		}
	}

	// Default pagination.
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	logs, total, err := a.core.GetDeliveryLogs(campaignID, subscriberID, listID, status, email, fromDate, toDate, offset, perPage)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{struct {
		Results []models.DeliveryLog `json:"results"`
		Total   int                  `json:"total"`
		Page    int                  `json:"page"`
		PerPage int                  `json:"per_page"`
	}{logs, total, page, perPage}})
}

// GetDeliveryLogByID returns a single delivery log by ID.
func (a *App) GetDeliveryLogByID(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, a.i18n.T("globals.messages.invalidID"))
	}

	log, err := a.core.GetDeliveryLog(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{log})
}

// ExportDeliveryLogs exports delivery logs as CSV.
func (a *App) ExportDeliveryLogs(c echo.Context) error {
	var (
		campaignID, _   = strconv.Atoi(c.QueryParam("campaign_id"))
		subscriberID, _ = strconv.Atoi(c.QueryParam("subscriber_id"))
		listID, _       = strconv.Atoi(c.QueryParam("list_id"))
		status          = c.QueryParam("status")
	)

	// Parse date filters.
	var fromDate, toDate *time.Time
	if from := c.QueryParam("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			fromDate = &t
		}
	}
	if to := c.QueryParam("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			toDate = &t
		}
	}

	logs, err := a.core.ExportDeliveryLogs(campaignID, subscriberID, listID, status, fromDate, toDate)
	if err != nil {
		return err
	}

	// Set CSV headers.
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=delivery-logs-%s.csv", time.Now().Format("2006-01-02")))

	// Write CSV.
	w := csv.NewWriter(c.Response().Writer)
	defer w.Flush()

	// Write header row.
	if err := w.Write([]string{
		"ID",
		"Campaign ID",
		"Campaign Name",
		"Subscriber ID",
		"Subscriber Name",
		"List ID",
		"List Name",
		"From Email",
		"To Email",
		"Subject",
		"Message ID",
		"SMTP Response",
		"SMTP Code",
		"Status",
		"Error",
		"Sent At",
	}); err != nil {
		return err
	}

	// Write data rows.
	for _, log := range logs {
		if err := w.Write([]string{
			strconv.FormatInt(log.ID, 10),
			strconv.FormatInt(log.CampaignID, 10),
			log.CampaignName,
			strconv.FormatInt(log.SubscriberID, 10),
			log.SubscriberName,
			strconv.FormatInt(log.ListID, 10),
			log.ListName,
			log.FromEmail,
			log.ToEmail,
			log.Subject,
			log.MessageID,
			log.SMTPResponse,
			strconv.Itoa(log.SMTPCode),
			log.Status,
			log.Error,
			log.SentAt.Format(time.RFC3339),
		}); err != nil {
			return err
		}
	}

	return nil
}
