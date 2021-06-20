.PHONY: tests
tests:
	go test ./...

.PHONY: int-tests
int-tests:
	export CONFIG=`pwd`/conf/app.yml;  \
	go test  -timeout 30s -run ^TestDaoSQLC$  github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/tests

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

.PHONY: wire-gen
wire-gen:
	wire

.PHONY: mocks-gen
mocks-gen:
	mockery --name=UserService --dir=./services/user_services  --output ./mocks 

.PHONY: app-run
app-run:
	go run main.go  
