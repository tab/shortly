-- +goose Up
CREATE SEQUENCE short_code_seq START 1;

-- +goose Down
DROP SEQUENCE short_code_seq;
