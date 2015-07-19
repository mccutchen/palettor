test:
	go test -v -cpu 1,2,4 -race ./...

benchmark:
	go test -v -cpu 1,2,4 -race -bench . -benchmem ./...

.PHONY: benchmark test
