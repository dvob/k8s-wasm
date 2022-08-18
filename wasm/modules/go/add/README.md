# add

This example uses [TinyGo](https://tinygo.org/) to build a WASM module from Go code.

* Build:
```
tinygo build -target wasm
```

* Run:
```
wasmtime add.wasm --invoke add 3 5
```
