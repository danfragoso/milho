# 🌽

Milho is a lisp-ish language created by [@celsobonutti](https://github.com/celsobonutti) and me ([@danfragoso](https://github.com/danfragoso)).

This repo hosts the GO implementation of the Milho interpreter. There's also a Rust implementation at [github.com/celsobonutti/milho](https://github.com/celsobonutti/milho).

## Try it online
[https://milho.fragoso.dev](https://milho.fragoso.dev)
## Running
Running the command bellow will compile and run the milho repl.
```
make repl
```

```
Milho REPL 🌽 v.6f29e9e_2021-04-22
Danilo Fragoso <danilo.fragoso@gmail.com> - 2021
🌽 > (def milho 10)
🍿 milho
🌽 > (def sq_milho (* milho milho))
🍿 sq_milho
🌽 > (prn sq_milho)
100
🍿 Nil
```