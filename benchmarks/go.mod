module github.com/quagmt/udecimal/benchmarks

go 1.23

replace github.com/quagmt/udecimal => ../

require (
	github.com/cockroachdb/apd/v3 v3.2.3
	github.com/ericlagergren/decimal v0.0.0-20240411145413-00de7ca16731
	github.com/govalues/decimal v0.1.36
	github.com/quagmt/udecimal v0.0.0-00010101000000-000000000000
	github.com/shopspring/decimal v1.6.0
	github.com/stretchr/testify v1.12.1
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect
