.PHONY: sqlc-gen

test:
	go test ./...

dao-test:
	export CONFIG=`pwd`/config/config.yml;  \
	go test  -timeout 30s -run ^TestDaoSQLC$  github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/sqlc/dao/users

start:
	docker-compose up -d

sqlc-gen:
	sqlc  generate

mysql-start:
	mysql.server start &

mysql-stop:
	mysql.server stop 