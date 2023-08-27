-- name: GetTokenByToken :one
SELECT *
FROM tokens
WHERE token = $1
LIMIT 1;

-- name: GetTokenByTokenAndType :one
SELECT *
FROM tokens
WHERE token = $1
  AND type = $2
LIMIT 1;

-- name: CreateToken :one
INSERT INTO tokens (email, token, type, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteExpiredTokens :exec
DELETE
FROM tokens
WHERE expires_at <= now();