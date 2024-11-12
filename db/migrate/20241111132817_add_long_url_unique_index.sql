-- +goose Up
ALTER TABLE public.urls ADD CONSTRAINT urls_long_url_key UNIQUE (long_url);

-- +goose Down
ALTER TABLE public.urls DROP CONSTRAINT urls_long_url_key;
