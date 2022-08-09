# Simple WASM Example

* Compile
```
clang --target=wasm32 --no-standard-libraries -Wl,--no-entry -Wl,--export=add -o add.wasm add.c
```

* Run from JavaScript
```
node wasm.js
```
