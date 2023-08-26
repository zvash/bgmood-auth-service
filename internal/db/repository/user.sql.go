// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const getUser = `-- name: GetUser :one
SELECT id, name, email, password, avatar, verified_at, created_at, deleted_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByAccessToken = `-- name: GetUserByAccessToken :one
SELECT u.id, u.name, u.email, u.password, u.avatar, u.verified_at, u.created_at, u.deleted_at
FROM users u,
     sessions s
WHERE u.id = s.user_id
  AND s.access_token = $1
`

func (q *Queries) GetUserByAccessToken(ctx context.Context, accessToken string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByAccessToken, accessToken)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const giveRoleToUser = `-- name: GiveRoleToUser :exec
INSERT INTO user_role (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`

type GiveRoleToUserParams struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID int64     `json:"role_id"`
}

func (q *Queries) GiveRoleToUser(ctx context.Context, arg GiveRoleToUserParams) error {
	_, err := q.db.Exec(ctx, giveRoleToUser, arg.UserID, arg.RoleID)
	return err
}

const registerUser = `-- name: RegisterUser :one
INSERT INTO users (id, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, password, avatar, verified_at, created_at, deleted_at
`

type RegisterUserParams struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (q *Queries) RegisterUser(ctx context.Context, arg RegisterUserParams) (User, error) {
	row := q.db.QueryRow(ctx, registerUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const unsetUserAvatar = `-- name: UnsetUserAvatar :one
UPDATE users
SET avatar = NULL
WHERE id = $1
RETURNING id, name, email, password, avatar, verified_at, created_at, deleted_at
`

func (q *Queries) UnsetUserAvatar(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, unsetUserAvatar, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET name     = COALESCE($1::varchar, name),
    avatar   = COALESCE($2::varchar, avatar),
    password = COALESCE($3::varchar, password)
WHERE id = $4::uuid
RETURNING id, name, email, password, avatar, verified_at, created_at, deleted_at
`

type UpdateUserParams struct {
	Name     pgtype.Text `json:"name"`
	Avatar   pgtype.Text `json:"avatar"`
	Password pgtype.Text `json:"password"`
	ID       uuid.UUID   `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Name,
		arg.Avatar,
		arg.Password,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const verifyEmail = `-- name: VerifyEmail :one
UPDATE users
SET verified_at = now()
WHERE id = $1
RETURNING id, name, email, password, avatar, verified_at, created_at, deleted_at
`

func (q *Queries) VerifyEmail(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, verifyEmail, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Avatar,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}
