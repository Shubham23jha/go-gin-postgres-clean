CREATE TABLE plan (
  id SERIAL PRIMARY KEY,

  plan_type VARCHAR,
  feature_json JSON,

  is_active BOOLEAN DEFAULT TRUE
);
