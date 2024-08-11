bench:
	go test -bench BenchmarkDiv -benchmem -memprofile mem.out -cpuprofile cpu.out
 
lint:
	golangci-lint run ./... -v

test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
