// kiwi.js - Ergonomic JavaScript wrapper for kiwi-go WebAssembly

// Registry to automatically clean up Go memory when JS objects are garbage collected.
const registry = new FinalizationRegistry((held) => {
    // If Go WASM was unloaded or died, the global might be gone
    if (typeof globalThis.kiwi_freeSolver !== 'function') return;

    if (held.type === 'variable') globalThis.kiwi_freeVariable(held.id);
    else if (held.type === 'constraint') globalThis.kiwi_freeConstraint(held.id);
    else if (held.type === 'solver') globalThis.kiwi_freeSolver(held.id);
});

export const Operator = {
    Le: 0,
    Ge: 1,
    Eq: 2
};

export const Strength = {
    create: (a, b, c, w = 1.0) => {
        let res = 0;
        res += Math.max(0, Math.min(1000, a * w)) * 1000000;
        res += Math.max(0, Math.min(1000, b * w)) * 1000;
        res += Math.max(0, Math.min(1000, c * w));
        return res;
    }
};

Strength.required = Strength.create(1000, 1000, 1000);
Strength.strong = Strength.create(1, 0, 0);
Strength.medium = Strength.create(0, 1, 0);
Strength.weak = Strength.create(0, 0, 1);

export class Variable {
    constructor(name = "") {
        this._name = name;
        this.id = globalThis.kiwi_createVariable(name);
        registry.register(this, { type: 'variable', id: this.id });
    }

    name() { return this._name; }
    value() { return globalThis.kiwi_getVariableValue(this.id); }

    plus(value) { return new Expression(this, value); }
    minus(value) { return new Expression(this, typeof value === 'number' ? -value : [-1, value]); }
    multiply(coeff) { return new Expression([coeff, this]); }
    divide(coeff) { return new Expression([1 / coeff, this]); }
}

export class Expression {
    constructor(...args) {
        this.terms = new Map();
        let constantObj = { val: 0 };
        this._parseArgs(args, this.terms, constantObj);
        this.constant = constantObj.val;
    }

    _parseArgs(args, terms, constant) {
        for (let item of args) {
            if (typeof item === 'number') {
                constant.val += item;
            } else if (item instanceof Variable) {
                terms.set(item, (terms.get(item) || 0) + 1);
            } else if (item instanceof Expression) {
                constant.val += item.constant;
                for (let [v, c] of item.terms) {
                    terms.set(v, (terms.get(v) || 0) + c);
                }
            } else if (Array.isArray(item) && item.length === 2) {
                let coeff = item[0];
                let target = item[1];
                if (typeof coeff !== 'number') throw new Error("array item 0 must be a number");

                if (target instanceof Variable) {
                    terms.set(target, (terms.get(target) || 0) + coeff);
                } else if (target instanceof Expression) {
                    constant.val += target.constant * coeff;
                    for (let [v, c] of target.terms) {
                        terms.set(v, (terms.get(v) || 0) + c * coeff);
                    }
                } else {
                    throw new Error("array item 1 must be a variable or expression");
                }
            } else {
                throw new Error("invalid Expression argument");
            }
        }
    }

    value() {
        let res = this.constant;
        for (let [v, c] of this.terms) {
            res += v.value() * c;
        }
        return res;
    }

    plus(value) { return new Expression(this, value); }
    minus(value) { return new Expression(this, typeof value === 'number' ? -value : [-1, value]); }
    multiply(coeff) { return new Expression([coeff, this]); }
    divide(coeff) { return new Expression([1 / coeff, this]); }
}

export class Constraint {
    constructor(expr, op, rhs, strength = Strength.required) {
        let finalExpr = new Expression(expr);
        if (rhs !== undefined && rhs !== null) {
            finalExpr = finalExpr.minus(rhs);
        }

        let args = [op, strength, finalExpr.constant];
        for (let [v, c] of finalExpr.terms) {
            args.push(v.id, c);
        }

        this.id = globalThis.kiwi_createConstraint(...args);
        registry.register(this, { type: 'constraint', id: this.id });
    }
}

export class Solver {
    constructor() {
        this.id = globalThis.kiwi_createSolver();
        registry.register(this, { type: 'solver', id: this.id });
    }

    addConstraint(cn) {
        let err = globalThis.kiwi_addConstraint(this.id, cn.id);
        if (err) throw new Error(err);
    }

    removeConstraint(cn) {
        let err = globalThis.kiwi_removeConstraint(this.id, cn.id);
        if (err) throw new Error(err);
    }

    addEditVariable(v, strength) {
        let err = globalThis.kiwi_addEditVariable(this.id, v.id, strength);
        if (err) throw new Error(err);
    }

    removeEditVariable(v) {
        let err = globalThis.kiwi_removeEditVariable(this.id, v.id);
        if (err) throw new Error(err);
    }

    suggestValue(v, val) {
        let err = globalThis.kiwi_suggestValue(this.id, v.id, val);
        if (err) throw new Error(err);
    }

    updateVariables() {
        globalThis.kiwi_updateVariables(this.id);
    }
}

let wasmInstantiated = false;

/**
 * Initialize the WASM module.
 * @param {Response | ArrayBuffer} wasmSource - A Response object from fetch(), or an ArrayBuffer containing the WASM binary.
 */
export async function initWasm(wasmSource) {
    if (wasmInstantiated) return;
    if (typeof globalThis.Go === 'undefined') {
        throw new Error("wasm_exec.js must be loaded globally before calling initWasm. You can load it via a <script> tag or require().");
    }

    const go = new globalThis.Go();
    let inst;

    // Use fast instantiateStreaming in browser if possible
    if (typeof Response !== 'undefined' && wasmSource instanceof Response) {
        const result = await WebAssembly.instantiateStreaming(wasmSource, go.importObject);
        inst = result.instance;
    } else {
        // Fallback for ArrayBuffer/Buffer (e.g., NodeJS)
        const result = await WebAssembly.instantiate(wasmSource, go.importObject);
        inst = result.instance;
    }

    // Do not await this, it blocks because Go main() waits on a channel indefinitely
    go.run(inst);
    wasmInstantiated = true;
}
