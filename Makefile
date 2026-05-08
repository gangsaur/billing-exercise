test:
	go test -tags mock ./...

generate-mock:
	mockery
