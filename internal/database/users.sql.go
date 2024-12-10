// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
	gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING id, created_at, updated_at, email
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

type CreateUserRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id, created_at, updated_at, email, hashed_password from users where email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
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

const updateUserEmailAndPassword = `-- name: UpdateUserEmailAndPassword :one
update users set email = $1, hashed_password = $2, updated_at = now()
where id = $3
returning id, created_at, updated_at, email
`

type UpdateUserEmailAndPasswordParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

type UpdateUserEmailAndPasswordRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

func (q *Queries) UpdateUserEmailAndPassword(ctx context.Context, arg UpdateUserEmailAndPasswordParams) (UpdateUserEmailAndPasswordRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserEmailAndPassword, arg.Email, arg.HashedPassword, arg.ID)
	var i UpdateUserEmailAndPasswordRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
	)
	return i, err
}
