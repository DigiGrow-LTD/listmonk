-- delivery logs

-- name: insert-delivery-log
INSERT INTO delivery_logs (campaign_id, subscriber_id, list_id, from_email, to_email, subject, message_id, smtp_response, smtp_code, status, error, sent_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;

-- name: get-delivery-logs
-- Retrieves delivery logs with optional filters for campaign, subscriber, list, status, and date range.
SELECT dl.*,
    c.name AS campaign_name,
    s.name AS subscriber_name,
    l.name AS list_name
FROM delivery_logs dl
LEFT JOIN campaigns c ON dl.campaign_id = c.id
LEFT JOIN subscribers s ON dl.subscriber_id = s.id
LEFT JOIN lists l ON dl.list_id = l.id
WHERE
    ($1 = 0 OR dl.campaign_id = $1)
    AND ($2 = 0 OR dl.subscriber_id = $2)
    AND ($3 = 0 OR dl.list_id = $3)
    AND ($4 = '' OR dl.status = $4)
    AND ($5 = '' OR dl.to_email ILIKE '%' || $5 || '%')
    AND ($6::TIMESTAMP IS NULL OR dl.sent_at >= $6)
    AND ($7::TIMESTAMP IS NULL OR dl.sent_at <= $7)
ORDER BY dl.sent_at DESC
OFFSET $8 LIMIT (CASE WHEN $9 < 1 THEN NULL ELSE $9 END);

-- name: get-delivery-logs-count
-- Gets the count of delivery logs with the same filters.
SELECT COUNT(*) FROM delivery_logs dl
WHERE
    ($1 = 0 OR dl.campaign_id = $1)
    AND ($2 = 0 OR dl.subscriber_id = $2)
    AND ($3 = 0 OR dl.list_id = $3)
    AND ($4 = '' OR dl.status = $4)
    AND ($5 = '' OR dl.to_email ILIKE '%' || $5 || '%')
    AND ($6::TIMESTAMP IS NULL OR dl.sent_at >= $6)
    AND ($7::TIMESTAMP IS NULL OR dl.sent_at <= $7);

-- name: get-delivery-log
-- Retrieves a single delivery log by ID.
SELECT dl.*,
    c.name AS campaign_name,
    s.name AS subscriber_name,
    l.name AS list_name
FROM delivery_logs dl
LEFT JOIN campaigns c ON dl.campaign_id = c.id
LEFT JOIN subscribers s ON dl.subscriber_id = s.id
LEFT JOIN lists l ON dl.list_id = l.id
WHERE dl.id = $1;

-- name: get-delivery-logs-by-message-id
-- Retrieves delivery logs by message ID (useful for bounce correlation).
SELECT dl.*,
    c.name AS campaign_name,
    s.name AS subscriber_name,
    l.name AS list_name
FROM delivery_logs dl
LEFT JOIN campaigns c ON dl.campaign_id = c.id
LEFT JOIN subscribers s ON dl.subscriber_id = s.id
LEFT JOIN lists l ON dl.list_id = l.id
WHERE dl.message_id = $1;

-- name: get-delivery-logs-for-export
-- Retrieves delivery logs for CSV export with all details.
SELECT
    dl.id,
    COALESCE(dl.campaign_id, 0) AS campaign_id,
    COALESCE(c.name, '') AS campaign_name,
    COALESCE(dl.subscriber_id, 0) AS subscriber_id,
    COALESCE(s.name, '') AS subscriber_name,
    COALESCE(dl.list_id, 0) AS list_id,
    COALESCE(l.name, '') AS list_name,
    dl.from_email,
    dl.to_email,
    dl.subject,
    dl.message_id,
    dl.smtp_response,
    dl.smtp_code,
    dl.status,
    COALESCE(dl.error, '') AS error,
    dl.sent_at
FROM delivery_logs dl
LEFT JOIN campaigns c ON dl.campaign_id = c.id
LEFT JOIN subscribers s ON dl.subscriber_id = s.id
LEFT JOIN lists l ON dl.list_id = l.id
WHERE
    ($1 = 0 OR dl.campaign_id = $1)
    AND ($2 = 0 OR dl.subscriber_id = $2)
    AND ($3 = 0 OR dl.list_id = $3)
    AND ($4 = '' OR dl.status = $4)
    AND ($5::TIMESTAMP IS NULL OR dl.sent_at >= $5)
    AND ($6::TIMESTAMP IS NULL OR dl.sent_at <= $6)
ORDER BY dl.sent_at DESC;

-- name: update-delivery-log-status
-- Updates the status of a delivery log (e.g., when a bounce is received).
UPDATE delivery_logs SET status = $2, error = $3 WHERE id = $1;

-- name: delete-delivery-logs-before
-- Deletes delivery logs older than the specified date.
DELETE FROM delivery_logs WHERE sent_at < $1;
