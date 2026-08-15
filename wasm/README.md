# kiwi-go WebAssembly

This directory provides WebAssembly bindings for `kiwi-go`. It allows using the fast Cassowary constraint solver directly in browser environments (or NodeJS).

## Build

To build the WASM file, simply run the build script:

```sh
./build.sh
```

This will output `kiwi.wasm` and copy `wasm_exec.js` from your Go installation. For the smallest bundle size, it is recommended to have `wasm-opt` (from [binaryen](https://github.com/WebAssembly/binaryen)) installed in your system PATH.

## Test

To run the internal bridge tests via NodeJS, use standard Go tools:

```sh
GOOS=js GOARCH=wasm go test -v .
```
