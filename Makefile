include .env

CONFIG_VARS = \
    PATH_CONFIG=${PATH_CONFIG} \
    SERVER_PORT=${SERVER_PORT} \
    QUOTES_FILE_PATH=${QUOTES_FILE_PATH}

start_server:
		$(CONFIG_VARS) go run cmd/server/main.go

start_client:
		$(CONFIG_VARS) go run cmd/client/main.go