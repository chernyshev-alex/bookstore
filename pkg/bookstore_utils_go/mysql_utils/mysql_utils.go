package mysql_utils

import (
	"strings"

	"github.com/chernyshev-alex/bookstore_utils_go/logger"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
	"github.com/go-sql-driver/mysql"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseErrors(err error) rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return rest_errors.NewNotFoundError("not found")
		}
		logger.Error(sqlErr.Error(), sqlErr)
		return rest_errors.NewInternalServerError("db error", err)
	}

	switch sqlErr.Number {
	case 1062:
		return rest_errors.NewBadRequestError("email is exists")
	}
	logger.Error(sqlErr.Error(), err)
	return rest_errors.NewInternalServerError("", err)
}
