# -v keeps the skip reasons (e.g. "cannot represent: ...") in the output; benchstat ignores those lines.
BENCH = go test -run NONE -bench=./$(1) -benchmem -count=10 -timeout=2h -v . | sed "s/$(1)\\///g" > bench-$(2).txt

bench-udec:
	@$(call BENCH,udec,udec)

bench-ss:
	@$(call BENCH,ss,ss)

bench-ed:
	@$(call BENCH,eric,ed)

bench-apd:
	@$(call BENCH,apd,apd)

bench-gov:
	@$(call BENCH,gov,gov)

bench-all: bench-udec bench-ss bench-ed bench-apd bench-gov

stats:
	@benchstat shopspring=bench-ss.txt udecimal=bench-udec.txt > bench-vs-shopspring.txt
	@benchstat ericlagergren=bench-ed.txt udecimal=bench-udec.txt > bench-vs-ericlagergren.txt
	@benchstat apd=bench-apd.txt udecimal=bench-udec.txt > bench-vs-apd.txt
	@benchstat govalues=bench-gov.txt udecimal=bench-udec.txt > bench-vs-govalues.txt

.PHONY: bench-udec bench-ss bench-ed bench-apd bench-gov bench-all stats
