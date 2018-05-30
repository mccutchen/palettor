test:
	golint .
	go vet .
	go test -v -cpu 1,2,4 -race ./...

benchmark:
	go test -cpu 1,2,4 -race -bench ./... -benchmem ./...

get-deps:
	go get -u github.com/golang/lint/golint

app:
	mkdir -p dist
	go build -o dist/palettor cmd/palettor/palettor.go

.PHONY: test benchmark get-deps
