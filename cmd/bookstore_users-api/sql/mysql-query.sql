
-- name: FindUser :one
SELECT id, first_name,last_name,email, date_created, status FROM users WHERE id = ?;

-- name: InsertUser :execresult
INSERT INTO users (first_name,last_name,email,date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateUser :execresult
UPDATE users SET first_name=?,last_name=?,email=? WHERE id = ?;

-- name: DeleteUser :execresult
DELETE FROM users WHERE id=?;

-- name: FindByStatus :many
SELECT id, first_name,last_name,email,date_created, status FROM users WHERE status=?;

-- name: FindByEMailAndPsw :one
SELECT id, first_name,last_name,email,date_created, status FROM users WHERE email=? and password=? and status=?;