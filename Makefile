GOPATH     := $(GOPATH)
GIT_HASH   := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

# Basics

all: tests run

clean:
	rm api/v1/example_files/plugins/mock_db.so
	rm api/v1/example_files/plugins/mock_img.so

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

.PHONY: ci-coverage
ci-coverage: | example-plugins
	docker build -t dairycoverage --file dockerfiles/coverage.Dockerfile .
	docker run --volume=$(GOPATH)/src/github.com/dairycart/dairycart:/output --rm -t dairycoverage

# Testing

.PHONY: tests
tests:
	make models-unit-tests api-unit-tests storage-unit-tests integration-tests

.PHONY: api-unit-tests
api-unit-tests: | example-plugins
	# api unit tests
	docker build -t api_test -f dockerfiles/unittest.Dockerfile .
	docker run --name api_test --rm api_test

.PHONY: models-unit-tests
models-unit-tests:
	go test -cover github.com/dairycart/dairycart/models/v1

.PHONY: storage-unit-tests
storage-unit-tests:
	go test -cover github.com/dairycart/dairycart/storage/database/postgres

.PHONY: integration-tests
integration-tests:
	docker-compose --file integration-tests.yml up --build --remove-orphans --force-recreate --abort-on-container-exit

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
	(cd storage/database && gnorm gen)

# What we're all here for

.PHONY: run
run:
	docker-compose --file docker-compose.yml up --build --remove-orphans --force-recreate --abort-on-container-exit
