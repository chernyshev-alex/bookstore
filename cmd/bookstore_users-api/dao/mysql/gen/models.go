// Code generated by sqlc. DO NOT EDIT.

package gen

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int32
	FirstName   sql.NullString
	LastName    sql.NullString
	Email       string
	DateCreated time.Time
	Status      sql.NullString
	Password    sql.NullString
}
