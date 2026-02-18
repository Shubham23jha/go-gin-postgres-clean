CREATE TABLE plan (
  "ID" SERIAL PRIMARY KEY,

  "planType" VARCHAR,
  "featureJson" JSON,

  "isActive" BOOLEAN DEFAULT TRUE
);
