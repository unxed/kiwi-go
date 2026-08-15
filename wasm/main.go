//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/unxed/kiwi-go"
)

// Maps to hold instances across the WebAssembly boundary
var (
	solvers      = make(map[int]*kiwi.Solver)
	variables    = make(map[int]*kiwi.Variable)
	constraints  = make(map[int]*kiwi.Constraint)
	nextSolverID = 1
)

func jsCreateSolver(this js.Value, args []js.Value) any {
	id := nextSolverID
	nextSolverID++
	solvers[id] = kiwi.NewSolver()
	return id
}

func jsFreeSolver(this js.Value, args []js.Value) any {
	id := args[0].Int()
	delete(solvers, id)
	return nil
}

func jsCreateVariable(this js.Value, args []js.Value) any {
	name := ""
	if len(args) > 0 && args[0].Type() == js.TypeString {
		name = args[0].String()
	}
	v := kiwi.NewVariable(name)
	variables[v.ID()] = v
	return v.ID()
}

func jsFreeVariable(this js.Value, args []js.Value) any {
	id := args[0].Int()
	delete(variables, id)
	return nil
}

func jsGetVariableValue(this js.Value, args []js.Value) any {
	id := args[0].Int()
	if v, ok := variables[id]; ok {
		return v.Value()
	}
	return 0.0
}

func jsSetVariableValue(this js.Value, args []js.Value) any {
	id := args[0].Int()
	val := args[1].Float()
	if v, ok := variables[id]; ok {
		v.SetValue(val)
	}
	return nil
}

func jsCreateConstraint(this js.Value, args []js.Value) any {
	op := kiwi.Operator(args[0].Int())
	strength := args[1].Float()
	constant := args[2].Float()

	exprArgs := make([]any, 0)
	if constant != 0 {
		exprArgs = append(exprArgs, constant)
	}

	// Terms are passed as flat pairs: [..., varId, coeff, varId, coeff]
	termsLen := len(args)
	for i := 3; i < termsLen; i += 2 {
		varId := args[i].Int()
		coeff := args[i+1].Float()
		if v, ok := variables[varId]; ok {
			exprArgs = append(exprArgs, []any{coeff, v})
		}
	}

	expr := kiwi.NewExpression(exprArgs...)
	cn := kiwi.NewConstraint(expr, op, nil, strength)
	constraints[cn.ID()] = cn
	return cn.ID()
}

func jsFreeConstraint(this js.Value, args []js.Value) any {
	id := args[0].Int()
	delete(constraints, id)
	return nil
}

func jsAddConstraint(this js.Value, args []js.Value) any {
	solverId := args[0].Int()
	cnId := args[1].Int()
	if s, sok := solvers[solverId]; sok {
		if cn, cok := constraints[cnId]; cok {
			err := s.AddConstraint(cn)
			if err != nil {
				return err.Error()
			}
			return nil
		}
	}
	return "solver or constraint not found"
}

func jsRemoveConstraint(this js.Value, args []js.Value) any {
	solverId := args[0].Int()
	cnId := args[1].Int()
	if s, sok := solvers[solverId]; sok {
		if cn, cok := constraints[cnId]; cok {
			err := s.RemoveConstraint(cn)
			if err != nil {
				return err.Error()
			}
			return nil
		}
	}
	return "solver or constraint not found"
}

func jsAddEditVariable(this js.Value, args []js.Value) any {
	solverId := args[0].Int()
	varId := args[1].Int()
	strength := args[2].Float()
	if s, sok := solvers[solverId]; sok {
		if v, vok := variables[varId]; vok {
			err := s.AddEditVariable(v, strength)
			if err != nil {
				return err.Error()
			}
			return nil
		}
	}
	return "solver or variable not found"
}

func jsRemoveEditVariable(this js.Value, args []js.Value) any {
	solverId := args[0].Int()
	varId := args[1].Int()
	if s, sok := solvers[solverId]; sok {
		if v, vok := variables[varId]; vok {
			err := s.RemoveEditVariable(v)
			if err != nil {
				return err.Error()
			}
			return nil
		}
	}
	return "solver or variable not found"
}

func jsSuggestValue(this js.Value, args []js.Value) any {
	solverId := args[0].Int()
	varId := args[1].Int()
	val := args[2].Float()
	if s, sok := solvers[solverId]; sok {
		if v, vok := variables[varId]; vok {
			err := s.SuggestValue(v, val)
			if err != nil {
				return err.Error()
			}
			return nil
		}
	}
	return "solver or variable not found"
}

func jsUpdateVariables(this js.Value, args []js.Value) any {
	solverId := args[0].Int()
	if s, sok := solvers[solverId]; sok {
		s.UpdateVariables()
	}
	return nil
}

func registerCallbacks() {
	js.Global().Set("kiwi_createSolver", js.FuncOf(jsCreateSolver))
	js.Global().Set("kiwi_freeSolver", js.FuncOf(jsFreeSolver))
	js.Global().Set("kiwi_createVariable", js.FuncOf(jsCreateVariable))
	js.Global().Set("kiwi_freeVariable", js.FuncOf(jsFreeVariable))
	js.Global().Set("kiwi_getVariableValue", js.FuncOf(jsGetVariableValue))
	js.Global().Set("kiwi_setVariableValue", js.FuncOf(jsSetVariableValue))
	js.Global().Set("kiwi_createConstraint", js.FuncOf(jsCreateConstraint))
	js.Global().Set("kiwi_freeConstraint", js.FuncOf(jsFreeConstraint))
	js.Global().Set("kiwi_addConstraint", js.FuncOf(jsAddConstraint))
	js.Global().Set("kiwi_removeConstraint", js.FuncOf(jsRemoveConstraint))
	js.Global().Set("kiwi_addEditVariable", js.FuncOf(jsAddEditVariable))
	js.Global().Set("kiwi_removeEditVariable", js.FuncOf(jsRemoveEditVariable))
	js.Global().Set("kiwi_suggestValue", js.FuncOf(jsSuggestValue))
	js.Global().Set("kiwi_updateVariables", js.FuncOf(jsUpdateVariables))
}

func main() {
	registerCallbacks()
	// Block main so that the exported functions remain available in WebAssembly
	<-make(chan struct{})
}
