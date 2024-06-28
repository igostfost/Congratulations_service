.PHONY: test

run:
	go run cmd/main.go

tests:
	go test -v ./pkg/repository/...