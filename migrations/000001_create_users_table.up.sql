-- USERS
CREATE TABLE users (
  id SERIAL PRIMARY KEY,

  first_name VARCHAR,
  last_name VARCHAR,

  email VARCHAR UNIQUE,
  phone_number VARCHAR UNIQUE,

  password VARCHAR NOT NULL,
  account_name VARCHAR,

  is_verified BOOLEAN DEFAULT FALSE,
  role VARCHAR DEFAULT 'user',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT users_email_phone_unique UNIQUE (email, phone_number)
);
