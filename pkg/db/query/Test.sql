-- name: Test_UpdatePasswordTokenCreatedAt :execresult
update password_tokens
set created_at = @createdAt
where id = @id;

-- name: Test_CountPasswordTokensByUser :one
select count(*)
from password_tokens
where password_token_user = @user;
