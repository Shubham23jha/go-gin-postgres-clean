ALTER TABLE email_logs ADD COLUMN message_id VARCHAR(255) NOT NULL DEFAULT '';
CREATE INDEX idx_email_logs_message_id ON email_logs(message_id);
