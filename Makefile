include .example.env

# Check if .env file exists
ifeq ($(wildcard .env),)
    # If .env file doesn't exist, print a message
    $(info .env file not found, skipping inclusion)
else
    # If .env file exists, include it. It need to be created in the root of the project, this import will override the variables in the .example.env file
    include .env
endif

export
PROJECT_ID=GCP_PROJECT_ID
SERVICE_NAME=kickstart
BINARY_NAME=app
GITHUB_ORGANIZATION=GITHUB_ORGANIZATION_NAME
IMPORT_PATH=github.com/${GITHUB_ORGANIZATION}/${SERVICE_NAME}
GIT_COMMIT_HASH=$(shell git rev-parse HEAD | cut -c -16)
BUILD_TIME=$(shell date +%s)
LDFLAGS = -X ${IMPORT_PATH}/global.ServiceName=${SERVICE_NAME}
LDFLAGS += -X ${IMPORT_PATH}/global.GitCommitHash=${GIT_COMMIT_HASH}
LDFLAGS += -X ${IMPORT_PATH}/global.BuildTime=${BUILD_TIME}
TAG=${PROJECT_ID}/${SERVICE_NAME}:${GIT_COMMIT_HASH}
IMAGE=asia.gcr.io/${TAG}
TEST_DEPLOY_NAME=DEV_DEPLOY_NAME
TEST_CLUSTER_NAME=DEV_CLUSTER_NAME
DEPLOY_NAME=kickstart
CLUSTER_NAME=PROD_CLUSTER_NAME
REGION=PROD_CLUSTER_REGION

.PHONY: clean doc deps run deploy

${BINARY_NAME}:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/${BINARY_NAME} -ldflags "$(LDFLAGS)" ./main.go

local-dev-up: ### Run docker-compose
	docker compose up -d --no-deps --build

local-dev-down: ### Run docker-compose
	docker compose down

doc: ### swag init # export PATH=$(go env GOPATH)/bin:$PATH, ref: https://github.com/swaggo/swag/issues/197
	swag init -g api/api.go --parseDependency --parseInternal

deps:
	go mod tidy
	go list -m all

clean:
	rm -f ./${BINARY_NAME}

mocks:
	find . -type f -name 'mock_*.go' -exec rm {} +
	./mockery --all --inpackage

# https://github.com/golang-migrate/migrate, check here for more detail on database migration
migration-script:
	./migrate create -ext sql -dir database/migrations -seq create_new_table

unit-test:
	go test -v -cover -short ./...

db-integration-test:
	go test -v -cover -run DBIntegration ./...

api-integration-test:
	go test -v -cover -run APIIntegration ./...

integration-test:
	go test -v -cover -run Integration ./...

full-test:
	go test -v -cover ./...

image: 
	docker build -t ${TAG} .
	docker tag ${TAG} ${IMAGE}
	docker tag ${TAG} asia.gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest

run: 
	GIN_MODE=debug go run main.go
