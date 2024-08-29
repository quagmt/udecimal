bench:
	@go test -bench BenchmarkString -benchmem -memprofile mem.out -cpuprofile cpu.out -run NONE
 
lint:
	@golangci-lint run ./... -v

test:
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out
