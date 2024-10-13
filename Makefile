.PHONY: test lint fuzz fuzz-all

lint:
	@golangci-lint --config=.golangci.yaml run ./... -v

test:
	# run all unit-tests, ignore fuzz tests
	@go test -tags='!fuzz' -v -race -failfast -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out

fuzz:
	$(eval fuzzName := $(filter-out $@,$(MAKECMDGOALS)))
	@go test -tags='fuzz' -v -run=Fuzz -fuzz=$(fuzzName) -fuzztime=30s -timeout=10m

fuzz-all:
	$(eval fuzzTime := $(filter-out $@,$(MAKECMDGOALS)))
	@sh scripts/fuzz-all.sh $(fuzzTime)

bench:
	@go test -bench=BenchmarkUnmarshalText -benchmem -benchmem -memprofile=mem.out -cpuprofile=cpu.out -run NONE

# https://stackoverflow.com/questions/6273608/how-to-pass-argument-to-makefile-from-command-line
%:
	@
