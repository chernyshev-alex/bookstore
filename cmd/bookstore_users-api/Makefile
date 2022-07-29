.PHONY: build test int-tests docker-start docker-stop sqlc-gen mysql-start mysql-stop wire-gen mocks-gen app-run

GOCMD:=$(shell which go)
GOLINT:=$(shell which golint)
GOIMPORT:=$(shell which goimports)
GOFMT:=$(shell which gofmt)
GOBUILD:=$(GOCMD) build
GOINSTALL:=$(GOCMD) install
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOLIST:=$(GOCMD) list
GOVET:=$(GOCMD) vet
GOPATH:=$(shell $(GOCMD) env GOPATH)
u := $(if $(update),-u)

# https://github.com/swaggo/swag/blob/master/Makefile
BINARY_NAME:=swag
PACKAGES:=$(shell $(GOLIST) github.com/swaggo/swag github.com/swaggo/swag/cmd/swag github.com/swaggo/swag/gen)
GOFILES:=$(shell find . -name "*.go" -type f)

export GO111MODULE := on

# test:
# 	echo "mode: count" > coverage.out
# 	for PKG in $(PACKAGES); do \
# 		$(GOCMD) test -v -covermode=count -coverprofile=profile.out $$PKG > tmp.out; \
# 		cat tmp.out; \
# 		if grep -q "^--- FAIL" tmp.out; then \
# 			rm tmp.out; \
# 			exit 1; \
# 		elif grep -q "build failed" tmp.out; then \
# 			rm tmp.out; \
# 			exit; \
# 		fi; \
# 		if [ -f profile.out ]; then \
# 			cat profile.out | grep -v "mode:" >> coverage.out; \
# 			rm profile.out; \
# 		fi; \
# 	done

build:
	go build ./...

tests:
	go test ./...

int-tests:
	export CONFIG=`pwd`/conf/app.yml;  \
	go test  -timeout 30s -run ^TestDaoSQLC$  github.com/chernyshev-alex/bookstore/cmd/bookstore_users_api/tests

docker-start:
	docker-compose up -d

docker-stop:
	docker-compose down

sqlc-gen:
	sqlc  generate

mysql-start:
	mysql.server start &

mysql-stop:
	mysql.server stop 

wire-gen:
	wire

mocks-gen:
	mockery --name=UserService --dir=./services/user_services  --output ./mocks 

app-run:
	go run main.go  
