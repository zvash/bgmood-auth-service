-- name: GetTokenByToken :one
SELECT *
FROM tokens
WHERE token = $1
LIMIT 1;

-- name: CreateToken :one
INSERT INTO tokens (email, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteExpiredTokens :exec
DELETE
FROM tokens
WHERE expires_at <= now();