import fs from 'fs';
import assert from 'assert';
import crypto from 'crypto';
import { performance } from 'perf_hooks';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

const __dirname = dirname(fileURLToPath(import.meta.url));

// Polyfill globalThis.crypto and performance for NodeJS environments
if (!globalThis.crypto) {
    globalThis.crypto = crypto.webcrypto || crypto;
}
if (!globalThis.performance) {
    globalThis.performance = performance;
}

// Execute wasm_exec.js in global context to define globalThis.Go
const wasmExecCode = fs.readFileSync(join(__dirname, 'wasm_exec.js'), 'utf8');
(new Function(wasmExecCode))();

// Dynamically import kiwi.js after Go polyfill is initialized
const { initWasm, Solver, Variable, Constraint, Operator, Strength } = await import('./kiwi.js');

async function test() {
    console.log("Loading WASM module...");
    const wasmBuffer = fs.readFileSync(join(__dirname, 'kiwi.wasm'));
    await initWasm(wasmBuffer);

    console.log("Initializing solver...");
    const solver = new Solver();
    const left = new Variable("left");
    const width = new Variable("width");
    const right = new Variable("right");

    solver.addEditVariable(left, Strength.strong);
    solver.addEditVariable(width, Strength.strong);

    solver.suggestValue(left, 100);
    solver.suggestValue(width, 400);

    // right == left + width
    console.log("Adding constraints...");
    const cn = new Constraint(right, Operator.Eq, left.plus(width));
    solver.addConstraint(cn);

    solver.updateVariables();

    console.log("Checking results...");
    assert.strictEqual(right.value(), 500, "right should be 500");
    console.log("✅ Integration test passed: right == 500");

    // Change value
    solver.suggestValue(width, 50);
    solver.updateVariables();
    assert.strictEqual(right.value(), 150, "right should be 150");
    console.log("✅ Integration test passed: right == 150 after width change");
}

test().catch(err => {
    console.error("❌ Test failed:", err);
    process.exit(1);
});
