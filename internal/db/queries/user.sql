-- name: RegisterUser :one
INSERT INTO users (id, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name     = COALESCE(sqlc.narg(name)::varchar, name),
    avatar   = COALESCE(sqlc.narg(avatar)::varchar, avatar),
    password = COALESCE(sqlc.narg(password)::varchar, password)
WHERE id = sqlc.arg(id)::uuid
RETURNING *;

-- name: UnsetUserAvatar :one
UPDATE users
SET avatar = NULL
WHERE id = $1
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByAccessToken :one
SELECT u.*
FROM users u,
     sessions s
WHERE u.id = s.user_id
  AND s.access_token = $1;

-- name: VerifyEmail :one
UPDATE users
SET verified_at = now()
WHERE id = $1
RETURNING *;