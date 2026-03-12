CREATE TABLE payment (
  id SERIAL PRIMARY KEY,

  user_id INT NOT NULL
    REFERENCES users(id),

  payment_amount DECIMAL(10,2),
  status VARCHAR,

  plan_id INT
    REFERENCES plan(id),

  created_at TIMESTAMP DEFAULT NOW()
);
