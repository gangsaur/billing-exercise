test:
	go test -timeout 30s -tags mock ./...

coverage:
	go test -timeout 30s -coverprofile=coverage.out -tags=mock ./...
	go tool cover -html=./coverage.out -o=./coverage.html

generate-mock:
	mockery
