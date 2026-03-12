-- ENUMS
CREATE TYPE campaign_status AS ENUM ('PENDING', 'RUNNING', 'PAUSED', 'COMPLETED');
CREATE TYPE outbox_status AS ENUM ('PENDING', 'PICKED_UP', 'PUBLISHED', 'FAILED');

-- CAMPAIGNS
CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    subject VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    status campaign_status DEFAULT 'PENDING',
    total_emails INTEGER DEFAULT 0,
    sent_emails INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- TRANSACTIONAL OUTBOX
CREATE TABLE outbox (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaigns(id) ON DELETE CASCADE,
    recipient VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    status outbox_status DEFAULT 'PENDING',
    retry_count INTEGER DEFAULT 0,
    last_error TEXT,
    message_id UUID DEFAULT gen_random_uuid(),
    failed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- EMAIL LOGS (Audit Trail)
CREATE TABLE email_logs (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaigns(id),
    recipient VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL, -- 'SUCCESS', 'BOUNCED', 'REJECTED'
    error_message TEXT,
    attempted_at TIMESTAMP DEFAULT NOW()
);

-- INDEXES
CREATE INDEX idx_outbox_status ON outbox(status);
CREATE INDEX idx_outbox_campaign ON outbox(campaign_id);
CREATE INDEX idx_email_logs_recipient ON email_logs(recipient);
