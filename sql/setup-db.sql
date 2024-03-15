CREATE TABLE levels(
  id SERIAL PRIMARY KEY,
  name VARCHAR(128) UNIQUE,
  data bytea,
  ts tsvector GENERATED ALWAYS AS (to_tsvector('english', name)) STORED
);
CREATE INDEX ts_idx ON levels USING GIN (ts);
