| Case | udecimal | shopspring | ericlagergren | apd | govalues |
|---|---|---|---|---|---|
| `Parse` `123` | 6.9 ns / 0 | 21.3 ns / 2 | 27.8 ns / 1 | 41.1 ns / 0 | 11.1 ns / 0 |
| `Parse` `123456.123456` | 13.3 ns / 0 | 37.1 ns / 2 | 56.7 ns / 1 | 74.5 ns / 1 | 28.2 ns / 0 |
| `Parse` `1234567890` | 10.7 ns / 0 | 32.4 ns / 2 | 46.5 ns / 1 | 43.7 ns / 0 | 23.2 ns / 0 |
| `Parse` `0.1234567890123456879` | 20.3 ns / 0 | 177 ns / 5 | 80.2 ns / 1 | 78.5 ns / 1 | 43.0 ns / 0 |
| `String` `123` | 13.7 ns / 1 | 39.9 ns / 2 | 38.2 ns / 4 | 13.5 ns / 0 | 6.7 ns / 0 |
| `String` `123456.123456` | 20.2 ns / 1 | 56.6 ns / 2 | 45.9 ns / 4 | 17.2 ns / 0 | 12.1 ns / 0 |
| `String` `1234567890` | 15.6 ns / 1 | 38.8 ns / 2 | 40.1 ns / 4 | 14.5 ns / 0 | 9.4 ns / 0 |
| `String` `0.1234567890123456879` | 22.2 ns / 1 | 59.4 ns / 2 | 46.0 ns / 4 | 36.0 ns / 2 | 18.0 ns / 0 |
| `MarshalText` `123` | 14.1 ns / 1 | 40.5 ns / 2 | 45.3 ns / 4 | 22.1 ns / 1 | 4.2 ns / 0 |
| `MarshalText` `123456.123456` | 20.5 ns / 1 | 57.1 ns / 2 | 50.0 ns / 4 | 25.3 ns / 1 | 10.5 ns / 0 |
| `MarshalText` `1234567890` | 15.9 ns / 1 | 40.2 ns / 2 | 44.4 ns / 4 | 21.5 ns / 1 | 7.8 ns / 0 |
| `MarshalText` `0.1234567890123456879` | 23.5 ns / 1 | 62.5 ns / 2 | 49.1 ns / 4 | 42.8 ns / 3 | 18.6 ns / 0 |
| `UnmarshalText` `123` | 8.0 ns / 0 | 29.8 ns / 3 | 31.1 ns / 1 | 48.8 ns / 1 | 12.6 ns / 0 |
| `UnmarshalText` `123456.123456` | 13.9 ns / 0 | 47.1 ns / 3 | 58.3 ns / 1 | 86.2 ns / 2 | 29.7 ns / 0 |
| `UnmarshalText` `1234567890` | 11.6 ns / 0 | 42.7 ns / 3 | 48.1 ns / 1 | 58.5 ns / 1 | 24.0 ns / 0 |
| `UnmarshalText` `0.1234567890123456879` | 20.0 ns / 0 | 189 ns / 6 | 82.3 ns / 1 | 94.8 ns / 2 | 41.9 ns / 0 |

Each cell shows the median ns/op and allocs/op.
