

tidy:
	go mod tidy 

fmt: tidy 
	goimports -w .

lint: fmt
	golangci-lint run 

test: fmt
	go test -v ./...

build-wasm:
	env GOOS=js GOARCH=wasm go build -o txnotify.wasm ./cmd/wasm/ && mv txnotify.wasm webui/public/
