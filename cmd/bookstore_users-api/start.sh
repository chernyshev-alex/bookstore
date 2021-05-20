#!/bin/bash

mysql.server start &

# integration test 
export CONFIG=`pwd`/config/config.yml
go test -run ^TestMain$  github.com/chernyshev-alex/bookstore_users-api/domain/users

# run
go run main.go -c config.yml

