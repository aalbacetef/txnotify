

tidy:
	go mod tidy 

fmt: tidy 
	goimports -w .

lint: fmt
	golangci-lint run 

test: fmt
	go test -v ./...


