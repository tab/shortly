-- +goose Up
CREATE TABLE urls (
  uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  long_url TEXT NOT NULL,
  short_code VARCHAR(255) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE urls;
