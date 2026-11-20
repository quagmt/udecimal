| Benchmark | udecimal | shopspring | ericlagergren | apd | govalues |
|---|---|---|---|---|---|
| Add | 4.4 ns | 5.4× | 2.8× | 4.4× | 1.1× |
| Sub | 6.3 ns | 4.5× | 1.9× | 3.1× | **0.87×** |
| Mul | 4.2 ns | 5.1× | 2.5× | 5.3× | 3.4× |
| Div | 9.9 ns | 7.7× | 5.1× | 8.1× | 20.6× |
| Pow | 124 ns | **0.83×** | 2.5× | 5.1× | 3.2× |
| Parse | 11.9 ns | 3.9× | 4.1× | 4.8× | 2.0× |
| String | 17.6 ns | 2.7× | 2.4× | 1.1× | **0.62×** |
| MarshalText | 18.1 ns | 2.7× | 2.6× | 1.5× | **0.49×** |
| UnmarshalText | 12.7 ns | 4.6× | 4.1× | 5.5× | 2.0× |
| MarshalJSON | 76.6 ns | 1.6× | 1.5× | 1.2× | 1.1× |
| UnmarshalJSON | 77.4 ns | 1.7× | 1.7× | 1.8× | 1.2× |

Values in bold are faster than udecimal.
