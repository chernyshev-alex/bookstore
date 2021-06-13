// Code generated by sqlc. DO NOT EDIT.
// source: mysql-query.sql

package user_dao

import (
	"context"
	"database/sql"
	"time"
)

const deleteUser = `-- name: DeleteUser :execresult
DELETE FROM users WHERE id=?
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) (sql.Result, error) {
	return q.exec(ctx, q.deleteUserStmt, deleteUser, id)
}

const findByEMailAndPsw = `-- name: FindByEMailAndPsw :one
SELECT id, first_name,last_name,email,date_created, status FROM users WHERE email=? and password=? and status=?
`

type FindByEMailAndPswParams struct {
	Email    string         `json:"email"`
	Password sql.NullString `json:"password"`
	Status   sql.NullString `json:"status"`
}

type FindByEMailAndPswRow struct {
	ID          int32          `json:"id"`
	FirstName   sql.NullString `json:"firstName"`
	LastName    sql.NullString `json:"lastName"`
	Email       string         `json:"email"`
	DateCreated time.Time      `json:"dateCreated"`
	Status      sql.NullString `json:"status"`
}

func (q *Queries) FindByEMailAndPsw(ctx context.Context, arg FindByEMailAndPswParams) (FindByEMailAndPswRow, error) {
	row := q.queryRow(ctx, q.findByEMailAndPswStmt, findByEMailAndPsw, arg.Email, arg.Password, arg.Status)
	var i FindByEMailAndPswRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.DateCreated,
		&i.Status,
	)
	return i, err
}

const findByStatus = `-- name: FindByStatus :many
SELECT id, first_name,last_name,email,date_created, status FROM users WHERE status=?
`

type FindByStatusRow struct {
	ID          int32          `json:"id"`
	FirstName   sql.NullString `json:"firstName"`
	LastName    sql.NullString `json:"lastName"`
	Email       string         `json:"email"`
	DateCreated time.Time      `json:"dateCreated"`
	Status      sql.NullString `json:"status"`
}

func (q *Queries) FindByStatus(ctx context.Context, status sql.NullString) ([]FindByStatusRow, error) {
	rows, err := q.query(ctx, q.findByStatusStmt, findByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindByStatusRow
	for rows.Next() {
		var i FindByStatusRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.DateCreated,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUser = `-- name: FindUser :one
SELECT id, first_name,last_name,email, date_created, status FROM users WHERE id = ?
`

type FindUserRow struct {
	ID          int32          `json:"id"`
	FirstName   sql.NullString `json:"firstName"`
	LastName    sql.NullString `json:"lastName"`
	Email       string         `json:"email"`
	DateCreated time.Time      `json:"dateCreated"`
	Status      sql.NullString `json:"status"`
}

func (q *Queries) FindUser(ctx context.Context, id int32) (FindUserRow, error) {
	row := q.queryRow(ctx, q.findUserStmt, findUser, id)
	var i FindUserRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.DateCreated,
		&i.Status,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :execresult
INSERT INTO users (first_name,last_name,email,date_created, status, password) VALUES (?, ?, ?, ?, ?, ?)
`

type InsertUserParams struct {
	FirstName   sql.NullString `json:"firstName"`
	LastName    sql.NullString `json:"lastName"`
	Email       string         `json:"email"`
	DateCreated time.Time      `json:"dateCreated"`
	Status      sql.NullString `json:"status"`
	Password    sql.NullString `json:"password"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (sql.Result, error) {
	return q.exec(ctx, q.insertUserStmt, insertUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.DateCreated,
		arg.Status,
		arg.Password,
	)
}

const updateUser = `-- name: UpdateUser :execresult
UPDATE users SET first_name=?,last_name=?,email=? WHERE id = ?
`

type UpdateUserParams struct {
	FirstName sql.NullString `json:"firstName"`
	LastName  sql.NullString `json:"lastName"`
	Email     string         `json:"email"`
	ID        int32          `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.exec(ctx, q.updateUserStmt, updateUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.ID,
	)
}
