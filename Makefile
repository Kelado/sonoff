BINARY_NAME=local-smart-home
PROD_BINARY_NAME=smart-home

build:
	GOARCH=${ARCH} GOOS=${OS} go build -o bin/${BINARY_NAME} main.go

run: build
	./bin/${BINARY_NAME}

push-to-server:
	GOARCH=amd64 GOOS=linux go build -o bin/${PROD_BINARY_NAME} main.go
	@echo "Pushing to server..."
	@scp bin/${PROD_BINARY_NAME} .env obito@konoha.local:~/apps/smart-home || echo "Failed to copy files"
