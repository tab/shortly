-- +goose Up
ALTER TABLE public.urls ADD COLUMN user_uuid UUID;
CREATE INDEX urls_user_uuid_idx ON public.urls(user_uuid);

-- +goose Down
DROP INDEX IF EXISTS urls_user_uuid_idx;
ALTER TABLE public.urls DROP COLUMN user_uuid;
