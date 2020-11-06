MODULE = $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")
PACKAGES := $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
LDFLAGS := -ldflags "-X main.Version=${VERSION}"
ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))

CONFIG_FILE ?= ./configs/local.yml
APP_DSN ?= $(shell sed -n 's/^dsn:[[:space:]]*"\(.*\)"/\1/p' $(CONFIG_FILE))
MIGRATE := docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "$(APP_DSN)"

PID_FILE := './.pid'
FSWATCH_FILE := './fswatch.cfg'

.PHONY: default
default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## run unit tests
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: test-cover
test-cover: test ## run unit tests and show test coverage information
	go tool cover -html=coverage-all.out

.PHONY: run
run: ## run the API server
	go run ${LDFLAGS} cmd/server/*.go

.PHONY: run-debug
run-debug: ## run the API server with debug mode
	go run ${LDFLAGS} cmd/server/server.go -debug=true

.PHONY: run-restart
run-restart: ## restart the API server
	@pkill -P `cat $(PID_FILE)` || true
	@printf '%*s\n' "80" '' | tr ' ' -
	@echo "Source file changed. Restarting server..."
	@go run ${LDFLAGS} cmd/server/server.go & echo $$! > $(PID_FILE)
	@printf '%*s\n' "80" '' | tr ' ' -

run-live: ## run the API server with live reload support (requires fswatch)
	@go run ${LDFLAGS} cmd/server/server.go & echo $$! > $(PID_FILE)
	@fswatch -x -o --event Created --event Updated --event Renamed -r internal pkg cmd config | xargs -n1 -I {} make run-restart

.PHONY: build
build:  ## build the API server binary
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o dist/server $(MODULE)/cmd/server

.PHONY: build-docker
build-docker: ## build the API server as a docker image
	docker build -f cmd/server/Dockerfile -t server .

.PHONY: clean
clean: ## remove temporary files
	rm -rf server coverage.out coverage-all.out

.PHONY: version
version: ## display the version of the API server
	@echo $(VERSION)


.PHONY: goimport
goimport: ## run goimports on all files
	goimports -l $(GOFILES)


.PHONY: db-start
db-start: ## start the database server
	@mkdir -p testdata/postgres
	docker run --rm --name postgres -v $(shell pwd)/testdata:/testdata \
		-v $(shell pwd)/testdata/postgres:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go_restful -d -p 5432:5432 postgres

.PHONY: db-stop
db-stop: ## stop the database server
	docker stop postgres

.PHONY: testdata
testdata: ## populate the database with test data
	make migrate-reset
	@echo "Populating test data..."
	@docker exec -it postgres psql "$(APP_DSN)" -f /testdata/testdata.sql

.PHONY: lint
lint: ## run golint on all Go package
	@golint $(PACKAGES)

.PHONY: fmt
fmt: ## run "go fmt" on all Go packages
	@go fmt $(PACKAGES)

.PHONY: swagger-to-go
swagger-to-go:
	@$(foreach file,$(shell find -type f -name "*.json"), \
		tire swagger-to-go $(ROOT)/${file} --pkg "auth" --out "$(ROOT)/internal/proto/v1/$(${file}).go" ;)

.PHONY: proto
proto:  ## run "prototool generate"
	@echo "Running prototool generate..."
	@prototool generate

.PHONY: gen
gen:  ## run proto and swagger-to-go
	@echo "generate..."
	proto
	swagget-to-go


.PHONY: migrate
migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir /migrations/ $${name// /_}

.PHONY: migrate-reset
migrate-reset: ## reset database and re-run all migrations
	@echo "Resetting database..."
	@$(MIGRATE) drop
	@echo "Running all database migrations..."
	@$(MIGRATE) up

.PHONY: install-swagger-ui
install-swagger-ui:
	git clone --depth 1 https://github.com/swagger-api/swagger-ui.git ./tmp/swagger-ui
	rm -rf ./third_party/swagger-ui
	mv ./tmp/swagger-ui/dist ./third_party/swagger-ui
	rm -rf ./tmp/swagger-ui
	sed -i -e 's/https:\/\/petstore.swagger.io\/v2\/swagger\.json/\/v1\/swagger\/index\.json/g' ./third_party/swagger-ui/index.html
	cd $(ROOT)/third_party/swagger-ui && go-bindata -nometadata -o $(ROOT)/pkg/grpcgw/swagger.gen.go -nomemcopy=true -pkg=grpcgw ./...
