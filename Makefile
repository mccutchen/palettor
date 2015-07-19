benchmark:
	go test -v -cpu 1,2,4 -race -bench . -benchmem ./...

get-deps:
	go get -u github.com/golang/lint/golint

test:
	golint .
	go vet .
	go test -v -cpu 1,2,4 -race ./...

.PHONY: benchmark get-deps test
