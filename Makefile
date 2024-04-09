BINARY_NAME=sonoff

build:
	bin/templ generate
	npx tailwindcss -i ./static/input.css -o ./static/output.css
	GOARCH=${ARCH} GOOS=${OS} go build -o bin/${BINARY_NAME} main.go

build-exec:
	GOARCH=${ARCH} GOOS=${OS} go build -o bin/${BINARY_NAME} main.go
	
w-view:
	bin/templ generate -watch

w-css:
	npx tailwindcss -i ./static/input.css -o ./static/output.css -w

run: build 
	./bin/${BINARY_NAME}

clean: 
	rm view/**/*.txt

clean-deep: 
	rm view/**/*.txt || true
	rm view/**/*.go  || true