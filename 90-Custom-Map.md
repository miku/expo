# Custom Map

* collision free, if keys are known
* example inspired by [FNV-1](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function)

> The Python programming language previously used a modified version of the FNV
> scheme for its default hash function. From Python 3.4, FNV has been replaced
> with SipHash to resist "hash flooding" denial-of-service attacks.

> One of FNV's key advantages is that it is very simple to implement.[8] Start
> with an initial hash value of FNV offset basis. For each byte in the input,
> multiply hash by the FNV prime, then XOR it with the byte from the input. The
> alternate algorithm, FNV-1a, reverses the multiply and XOR steps.

see:

* [x/xmap](x/xmap)

## Misc

* https://nnethercote.github.io/2021/12/08/a-brutally-effective-hash-function-in-rust.html
* https://github.com/JuliaLang/julia/issues/52440
