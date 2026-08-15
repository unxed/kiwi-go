//go:build js && wasm

package main

import (
	"syscall/js"
	"testing"
)

func TestWASMBridge(t *testing.T) {
	registerCallbacks()

	// 1. Test create solver
	solverIdVal := js.Global().Call("kiwi_createSolver")
	solverId := solverIdVal.Int()
	if solvers[solverId] == nil {
		t.Fatal("solver not created")
	}

	// 2. Test create variable
	varIdVal := js.Global().Call("kiwi_createVariable", "x")
	varId := varIdVal.Int()
	if variables[varId] == nil {
		t.Fatal("variable not created")
	}

	// 3. Test create constraint: x >= 10 (x - 10 >= 0)
	// args: op(1=Ge), strength(1e9), constant(-10), term1VarId(varId), term1Coeff(1.0)
	cnIdVal := js.Global().Call("kiwi_createConstraint", 1, 1000000000.0, -10.0, varId, 1.0)
	cnId := cnIdVal.Int()
	if constraints[cnId] == nil {
		t.Fatal("constraint not created")
	}

	// 4. Test add constraint to solver
	errVal := js.Global().Call("kiwi_addConstraint", solverId, cnId)
	if !errVal.IsNull() {
		t.Fatalf("addConstraint failed: %s", errVal.String())
	}

	// 5. Update and assert value
	js.Global().Call("kiwi_updateVariables", solverId)
	val := js.Global().Call("kiwi_getVariableValue", varId).Float()
	if val != 10.0 {
		t.Fatalf("expected 10.0, got %f", val)
	}

	// 6. Cleanup & check for leaks
	js.Global().Call("kiwi_freeSolver", solverId)
	js.Global().Call("kiwi_freeVariable", varId)
	js.Global().Call("kiwi_freeConstraint", cnId)

	if solvers[solverId] != nil || variables[varId] != nil || constraints[cnId] != nil {
		t.Fatal("memory leak after free")
	}
}
