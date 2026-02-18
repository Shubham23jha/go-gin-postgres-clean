CREATE TABLE "userSessions" (
  "ID" SERIAL PRIMARY KEY,

  "userID" INT NOT NULL,

  "refreshToken" VARCHAR UNIQUE NOT NULL,

  "deviceID" VARCHAR,
  "deviceName" VARCHAR,
  "browser" VARCHAR,
  "ipAddress" VARCHAR,

  "isActive" BOOLEAN DEFAULT TRUE,

  "expiresAt" TIMESTAMP,
  "createdAt" TIMESTAMP DEFAULT NOW(),

  CONSTRAINT fk_user_sessions_user
FOREIGN KEY ("userID")
REFERENCES users("ID")
ON DELETE CASCADE
);

-- Index for device tracking
CREATE INDEX idx_user_device
ON "userSessions"("userID", "deviceID");
