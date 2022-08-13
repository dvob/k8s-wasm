# Simple WASM Example

* Compile WASM module `add.wasm`
```
clang --target=wasm32 --no-standard-libraries -Wl,--no-entry -Wl,--export=add -o add.wasm add.c
```

* Use compiled `add.wasm` from JavaScript
```
node wasm.js
```
