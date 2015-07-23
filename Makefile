benchmark:
	go test -bench . -benchmem ./...

get-deps:
	go get -u github.com/golang/lint/golint

test:
	golint .
	go vet .
	go test -v ./...

.PHONY: benchmark get-deps test
