CREATE TABLE user_sessions (
  id SERIAL PRIMARY KEY,

  user_id INT NOT NULL,

  refresh_token VARCHAR UNIQUE NOT NULL,

  device_id VARCHAR,
  device_name VARCHAR,
  browser VARCHAR,
  ip_address VARCHAR,

  is_active BOOLEAN DEFAULT TRUE,

  expires_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT fk_user_sessions_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- Index for device tracking
CREATE INDEX idx_user_device
ON user_sessions(user_id, device_id);
