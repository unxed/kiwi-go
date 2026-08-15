package kiwi

import "math"

// ApportionGroup defines a group of variables constrained to sum to a target variable or constant.
type ApportionGroup struct {
	Vars        []*Variable
	TargetVar   *Variable
	TargetConst int
	MinSizes    map[*Variable]int
}

// DiscreteSolver integrates a continuous Cassowary Solver with FreeType autohinting heuristics
// and TrueType rule directives for discrete TUI auto-layout.
type DiscreteSolver struct {
	solver     *Solver
	groups     []ApportionGroup
	tracked    []*Variable
	minSizes   map[*Variable]int
	maxSizes   map[*Variable]int
	ruleHinter *RuleHinter
}

// NewDiscreteSolver creates a new DiscreteSolver wrapping a Cassowary Solver instance.
// If solver is nil, a new Solver instance is created automatically.
func NewDiscreteSolver(solver ...*Solver) *DiscreteSolver {
	var s *Solver
	if len(solver) > 0 && solver[0] != nil {
		s = solver[0]
	} else {
		s = NewSolver()
	}
	return &DiscreteSolver{
		solver:     s,
		groups:     make([]ApportionGroup, 0),
		tracked:    make([]*Variable, 0),
		minSizes:   make(map[*Variable]int),
		maxSizes:   make(map[*Variable]int),
		ruleHinter: NewRuleHinter(),
	}
}

// Solver returns the underlying Cassowary continuous Solver instance.
func (ds *DiscreteSolver) Solver() *Solver {
	return ds.solver
}

// RuleHinter returns the internal RuleHinter for adding TrueType-style hint directives.
func (ds *DiscreteSolver) RuleHinter() *RuleHinter {
	return ds.ruleHinter
}

// TrackVariable adds variables to be included in the DiscreteResult.
func (ds *DiscreteSolver) TrackVariable(vars ...*Variable) *DiscreteSolver {
	ds.tracked = append(ds.tracked, vars...)
	return ds
}

// SetMinSize registers a minimum character dimension for a variable.
func (ds *DiscreteSolver) SetMinSize(v *Variable, minSize int) *DiscreteSolver {
	if v != nil {
		ds.minSizes[v] = minSize
		ds.TrackVariable(v)
	}
	return ds
}

// SetMaxSize registers a maximum character dimension for a variable.
func (ds *DiscreteSolver) SetMaxSize(v *Variable, maxSize int) *DiscreteSolver {
	if v != nil {
		ds.maxSizes[v] = maxSize
		ds.TrackVariable(v)
	}
	return ds
}

// AddApportionGroup adds a group of variables that sum to a target variable or constant.
func (ds *DiscreteSolver) AddApportionGroup(group ApportionGroup) *DiscreteSolver {
	ds.groups = append(ds.groups, group)
	for _, v := range group.Vars {
		ds.TrackVariable(v)
	}
	if group.TargetVar != nil {
		ds.TrackVariable(group.TargetVar)
	}
	return ds
}

// AddDirective adds TrueType-style hint directives to the discrete solver pipeline.
func (ds *DiscreteSolver) AddDirective(directives ...HintDirective) *DiscreteSolver {
	ds.ruleHinter.AddDirective(directives...)
	return ds
}

// SolveDiscrete solves continuous constraints and applies autohinting and directives to return integer grid layout values.
func (ds *DiscreteSolver) SolveDiscrete() DiscreteResult {
	// 1. Solve continuous Cassowary constraints
	ds.solver.UpdateVariables()

	result := make(DiscreteResult)
	groupedVars := make(map[*Variable]bool)

	// 2. Process Apportionment Groups (FreeType-style Largest Remainder method)
	for _, group := range ds.groups {
		if len(group.Vars) == 0 {
			continue
		}

		targetSum := group.TargetConst
		if group.TargetVar != nil {
			targetSum = int(math.Round(group.TargetVar.Value()))
		}

		minMap := make(map[*Variable]int)
		for k, v := range ds.minSizes {
			minMap[k] = v
		}
		for k, v := range group.MinSizes {
			minMap[k] = v
		}

		groupRes := ApportionSum(group.Vars, nil, targetSum, minMap)
		for _, v := range group.Vars {
			result[v] = groupRes.Get(v)
			groupedVars[v] = true
		}
		if group.TargetVar != nil {
			result[group.TargetVar] = targetSum
			groupedVars[group.TargetVar] = true
		}
	}

	// 3. Process remaining tracked variables with standard rounding
	var remaining []*Variable
	for _, v := range ds.tracked {
		if !groupedVars[v] {
			remaining = append(remaining, v)
			groupedVars[v] = true
		}
	}

	if len(remaining) > 0 {
		roundRes := RoundValues(remaining, nil, ds.minSizes)
		for _, v := range remaining {
			result[v] = roundRes.Get(v)
		}
	}

	// 4. Apply MaxSize clamping if configured
	for v, maxVal := range ds.maxSizes {
		if val, ok := result[v]; ok && val > maxVal {
			result[v] = maxVal
		}
	}

	// 5. Apply TrueType-style Rule Directives
	ds.ruleHinter.Apply(result)

	return result
}
