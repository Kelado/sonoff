BINARY_NAME=smart-home

build:
	GOARCH=${ARCH} GOOS=${OS} go build -o bin/${BINARY_NAME} main.go

run: build 
	./bin/${BINARY_NAME}
