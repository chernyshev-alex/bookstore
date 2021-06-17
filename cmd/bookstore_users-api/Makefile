.PHONY: tests
tests:
	go test ./...

.PHONY: int-tests
int-tests:
	export CONFIG=`pwd`/config/config.yml;  \
	go test  -timeout 30s -run ^TestDaoSQLC$  github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/sqlc/dao/users

.PHONY: docker-start
docker-start:
	docker-compose up -d

.PHONY: docker-stop
docker-stop:
	docker-compose down

.PHONY: sqlc-gen
sqlc-gen:
	sqlc  generate

.PHONY: mysql-start
mysql-start:
	mysql.server start &

.PHONY: mysql-stop
mysql-stop:
	mysql.server stop 