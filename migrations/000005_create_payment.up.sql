CREATE TABLE payment (
  "ID" SERIAL PRIMARY KEY,

  "userID" INT NOT NULL
    REFERENCES users(id),

  "paymentAmount" DECIMAL(10,2),
  status VARCHAR,

  "planID" INT
    REFERENCES plan("ID"),

  "createdAt" TIMESTAMP DEFAULT NOW()
);
