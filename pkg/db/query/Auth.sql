-- name: Auth_GetUserByEmail :one
select id, name, email, password, verified, created_at
from users
where lower(email) = @email
order by id;

-- name: Auth_GetUserById :one
select id, name, email, password, verified, created_at
from users
where id = @id
limit 1;

-- name: Auth_SaveUser :one
insert into users(name, email, password, created_at)
values (@name, @email, @password, date())
returning *;

-- name: Auth_MarkUserAsVerified :one
update users
set verified = true
where id = @id
returning *;

-- name: Auth_UpdatePassword :one
update users
set password = @password
where id = @id
returning *;

-- name: Auth_SavePasswordToken :one
insert into password_tokens(password_token_user, hash, created_at)
values (@user, @hash, date())
returning *;

-- name: Auth_GetTokenById :one
select id, hash, created_at, password_token_user
from password_tokens
where id = @id
  and password_token_user = @user
  and created_at >= @createdAt
limit 1;

-- name: Auth_DeleteAllPasswordTokensByUser :exec
delete
from password_tokens
where password_token_user = @user;
