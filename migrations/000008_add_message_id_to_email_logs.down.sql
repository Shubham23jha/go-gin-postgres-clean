DROP INDEX IF EXISTS idx_email_logs_message_id;
ALTER TABLE email_logs DROP COLUMN IF EXISTS message_id;
