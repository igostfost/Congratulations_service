.PHONY: test

run:
	go run cmd/main.go

test:
	go test -v ./pkg/repository/...