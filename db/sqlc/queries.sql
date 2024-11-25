-- name: HealthCheck :one
SELECT 1;

-- name: CreateURL :one
INSERT INTO urls (uuid, long_url, short_code, user_uuid)
VALUES ($1, $2, $3, $4)
ON CONFLICT (long_url) DO UPDATE SET short_code = urls.short_code
RETURNING uuid, long_url, short_code;

-- name: GetURLByShortCode :one
SELECT uuid, long_url, short_code FROM urls WHERE short_code = $1 AND deleted_at IS NULL;

-- name: GetURLsByUserID :many
WITH counter AS (
  SELECT COUNT(*) AS total
  FROM urls
  WHERE user_uuid = $1 AND deleted_at IS NULL
)
SELECT
  u.uuid,
  u.long_url,
  u.short_code,
  counter.total
FROM urls AS u
RIGHT JOIN counter ON TRUE
WHERE u.user_uuid = $1 AND deleted_at IS NULL
ORDER BY u.created_at DESC LIMIT $2 OFFSET $3;

-- name: DeleteURLsByUserIDAndShortCodes :exec
UPDATE urls
SET deleted_at = NOW()
WHERE user_uuid = $1 AND short_code = ANY($2::varchar[]) AND deleted_at IS NULL;
