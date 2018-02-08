GOPATH     := $(GOPATH)
GIT_HASH   := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

all: tests

.PHONY: coverage
coverage:
	if [ -f coverage.out ]; then rm coverage.out; fi
	docker build -t dairycoverage --file coverage.Dockerfile .
	docker run --volume=$(GOPATH)/src/github.com/dairycart/dairycart:/output --rm -t dairycoverage
	go tool cover -html=coverage.out
	if [ -f coverage.out ]; then rm coverage.out; fi

.PHONY: tests
tests:
	# docker system prune -f
	docker build -t api_test -f test.Dockerfile .
	docker run --name api_test --rm api_test

.PHONY: integration
integration:
	docker-compose --file integration_tests.yml up --abort-on-container-exit --build --remove-orphans --force-recreate

.PHONY: vendor
vendor:
	dep ensure -update -v

.PHONY: revendor
revendor:
	rm -rf vendor
	rm Gopkg.*
	dep init -v

.PHONY: storage
storage:
	(cd storage && gnorm gen)

.PHONY: run
run:
	docker-compose --file docker-compose.yml up --build --remove-orphans --force-recreate --abort-on-container-exit