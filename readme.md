# 🌽

Milho is a lisp-ish language created by [@celsobonutti](https://github.com/celsobonutti) and me ([@danfragoso](https://github.com/danfragoso)).

This repo hosts the GO implementation of the Milho interpreter and compiler. There's also [Rust](https://github.com/celsobonutti/milho-rust), [Haskell](https://github.com/celsobonutti/milho-rust) and [OCaml](https://github.com/renatoalencar/milho-ocaml) implementations.

## Try it online (WebAssembly version)
[https://milho.fragoso.dev](https://milho.fragoso.dev)

## How does it work?
```
Tokenizer → Parser → MIR
                      ├ → Interpreter (Linux, BSD, WebAssembly)
                      │                or anything you manage to run GO on.
                      │
                      └ → Compiler/Codegen
                            ├ → JavaScript source code
                            └ → LLVM IR code
```

## Building
```
make build
```

## Installing
```
make install
```
## Running a milho file
```
milho hello.milho
```
## Creating scripts
If you want to run milho files from the command line, you can create a milho script.
You just need to add ```#!/bin/milho``` to the top of your file and run with ```./script.milho```.

```clojure
#!/bin/milho

(defn getUserName ()
  (exec "git" "config" "--get" "user.name"))

(defn buildHelloString (who)
  (str "Hello " who "! 🌽"))

(def userName (getUserName))

(if (! (= userName ""))
  (println (buildHelloString userName))
  (exit 1)) ;; Exit with error if no user.name.

```
## Compiling a milho file
```
// To JavaScript
milho -cJS examples/hello.milho > hello.js
node hello.js
```
```
// To LLVM IR
milho -cLLVM examples/hello.milho > hello.ll
lli hello.ll
```
## Configuring git hooks
```
make hooks
```

## Running the repl
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