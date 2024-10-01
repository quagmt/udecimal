
bench:
	@go test -bench BenchmarkNew -benchmem -memprofile mem.out -cpuprofile cpu.out -run NONE

bench-udec:
	@rm -f bench-udec.txt
	@go test -bench=./udec -benchmem -count=10 -timeout=1h> bench-udec.txt
	@sed -i 's/udec\///g' bench-udec.txt

bench-ss:
	@rm -f bench-ss.txt
	@go test -bench=./ss -benchmem -count=10 -timeout=1h> bench-ss.txt
	@sed -i 's/ss\///g' bench-ss.txt

stats: bench-ss bench-udec
	@benchstat shopspring=bench-ss.txt udecimal=bench-udec.txt > benchstat.txt
