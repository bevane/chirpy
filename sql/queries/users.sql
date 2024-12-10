-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
	gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING id, created_at, updated_at, email;


-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: UpdateUserEmailAndPassword :one
update users set email = $1, hashed_password = $2, updated_at = now()
where id = $3
returning id, created_at, updated_at, email;
