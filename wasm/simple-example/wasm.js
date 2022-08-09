const fs = require("fs");

const wasmModule = new WebAssembly.Module(fs.readFileSync("add.wasm"));

(async () => {
	const instance = await WebAssembly.instantiate(wasmModule, {})
	console.log(instance.exports.add(11,31))
})();
