-- name: CreateSession :one
INSERT INTO sessions (id,
                      user_id,
                      access_token,
                      refresh_token,
                      user_agent,
                      client_ip,
                      is_blocked,
                      expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetSession :one
SELECT *
FROM sessions
WHERE id = $1
LIMIT 1;

-- name: TerminateOtherSessions :exec
DELETE
FROM sessions
WHERE id <> $1
  AND user_id = $2;

-- name: GetSessionByAccessToken :one
SELECT *
FROM sessions
WHERE access_token = $1;

-- name: ListActiveSessions :many
SELECT *
FROM sessions
WHERE user_id = $1
  AND expires_at > now()
  AND is_blocked = false
ORDER BY access_token = $2, expires_at DESC;