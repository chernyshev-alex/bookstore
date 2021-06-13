#!/bin/bash

mysql.server start &

# integration test 
export CONFIG=`pwd`/config/config.yml
#@deprecated
go test -run ^TestMain$  github.com/chernyshev-alex/bookstore_users-api/domain/users
#use this
go test  -run ^TestMain$ github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/sqlc/dao/users

# generate sqlc/dao/users
sqlc  generate

# run
go run main.go -c config.yml

