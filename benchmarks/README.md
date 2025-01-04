# Benchmarks

This section provides benchmarks for the `udecimal` library in comparision with other libraries. Some benchmarks have `fallback` suffix which means the operation can't be done by using `u128` and falls back to `big.Int` API.

## Benchmark Results

<i>**NOTE**: The results are for references and can be varied depending on the hardware</i>

### Environment

```powershell
goos: linux
goarch: amd64
pkg: github.com/quagmt/udecimal/benchmarks
cpu: Intel(R) Core(TM) i9-14900HX
```

### Results

Full benchmark details can be found in [bench-udec.txt](bench-udec.txt)

```powershell

# Arithmetic operations
BenchmarkAdd/1234567890123456789.1234567890123456879.Add(1111.1789)-32                     	100000000	 10.91 ns/op	   0 B/op	 0 allocs/op
BenchmarkAdd/548751.15465466546.Add(1542.456487)-32                                      	100000000	 11.44 ns/op	   0 B/op	 0 allocs/op
BenchmarkSub/1234567890123456789.1234567890123456879.Sub(1111.1789)-32                   	 96477370	 13.21 ns/op	   0 B/op	 0 allocs/op
BenchmarkSub/123.456.Sub(0.123)-32                                                       	100000000	 10.75 ns/op	   0 B/op	 0 allocs/op
BenchmarkMul/1234.1234567890123456879.Mul(1111.1789)-32                                  	100000000	 10.43 ns/op	   0 B/op	 0 allocs/op
BenchmarkMul/1234.1234567890123456879.Mul(1111.1234567890123456789)-32                   	100000000	 10.49 ns/op	   0 B/op	 0 allocs/op
BenchmarkDiv/1234567890123456789.1234567890123456879.Div(1111.1789)-32                    	100000000	 10.46 ns/op	   0 B/op	 0 allocs/op
BenchmarkDiv/12345.1234567890123456879.Div(1111.1234567890123456789)-32                   	 50097552	 25.23 ns/op	   0 B/op	 0 allocs/op
BenchmarkDiv/123.456.Div(0.123)-32                                                        	111231829	 10.70 ns/op	   0 B/op	 0 allocs/op
BenchmarkDiv/548751.15465466546.Div(1542.456487)-32                                       	100000000	 10.30 ns/op	   0 B/op	 0 allocs/op
BenchmarkDivFallback/12345679012345679890123456789.1234567890123456789.Div(999999)-32     	  4187094	 302.1 ns/op	 264 B/op	 7 allocs/op
BenchmarkDivFallback/1234.Div(12345679012345679890123456789.1234567890123456789)-32       	  3901846	 306.7 ns/op	 320 B/op	 7 allocs/op
BenchmarkPow/1.01.Pow(10)-32                                                              	 30365348	 41.14 ns/op	   0 B/op	 0 allocs/op
BenchmarkPow/1.01.Pow(100)-32                                                            	   994129	  1137 ns/op	 817 B/op	13 allocs/op

# Parsing string
BenchmarkParse/1234567890123456789.1234567890123456879-32                                  	 32111433	 38.21 ns/op	   0 B/op	 0 allocs/op
BenchmarkParse/1234567890-32                                                               	 98585916	 12.58 ns/op	   0 B/op	 0 allocs/op
BenchmarkParse/0.1234567890123456879-32                                                    	 44339668	 26.45 ns/op	   0 B/op	 0 allocs/op
BenchmarkParseFallBack/123456789123456789123456.1234567890123456-32                        	  2805122	 473.3 ns/op	 192 B/op	 6 allocs/op
BenchmarkParseFallBack/111222333444555666777888999.1234567890123456789-32                  	  2442004	 500.8 ns/op	 216 B/op	 6 allocs/op
BenchmarkString/1234567890123456789.1234567890123456879-32                                 	 14577884	 76.50 ns/op	  48 B/op	 1 allocs/op
BenchmarkString/0.1234567890123456879-32                                                   	 41109242	 40.02 ns/op	  24 B/op	 1 allocs/op
BenchmarkStringFallBack/123456789123456789123456.1234567890123456-32                       	  4147044	 256.2 ns/op	 208 B/op	 4 allocs/op
BenchmarkStringFallBack/111222333444555666777888999.1234567890123456789-32                 	  3808071	 313.3 ns/op	 208 B/op	 4 allocs/op

# Marshal/Unmarshal
BenchmarkMarshalJSON/1234567890123456789.1234567890123456879-32                            	 13965998	 77.22 ns/op	  48 B/op	 1 allocs/op
BenchmarkMarshalJSON/0.1234567890123456879-32                                              	 24039360	 43.57 ns/op	  24 B/op	 1 allocs/op
BenchmarkMarshalJSON/12345678901234567891234567890123456789.1234567890123456879-32         	  3445560	 291.6 ns/op	 320 B/op	 5 allocs/op
BenchmarkUnmarshalJSON/1234567890123456789.1234567890123456879-32                          	 15943234	 73.77 ns/op	   0 B/op	 0 allocs/op
BenchmarkUnmarshalJSON/123456.123456-32                                                    	 46983879	 26.55 ns/op	   0 B/op	 0 allocs/op
BenchmarkUnmarshalJSON/12345678901234567891234567890123456789.1234567890123456879-32       	  2267604	 517.4 ns/op	 264 B/op	 6 allocs/op
BenchmarkMarshalBinary/1234567890123456789.1234567890123456879-32                          	 50875198	 25.97 ns/op	  24 B/op	 1 allocs/op
BenchmarkMarshalBinary/0.1234567890123456879-32                                            	 54470340	 20.91 ns/op	  16 B/op	 1 allocs/op
BenchmarkMarshalBinary/12345678901234567891234567890123456789.1234567890123456879-32       	 21138375	 48.85 ns/op	  32 B/op	 1 allocs/op
BenchmarkUnmarshalBinary/1234567890123456789.1234567890123456879-32                       	554818034	 2.034 ns/op	   0 B/op	 0 allocs/op
BenchmarkUnmarshalBinary/0.1234567890123456879-32                                          	637610913	 1.822 ns/op	   0 B/op	 0 allocs/op
```

## Comparision with other libraries

Same benchmarks are performed on other libraries and the results are compared with `udecimal` using benchstat tool.
Libraries used for comparision are:

- [shopspring/decimal](https://github.com/shopspring/decimal)
- [ericlagergren/decimal](https://github.com/ericlagergren/decimal)

### Results

- [bench-vs-shopspring.txt](bench-vs-shopspring.txt)
- [bench-vs-ericlagergren.txt](bench-vs-ericlagergren.txt)
