#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "Building kiwi.wasm..."
GOOS=js GOARCH=wasm go build -o kiwi.wasm .

echo "Locating wasm_exec.js..."
WASM_EXEC_PATH=""

# Check both new (Go 1.24+) and old (Go <= 1.23) standard paths
PATHS=(
    "$(go env GOROOT)/lib/wasm/wasm_exec.js"
    "$(go env GOROOT)/misc/wasm/wasm_exec.js"
)

for p in "${PATHS[@]}"; do
    if [ -f "$p" ]; then
        WASM_EXEC_PATH="$p"
        break
    fi
done

if [ -n "$WASM_EXEC_PATH" ]; then
    echo "Copying wasm_exec.js from $WASM_EXEC_PATH..."
    cp "$WASM_EXEC_PATH" .
else
    echo "Notice: wasm_exec.js not found in local GOROOT (common for some stripped Go toolchains)."

    # Attempt to download from upstream Go matching the active version
    GO_VER=$(go version | awk '{print $3}')
    GO_VER_CLEAN=$(echo "$GO_VER" | grep -oE 'go[0-9]+\.[0-9]+(\.[0-9]+)?' || echo "master")

    # Determine directory structure based on major version (Go 1.24+ uses lib/wasm)
    MAJOR_VER=$(echo "$GO_VER_CLEAN" | cut -d. -f2 || echo "24")
    REMOTE_DIR="misc/wasm"
    if [ "$MAJOR_VER" -ge 24 ] 2>/dev/null || [ "$GO_VER_CLEAN" = "master" ]; then
        REMOTE_DIR="lib/wasm"
    fi

    URL="https://raw.githubusercontent.com/golang/go/${GO_VER_CLEAN}/${REMOTE_DIR}/wasm_exec.js"
    echo "Attempting to download matching version from upstream Go: $URL"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$URL" -o wasm_exec.js && echo "Successfully downloaded wasm_exec.js."
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$URL" -O wasm_exec.js && echo "Successfully downloaded wasm_exec.js."
    else
        echo "Error: Could not copy or download wasm_exec.js (neither curl nor wget is installed)."
        echo "Please download it manually from $URL and place it in this directory."
        exit 1
    fi
fi

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
