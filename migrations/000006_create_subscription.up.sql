CREATE TABLE subscription (
  "ID" SERIAL PRIMARY KEY,

  "userID" INT NOT NULL
    REFERENCES users("ID"),

  "paymentID" INT
    REFERENCES payment("ID"),

  "startDate" TIMESTAMP,
  "endDate" TIMESTAMP,

  "isActive" BOOLEAN DEFAULT TRUE,

  "createdAt" TIMESTAMP DEFAULT NOW()
);
