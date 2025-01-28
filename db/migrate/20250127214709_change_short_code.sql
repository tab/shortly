-- +goose Up
ALTER TABLE urls ALTER COLUMN short_code SET DEFAULT generate_short_code();

-- +goose Down
ALTER TABLE urls ALTER COLUMN short_code DROP DEFAULT;
