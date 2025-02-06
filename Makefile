include .env

IMG_NAME_POW_SERVER=pow-server
IMG_NAME_POW_CLIENT=pow-client
NETWORK_NAME=pow-network

CONFIG_VARS = \
    PATH_CONFIG=${PATH_CONFIG} \
    SERVER_PORT=${SERVER_PORT} \
    QUOTES_FILE_PATH=${QUOTES_FILE_PATH}

test:
	go test ./...

local_run_server:
	${CONFIG_VARS} go run cmd/server/main.go

local_run_client:
	${CONFIG_VARS} go run cmd/client/main.go localhost:${SERVER_PORT}

build_server:
	docker build -t ${IMG_NAME_POW_SERVER} -f .deploy/Dockerfile.server .

start_server:
	docker run --rm --network=${NETWORK_NAME} --env-file .env ${IMG_NAME_POW_SERVER}

build_client:
	docker build -t ${IMG_NAME_POW_CLIENT} -f .deploy/Dockerfile.client .

start_client:
	docker run --rm --network=${NETWORK_NAME} --env-file .env ${IMG_NAME_POW_CLIENT}
