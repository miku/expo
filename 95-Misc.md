# Misc

## SWAR

* simd with a register
* see: [x/swar](x/swar)

Turn a bytes slice to a uint64:

```
data := *(*uint64)(unsafe.Add(unsafe.Pointer(bytes), offset))
```

Read out ascii chars from memory ([gist](https://gist.github.com/miku/9e02b083d0dcf45e12fffe8ea3cb9eec)):

```go
package main

import (
        "fmt"
        "log"
        "unsafe"
)

func main() {
        defer func() {
                if err := recover(); err != nil {
                        log.Println(err)
                }
        }()
        s := "hi!"
        b := unsafe.StringData(s) // *uint8
        i := 0
        for {
                x := unsafe.Add(unsafe.Pointer(b), i) // pointer arithmetic
                c := *(*uint8)(x)
                if c <= 122 && c >= 65 {
                        fmt.Printf("%c", c) // 0x4ba06b
                }
                i++
        }
        fmt.Println()
        // 00843609
        // unexpected fault address 0x578000
}
```

### Some notions

ASCII occupies the lower 7 bits, so bitwise `AND` with `0x80` will be greater 0, if value was not ASCII.

```
const asciiMask = 0x8080808080808080 // 8 bytes
```

Other:

```
const BroadcastSemicolon = 0x3B3B3B3B3B3B3B3B
const Broadcast0x01 = 0x0101010101010101
const Broadcast0x80 = 0x8080808080808080
```

Semi:

```
In [1]: chr(0x3b)
Out[1]: ';'
```




* [Faster utf8.Valid using multi-byte processing without SIMD.](https://sugawarayuuta.github.io/charcoal/)

## Readings

* [Validating UTF-8 In Less Than One Instruction Per Byte](https://arxiv.org/pdf/2010.03090)
