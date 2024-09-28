.PHONY: all build clean test bench lint fuzz

bench:
	@go test -bench BenchmarkMarshalBinary -benchmem -memprofile mem.out -cpuprofile cpu.out -run NONE
 
lint:
	@golangci-lint --config=.golangci.yaml run ./... -v

test:
	# run all unit-tests, ignore fuzz tests
	@go test -tags='!fuzz' -v -race -failfast -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out

fuzz:
	$(eval fuzzName := $(filter-out $@,$(MAKECMDGOALS)))
	@go test -tags='fuzz' -v -run=Fuzz -fuzz=$(fuzzName) -fuzztime=30s -timeout=10m ./...

fuzz-all:
	echo "Run all fuzz tests"
	for fuzz_test in $(shell go test -list "^Fuzz" $$fuzz_pkg | grep "^Fuzz"); do \
		echo "Fuzzing $$fuzz_test in $$fuzz_pkg ..."; \
		go test -tags='fuzz' -run=Fuzz -fuzz=$$fuzz_test -fuzztime=30s $$fuzz_pkg -timeout=15m || exit 1; \
	done \
