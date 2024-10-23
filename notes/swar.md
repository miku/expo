# simd

what is simd? single instruction multiple data.

* sisd - single instruction single data
* simd: 1 instruction, 4 operands

SIMD requires special HW, AVX, NEON.

* with SWAR no extra HW, just broad enough registers
* cheapest way is SWAR

> CUDA/HW-SIMD/SW-SIMD

Pretend to have lanes:

```
constexpr u64 a = 0001'0001; // 1 | 1
constexpr u64 b = 0010'0010; // 2 | 2
constexpr u64 c = 0011'0011; // 3 | 3
```

----

Example applications:

* determine if a string is just ASCII
* find specific char in string, e.g. ";" // 3B
* add 4x 8 bit numbers at once
