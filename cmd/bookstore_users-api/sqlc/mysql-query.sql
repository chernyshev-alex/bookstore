
-- name: FindUser :one
SELECT id, first_name,last_name,email, date_created, status FROM users WHERE id=$1;

-- name: InsertUser :exec
INSERT INTO users (first_name,last_name,email,date_created, status, password) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUser :exec
UPDATE users SET first_name=$2,last_name=$3,email=$4 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1;

-- name: FindByStatus :many
SELECT id, first_name,last_name,email,date_created, status FROM users WHERE status=$1;

-- name: FindByEMailAndPsw :many
SELECT id, first_name,last_name,email,date_created, status FROM users WHERE email=$1 and password=$2 and status=$3;