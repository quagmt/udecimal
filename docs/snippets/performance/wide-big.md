| Case | udecimal | shopspring | ericlagergren | apd | govalues |
|---|---|---|---|---|---|
| `ParseWide` `123456789123456789123456.1234567890123456` | 290 ns / 6 | 284 ns / 5 | 394 ns / 4 | 382 ns / 6 | can't represent |
| `ParseWide` `111222333444555666777888999.1234567890123456789` | 292 ns / 6 | 319 ns / 5 | 440 ns / 4 | 460 ns / 6 | can't represent |
| `StringWide` `123456789123456789123456.1234567890123456` | 123 ns / 4 | 122 ns / 4 | 143 ns / 6 | 159 ns / 6 | can't represent |
| `StringWide` `111222333444555666777888999.1234567890123456789` | 125 ns / 4 | 122 ns / 4 | 145 ns / 6 | 164 ns / 6 | can't represent |
| `DivWide` `12345679012345679890123456789.1234567890123456789.Div(999999)` | 71.5 ns / 6 | 78.0 ns / 7 | 62.1 ns / 2 | 260 ns / 9 | can't represent |
| `DivWide` `1234.Div(12345679012345679890123456789.1234567890123456789)` | 74.3 ns / 7 | 65.0 ns / 5 | 152 ns / 4 | 287 ns / 9 | can't represent |

Each cell shows the median ns/op and allocs/op.
