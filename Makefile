.PHONY: mock

test: mock
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

mock: mock-init
	go generate ./...

mock-init:
	go install github.com/vektra/mockery/v2@v2.16.0
