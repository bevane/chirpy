-- name: CreateRefreshToken :exec
insert into refresh_tokens (token, created_at, updated_at, user_id, expires_at)
values (
	$1, NOW(), NOW(), $2, $3
);


-- name: GetRefreshToken :one
select * from refresh_tokens where token = $1;

-- name: RevokeRefreshToken :exec

update refresh_tokens set updated_at = NOW(), revoked_at = NOW()
where token = $1;
