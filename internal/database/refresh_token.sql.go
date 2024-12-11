// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_token.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_token (
        id,
        created_at,
        updated_at,
        user_id,
        expires_at
    )
VALUES (
        $1,
        NOW(),
        NOW(),
        $2,
        $3
    )
RETURNING id, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRefreshTokenParams struct {
	ID        string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.ID, arg.UserID, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT users.id, users.created_at, users.updated_at, users.email, users.hashed_password
FROM users
    JOIN refresh_token ON users.id = refresh_token.user_id
WHERE refresh_token.id = $1
    AND revoked_at IS NULL
    AND expires_at > NOW()
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :one
UPDATE refresh_token
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, user_id, expires_at, revoked_at
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, id string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, revokeRefreshToken, id)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}