export declare enum Operator {
    Le = 0,
    Ge = 1,
    Eq = 2
}

export declare const Strength: {
    create: (a: number, b: number, c: number, w?: number) => number;
    required: number;
    strong: number;
    medium: number;
    weak: number;
};

export declare class Variable {
    constructor(name?: string);
    name(): string;
    value(): number;
    plus(value: any): Expression;
    minus(value: any): Expression;
    multiply(coeff: number): Expression;
    divide(coeff: number): Expression;
}

export declare class Expression {
    constructor(...args: any[]);
    value(): number;
    plus(value: any): Expression;
    minus(value: any): Expression;
    multiply(coeff: number): Expression;
    divide(coeff: number): Expression;
}

export declare class Constraint {
    constructor(expr: any, op: Operator, rhs?: any, strength?: number);
}

export declare class Solver {
    constructor();
    addConstraint(cn: Constraint): void;
    removeConstraint(cn: Constraint): void;
    addEditVariable(v: Variable, strength: number): void;
    removeEditVariable(v: Variable): void;
    suggestValue(v: Variable, val: number): void;
    updateVariables(): void;
}

/**
 * Initializes the Kiwi WebAssembly engine.
 * @param wasmSource A Response (e.g. from fetch('kiwi.wasm')) or an ArrayBuffer containing the WebAssembly binary.
 */
export declare function initWasm(wasmSource: Response | ArrayBuffer): Promise<void>;
