-- +goose Up
ALTER TABLE public.urls ADD COLUMN deleted_at TIMESTAMP;
CREATE INDEX urls_user_uuid_deleted_at_idx ON public.urls(user_uuid, deleted_at);

-- +goose Down
DROP INDEX IF EXISTS urls_user_uuid_deleted_at_idx;
ALTER TABLE public.urls DROP COLUMN deleted_at;
