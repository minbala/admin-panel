-- name: CreateSession :exec
insert into sessions ( access_token, user_id) values ($1,$2);

-- name: DeleteSession :exec
delete from sessions where user_id =$1 and access_token =$2;

-- name: GetSession :one
select * from  sessions where user_id =$1 and access_token =$2;

-- name: GetUsers :many
select * from users where (id = $1 or $1 = 0) and (email = $2 or $2 = '') and (user_role::text like $3 or $3 = '%%') and (name like $4 or $4 = '%%') order by id desc limit $5 offset $6;

-- name: GetUsersCount :one
select count(*) from users where (id = $1 or $1 = 0) and (email = $2 or $2 = '') and (user_role::text like $3 or $3 = '%%') and (name like $4 or $4 = '%%');

-- name: CreateUser :one
insert into users (email, name, password, user_role) values ($1,$2,$3,$4) returning id;

-- name: DeleteUser :exec
delete  from users where  id =  $1;

-- name: UpdateUser :exec
update users set email = $1, name = $2, user_role = $3, password =$4, updated_at = clock_timestamp() where id = $5;

-- name: GetUserLogs :many
select * from user_logs where (user_id = $1 or $1 is  null) order by id desc limit  $2 offset $3;

-- name: GetUserLogsCount :one
select count(*) from user_logs where (user_id = $1 or $1 is null);

-- name: CreateUserLog :exec
insert into user_logs (user_id, event, request_url, data, status, error_message) values ($1,$2,$3,$4,$5,$6);
