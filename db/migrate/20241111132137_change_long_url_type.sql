-- +goose Up
ALTER TABLE public.urls ALTER COLUMN long_url TYPE character varying(2048);

-- +goose Down
ALTER TABLE public.urls ALTER COLUMN long_url TYPE TEXT;
