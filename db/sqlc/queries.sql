-- name: HealthCheck :one
SELECT 1;

-- name: CreateURL :one
INSERT INTO urls (uuid, long_url, short_code)
VALUES ($1, $2, $3)
ON CONFLICT (long_url) DO UPDATE SET short_code = urls.short_code
RETURNING uuid, long_url, short_code;

-- name: GetURLByShortCode :one
SELECT uuid, long_url, short_code FROM urls WHERE short_code = $1;
