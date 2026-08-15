#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "Building kiwi.wasm..."
GOOS=js GOARCH=wasm go build -o kiwi.wasm .

echo "Copying wasm_exec.js from Go standard library..."
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

if command -v wasm-opt >/dev/null 2>&1; then
    echo "Optimizing with wasm-opt..."
    wasm-opt -O3 kiwi.wasm -o kiwi_opt.wasm
    mv kiwi_opt.wasm kiwi.wasm
    echo "Optimization complete."
else
    echo "Notice: wasm-opt not found, skipping optimization."
    echo "Install binaryen (https://github.com/WebAssembly/binaryen) for smaller bundle size."
fi

echo "Build complete. Wasm file is ready."
