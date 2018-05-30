GOPATH     := $(GOPATH)
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD)
GIT_HASH   := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

# Basics

all: tests run

clean:
	rm -f api/v1/example_files/plugins/mock_db.so
	rm -f api/v1/example_files/plugins/mock_img.so

.PHONY: tools
tools:
	go get -u gnorm.org/gnorm
	go get -u github.com/rakyll/statik

# Coverage Reports

.PHONY: coverage-report
coverage-report:
	make api-coverage-report models-coverage-report

.PHONY: api-coverage-report
api-coverage-report: | example-plugins
	if [ -f coverage.out ]; then rm coverage.out; fi
	docker build -t dairycoverage --file dockerfiles/coverage.Dockerfile .
	docker run --volume=$(GOPATH)/src/github.com/dairycart/dairycart:/output --rm -t dairycoverage
	go tool cover -html=coverage.out
	if [ -f coverage.out ]; then rm coverage.out; fi

.PHONY: models-coverage-report
models-coverage-report:
	if [ -f coverage.out ]; then rm coverage.out; fi
	go test -coverprofile=coverage.out github.com/dairycart/dairycart/models/v1
	go tool cover -html=coverage.out
	if [ -f coverage.out ]; then rm coverage.out; fi

.PHONY: storage-postgres-coverage-report
storage-postgres-coverage-report:
	if [ -f coverage.out ]; then rm coverage.out; fi
	go test -coverprofile=coverage.out github.com/dairycart/dairycart/storage/v1/database/postgres
	go tool cover -html=coverage.out
	if [ -f coverage.out ]; then rm coverage.out; fi

.PHONY: ci-coverage
ci-coverage: | example-plugins
	docker build -t dairycoverage --file dockerfiles/coverage.Dockerfile .
	docker run --volume=$(GOPATH)/src/github.com/dairycart/dairycart:/output --rm -t dairycoverage

# Testing

.PHONY: tests
tests:
	make unit-tests integration-tests

.PHONY: unit-tests
unit-tests:
	make models-unit-tests api-unit-tests client-unit-tests storage-postgres-unit-tests

.PHONY: api-unit-tests
api-unit-tests: | example-plugins
	# api unit tests
	docker build -t api_test -f dockerfiles/unittest.Dockerfile .
	docker run --name api_test --rm api_test

.PHONY: client-unit-tests
client-unit-tests:
	go test -cover github.com/dairycart/dairycart/client/v1

.PHONY: models-unit-tests
models-unit-tests:
	go test -cover github.com/dairycart/dairycart/models/v1

.PHONY: storage-postgres-unit-tests
storage-postgres-unit-tests:
	go test -cover github.com/dairycart/dairycart/storage/v1/database/postgres

.PHONY: storage-local-image-unit-tests
storage-local-image-unit-tests:
	go test -cover github.com/dairycart/dairycart/storage/v1/images/local

.PHONY: integration-tests
integration-tests:
	docker-compose --file compose-files/integration-tests.yml up --build --remove-orphans --force-recreate --abort-on-container-exit

.PHONY: quick-coverage
quick-coverage:
	go test -cover \
		github.com/dairycart/dairycart/api/v1 \
		github.com/dairycart/dairycart/client/v1 \
		github.com/dairycart/dairycart/cmd/admin/v1/server \
		github.com/dairycart/dairycart/storage/v1/database/postgres \
		github.com/dairycart/dairycart/storage/v1/images/local \

# Dependency management

.PHONY: vendor
vendor:
	dep ensure -update -v

.PHONY: revendor
revendor:
	rm -rf vendor
	rm Gopkg.*
	dep init -v

# Plugin files

.PHONY: example-plugins
example-plugins:
	make api/v1/example_files/plugins/mock_db.so api/v1/example_files/plugins/mock_img.so

api/v1/example_files/plugins/mock_db.so api/v1/example_files/plugins/mock_img.so:
	docker build -t plugins --file dockerfiles/example_plugins.Dockerfile .
	docker run --volume=$(GOPATH)/src/github.com/dairycart/dairycart/api/v1/example_files/plugins:/output --rm -t plugins

# Generated Code

.PHONY: models
models:
	(cd models/v1 && gnorm gen --config="gnorm_postgres.toml")

.PHONY: storage
storage:
	docker-compose --file compose-files/storage-gen.yml up --build --remove-orphans --force-recreate --abort-on-container-exit

.PHONY: assets
assets:
	statik -src=cmd/admin/v1/server/html -dest=cmd/admin/v1/server/html -f

# Docker images

.PHONY: dairydev
dairydev:
	docker build --tag dairycart/dairydev:latest --file=dockerfiles/dairydev.Dockerfile .
	docker push dairycart/dairydev:latest

# What we're all here for

.PHONY: run
run:
	docker-compose --file compose-files/docker-compose.yml up --build --remove-orphans --force-recreate --abort-on-container-exit

