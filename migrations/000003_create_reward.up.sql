CREATE TABLE reward (
  id SERIAL PRIMARY KEY,

  user_id INT NOT NULL UNIQUE,

  coin BIGINT DEFAULT 0,

  CONSTRAINT fk_reward
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);
