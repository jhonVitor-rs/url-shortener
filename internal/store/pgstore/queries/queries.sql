-- name: GetUsers :many
SELECT
  *
FROM
  users;
-- name: GetUser :one
SELECT
  *
FROM
  users
WHERE
  id = $1;
-- name: GetUserByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = $1;
-- name: InsertUser :one
INSERT INTO
  users (NAME, email)
VALUES
  ($1, $2) RETURNING id;
-- name: UpdateUser :one
UPDATE
  users
SET
  name = $2,
  email = $3
WHERE
  id = $1 RETURNING id;
-- name: DeleteUser :exec
DELETE FROM
  users
WHERE
  id = $1;
-- name: GetShortUrlBySlug :one
SELECT
  *
FROM
  short_urls
WHERE
  slug = $1;
-- name: GetShortUrlsByUserId :many
SELECT
  *
FROM
  short_urls
WHERE
  user_id = $1;
-- name: GetShortUrlById :one
SELECT
  *
FROM
  short_urls
WHERE
  id = $1;
-- name: CreateShortUrl :one
INSERT INTO
  short_urls (user_id, slug, original_url, expires_at)
VALUES
  ($1, $2, $3, $4) RETURNING id;
-- name: UpdateShortUrl :one
UPDATE
  short_urls
SET
  slug = $2,
  original_url = $3,
  expires_at = $4
WHERE
  id = $1 RETURNING id;
-- name: IncrementAccessCount :exec
UPDATE
  short_urls
SET
  access_count = access_count + $2
WHERE
  id = $1;
-- name: DeleteShortUrl :exec
DELETE FROM
  short_urls
WHERE
  id = $1;