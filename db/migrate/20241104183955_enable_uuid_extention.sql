-- +goose Up
CREATE EXTENSION "uuid-ossp";

-- +goose Down
DROP EXTENSION IF EXISTS "uuid-ossp";
