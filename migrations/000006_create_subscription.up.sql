CREATE TABLE subscription (
  id SERIAL PRIMARY KEY,

  user_id INT NOT NULL
    REFERENCES users(id),

  payment_id INT
    REFERENCES payment(id),

  start_date TIMESTAMP,
  end_date TIMESTAMP,

  is_active BOOLEAN DEFAULT TRUE,

  created_at TIMESTAMP DEFAULT NOW()
);
