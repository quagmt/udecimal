| Case | udecimal | shopspring | ericlagergren | apd | govalues |
|---|---|---|---|---|---|
| `MarshalJSON` `123` | 69.6 ns / 2 | 113 ns / 4 | 113 ns / 5 | 82.8 ns / 2 | 77.5 ns / 2 |
| `MarshalJSON` `123456.123456` | 79.5 ns / 2 | 127 ns / 4 | 126 ns / 5 | 89.8 ns / 2 | 82.9 ns / 2 |
| `MarshalJSON` `1234567890` | 72.4 ns / 2 | 108 ns / 4 | 114 ns / 5 | 86.9 ns / 2 | 76.3 ns / 2 |
| `MarshalJSON` `0.1234567890123456879` | 86.2 ns / 2 | 136 ns / 4 | 118 ns / 5 | 106 ns / 4 | 92.6 ns / 2 |
| `UnmarshalJSON` `123` | 64.2 ns / 0 | 88.2 ns / 3 | 100 ns / 1 | 116 ns / 1 | 70.2 ns / 0 |
| `UnmarshalJSON` `123456.123456` | 77.5 ns / 0 | 118 ns / 3 | 134 ns / 1 | 165 ns / 2 | 97.7 ns / 0 |
| `UnmarshalJSON` `1234567890` | 72.9 ns / 0 | 105 ns / 3 | 122 ns / 1 | 130 ns / 1 | 89.2 ns / 0 |
| `UnmarshalJSON` `0.1234567890123456879` | 98.7 ns / 0 | 265 ns / 6 | 168 ns / 1 | 168 ns / 2 | 117 ns / 0 |

Each cell shows the median ns/op and allocs/op.
