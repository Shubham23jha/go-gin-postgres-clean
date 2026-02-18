-- USERS
CREATE TABLE users (
  "ID" SERIAL PRIMARY KEY,

  "firstName" VARCHAR,
  "lastName" VARCHAR,

  "email" VARCHAR UNIQUE,
  "phoneNumber" VARCHAR UNIQUE,

  "password" VARCHAR NOT NULL,
  "accountName" VARCHAR,

  "isVerified" BOOLEAN DEFAULT FALSE,
  "role" VARCHAR DEFAULT 'user',
  "createdAt" TIMESTAMP DEFAULT NOW(),
  "updatedAt" TIMESTAMP DEFAULT NOW(),

  CONSTRAINT users_email_phone_unique UNIQUE ("email", "phoneNumber")
);


