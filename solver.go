package kiwi

import (
	"errors"
	"math"
)

// Tag holds the marker and error symbols associated with a constraint.
type Tag struct {
	Marker *Symbol
	Other  *Symbol
}

// EditInfo holds metadata for an edit variable in the solver.
type EditInfo struct {
	Tag        Tag
	Constraint *Constraint
	Constant   float64
}

type rowCreationResult struct {
	row *Row
	tag Tag
}

// Solver is the primary constraint solver class using the Cassowary algorithm.
type Solver struct {
	MaxIterations  int
	cnMap          *IndexedMap[*Constraint, Tag]
	rowMap         *IndexedMap[*Symbol, *Row]
	varMap         *IndexedMap[*Variable, *Symbol]
	editMap        *IndexedMap[*Variable, EditInfo]
	infeasibleRows []*Symbol
	objective      *Row
	artificial     *Row
	idTick         int
}

// NewSolver constructs a new Solver instance with default max iterations (1000).
func NewSolver() *Solver {
	return &Solver{
		MaxIterations:  1000,
		cnMap:          NewIndexedMap[*Constraint, Tag](),
		rowMap:         NewIndexedMap[*Symbol, *Row](),
		varMap:         NewIndexedMap[*Variable, *Symbol](),
		editMap:        NewIndexedMap[*Variable, EditInfo](),
		infeasibleRows: make([]*Symbol, 0),
		objective:      NewRow(),
		artificial:     nil,
		idTick:         0,
	}
}

// CreateConstraint creates and adds a new constraint to the solver.
func (s *Solver) CreateConstraint(lhs any, operator Operator, rhsAndStrength ...any) (*Constraint, error) {
	cn := NewConstraint(lhs, operator, rhsAndStrength...)
	err := s.AddConstraint(cn)
	if err != nil {
		return nil, err
	}
	return cn, nil
}

// MustCreateConstraint creates and adds a new constraint, panicking on error.
func (s *Solver) MustCreateConstraint(lhs any, operator Operator, rhsAndStrength ...any) *Constraint {
	cn, err := s.CreateConstraint(lhs, operator, rhsAndStrength...)
	if err != nil {
		panic(err)
	}
	return cn
}

// AddConstraint adds a constraint to the solver.
func (s *Solver) AddConstraint(constraint *Constraint) error {
	if s.cnMap.Contains(constraint) {
		return errors.New("duplicate constraint")
	}

	data := s.createRow(constraint)
	row := data.row
	tag := data.tag
	subject := s.chooseSubject(row, tag)

	if subject.Type() == SymbolInvalid && row.AllDummies() {
		if !nearZero(row.Constant()) {
			return errors.New("unsatisfiable constraint")
		}
		subject = tag.Marker
	}

	if subject.Type() == SymbolInvalid {
		if !s.addWithArtificialVariable(row) {
			return errors.New("unsatisfiable constraint")
		}
	} else {
		row.SolveFor(subject)
		s.substitute(subject, row)
		s.rowMap.Insert(subject, row)
	}

	s.cnMap.Insert(constraint, tag)
	return s.optimize(s.objective)
}

// RemoveConstraint removes a constraint from the solver.
func (s *Solver) RemoveConstraint(constraint *Constraint) error {
	cnPair, ok := s.cnMap.Erase(constraint)
	if !ok {
		return errors.New("unknown constraint")
	}

	s.removeConstraintEffects(constraint, cnPair.Second)

	marker := cnPair.Second.Marker
	rowPair, ok := s.rowMap.Erase(marker)
	if !ok {
		leaving := s.getMarkerLeavingSymbol(marker)
		if leaving.Type() == SymbolInvalid {
			return errors.New("failed to find leaving row")
		}
		rowPair, ok = s.rowMap.Erase(leaving)
		if !ok {
			return errors.New("failed to find leaving row")
		}
		rowPair.Second.SolveForEx(leaving, marker)
		s.substitute(marker, rowPair.Second)
	}

	return s.optimize(s.objective)
}

// HasConstraint returns true if the solver contains the given constraint.
func (s *Solver) HasConstraint(constraint *Constraint) bool {
	return s.cnMap.Contains(constraint)
}

// GetConstraints returns a slice of all current constraints in the solver.
func (s *Solver) GetConstraints() []*Constraint {
	arr := s.cnMap.Array()
	result := make([]*Constraint, len(arr))
	for i, pair := range arr {
		result[i] = pair.First
	}
	return result
}

// AddEditVariable adds an edit variable to the solver with the given strength.
func (s *Solver) AddEditVariable(variable *Variable, strength float64) error {
	if s.editMap.Contains(variable) {
		return errors.New("duplicate edit variable")
	}
	strength = ClipStrength(strength)
	if strength == StrengthRequired {
		return errors.New("bad required strength")
	}
	expr := NewExpression(variable)
	cn := NewConstraint(expr, OpEq, nil, strength)
	err := s.AddConstraint(cn)
	if err != nil {
		return err
	}
	tagPair, _ := s.cnMap.Find(cn)
	info := EditInfo{
		Tag:        tagPair.Second,
		Constraint: cn,
		Constant:   0.0,
	}
	s.editMap.Insert(variable, info)
	return nil
}

// RemoveEditVariable removes an edit variable from the solver.
func (s *Solver) RemoveEditVariable(variable *Variable) error {
	editPair, ok := s.editMap.Erase(variable)
	if !ok {
		return errors.New("unknown edit variable")
	}
	return s.RemoveConstraint(editPair.Second.Constraint)
}

// HasEditVariable returns true if the solver contains the edit variable.
func (s *Solver) HasEditVariable(variable *Variable) bool {
	return s.editMap.Contains(variable)
}

// SuggestValue suggests a new value for an edit variable.
func (s *Solver) SuggestValue(variable *Variable, value float64) error {
	editPair, ok := s.editMap.Find(variable)
	if !ok {
		return errors.New("unknown edit variable")
	}

	rows := s.rowMap
	info := editPair.Second
	delta := value - info.Constant
	info.Constant = value
	s.editMap.Insert(variable, info)

	marker := info.Tag.Marker
	rowPair, ok := rows.Find(marker)
	if ok {
		if rowPair.Second.Add(-delta) < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, marker)
		}
		return s.dualOptimize()
	}

	other := info.Tag.Other
	rowPair, ok = rows.Find(other)
	if ok {
		if rowPair.Second.Add(delta) < 0.0 {
			s.infeasibleRows = append(s.infeasibleRows, other)
		}
		return s.dualOptimize()
	}

	for i := 0; i < rows.Size(); i++ {
		rPair := rows.ItemAt(i)
		row := rPair.Second
		coeff := row.CoefficientFor(marker)
		if coeff != 0.0 && row.Add(delta*coeff) < 0.0 && rPair.First.Type() != SymbolExternal {
			s.infeasibleRows = append(s.infeasibleRows, rPair.First)
		}
	}
	return s.dualOptimize()
}

// UpdateVariables calculates and updates the computed values of all variables.
func (s *Solver) UpdateVariables() {
	vars := s.varMap
	rows := s.rowMap
	for i := 0; i < vars.Size(); i++ {
		pair := vars.ItemAt(i)
		rowPair, ok := rows.Find(pair.Second)
		if ok {
			pair.First.SetValue(rowPair.Second.Constant())
		} else {
			pair.First.SetValue(0.0)
		}
	}
}

func (s *Solver) makeSymbol(symType SymbolType) *Symbol {
	sym := NewSymbol(symType, s.idTick)
	s.idTick++
	return sym
}

func (s *Solver) getVarSymbol(variable *Variable) *Symbol {
	factory := func() *Symbol {
		return s.makeSymbol(SymbolExternal)
	}
	return s.varMap.SetDefault(variable, factory).Second
}

func (s *Solver) createRow(constraint *Constraint) rowCreationResult {
	expr := constraint.Expression()
	row := NewRow(expr.Constant())

	terms := expr.Terms()
	for i := 0; i < terms.Size(); i++ {
		termPair := terms.ItemAt(i)
		if !nearZero(termPair.Second) {
			symbol := s.getVarSymbol(termPair.First)
			basicPair, ok := s.rowMap.Find(symbol)
			if ok {
				row.InsertRow(basicPair.Second, termPair.Second)
			} else {
				row.InsertSymbol(symbol, termPair.Second)
			}
		}
	}

	objective := s.objective
	strength := constraint.Strength()
	tag := Tag{Marker: invalidSymbol, Other: invalidSymbol}

	switch constraint.Op() {
	case OpLe, OpGe:
		coeff := 1.0
		if constraint.Op() == OpGe {
			coeff = -1.0
		}
		slack := s.makeSymbol(SymbolSlack)
		tag.Marker = slack
		row.InsertSymbol(slack, coeff)
		if strength < StrengthRequired {
			errorSym := s.makeSymbol(SymbolError)
			tag.Other = errorSym
			row.InsertSymbol(errorSym, -coeff)
			objective.InsertSymbol(errorSym, strength)
		}
	case OpEq:
		if strength < StrengthRequired {
			errplus := s.makeSymbol(SymbolError)
			errminus := s.makeSymbol(SymbolError)
			tag.Marker = errplus
			tag.Other = errminus
			row.InsertSymbol(errplus, -1.0)
			row.InsertSymbol(errminus, 1.0)
			objective.InsertSymbol(errplus, strength)
			objective.InsertSymbol(errminus, strength)
		} else {
			dummy := s.makeSymbol(SymbolDummy)
			tag.Marker = dummy
			row.InsertSymbol(dummy)
		}
	}

	if row.Constant() < 0.0 {
		row.ReverseSign()
	}

	return rowCreationResult{row: row, tag: tag}
}

func (s *Solver) chooseSubject(row *Row, tag Tag) *Symbol {
	cells := row.Cells()
	for i := 0; i < cells.Size(); i++ {
		pair := cells.ItemAt(i)
		if pair.First.Type() == SymbolExternal {
			return pair.First
		}
	}
	t := tag.Marker.Type()
	if t == SymbolSlack || t == SymbolError {
		if row.CoefficientFor(tag.Marker) < 0.0 {
			return tag.Marker
		}
	}
	t = tag.Other.Type()
	if t == SymbolSlack || t == SymbolError {
		if row.CoefficientFor(tag.Other) < 0.0 {
			return tag.Other
		}
	}
	return invalidSymbol
}

func (s *Solver) addWithArtificialVariable(row *Row) bool {
	art := s.makeSymbol(SymbolSlack)
	s.rowMap.Insert(art, row.Copy())
	s.artificial = row.Copy()

	err := s.optimize(s.artificial)
	if err != nil {
		s.artificial = nil
		return false
	}
	success := nearZero(s.artificial.Constant())
	s.artificial = nil

	pair, ok := s.rowMap.Erase(art)
	if ok {
		basicRow := pair.Second
		if basicRow.IsConstant() {
			return success
		}
		entering := s.anyPivotableSymbol(basicRow)
		if entering.Type() == SymbolInvalid {
			return false
		}
		basicRow.SolveForEx(art, entering)
		s.substitute(entering, basicRow)
		s.rowMap.Insert(entering, basicRow)
	}

	rows := s.rowMap
	for i := 0; i < rows.Size(); i++ {
		rows.ItemAt(i).Second.RemoveSymbol(art)
	}
	s.objective.RemoveSymbol(art)
	return success
}

func (s *Solver) substitute(symbol *Symbol, row *Row) {
	rows := s.rowMap
	for i := 0; i < rows.Size(); i++ {
		pair := rows.ItemAt(i)
		pair.Second.Substitute(symbol, row)
		if pair.Second.Constant() < 0.0 && pair.First.Type() != SymbolExternal {
			s.infeasibleRows = append(s.infeasibleRows, pair.First)
		}
	}
	s.objective.Substitute(symbol, row)
	if s.artificial != nil {
		s.artificial.Substitute(symbol, row)
	}
}

func (s *Solver) optimize(objective *Row) error {
	iterations := 0
	for iterations < s.MaxIterations {
		entering := s.getEnteringSymbol(objective)
		if entering.Type() == SymbolInvalid {
			return nil
		}
		leaving := s.getLeavingSymbol(entering)
		if leaving.Type() == SymbolInvalid {
			return errors.New("the objective is unbounded")
		}

		pair, ok := s.rowMap.Erase(leaving)
		if !ok {
			return errors.New("leaving row not found")
		}
		row := pair.Second
		row.SolveForEx(leaving, entering)
		s.substitute(entering, row)
		s.rowMap.Insert(entering, row)

		iterations++
	}
	return errors.New("solver iterations exceeded")
}

func (s *Solver) dualOptimize() error {
	rows := s.rowMap
	for len(s.infeasibleRows) != 0 {
		leaving := s.infeasibleRows[len(s.infeasibleRows)-1]
		s.infeasibleRows = s.infeasibleRows[:len(s.infeasibleRows)-1]

		pair, ok := rows.Find(leaving)
		if ok && pair.Second.Constant() < 0.0 {
			entering := s.getDualEnteringSymbol(pair.Second)
			if entering.Type() == SymbolInvalid {
				return errors.New("dual optimize failed")
			}

			row := pair.Second
			rows.Erase(leaving)
			row.SolveForEx(leaving, entering)
			s.substitute(entering, row)
			rows.Insert(entering, row)
		}
	}
	return nil
}

func (s *Solver) getEnteringSymbol(objective *Row) *Symbol {
	cells := objective.Cells()
	for i := 0; i < cells.Size(); i++ {
		pair := cells.ItemAt(i)
		symbol := pair.First
		if pair.Second < 0.0 && symbol.Type() != SymbolDummy {
			return symbol
		}
	}
	return invalidSymbol
}

func (s *Solver) getDualEnteringSymbol(row *Row) *Symbol {
	ratio := math.MaxFloat64
	entering := invalidSymbol
	cells := row.Cells()
	for i := 0; i < cells.Size(); i++ {
		pair := cells.ItemAt(i)
		symbol := pair.First
		c := pair.Second
		if c > 0.0 && symbol.Type() != SymbolDummy {
			coeff := s.objective.CoefficientFor(symbol)
			r := coeff / c
			if r < ratio {
				ratio = r
				entering = symbol
			}
		}
	}
	return entering
}

func (s *Solver) getLeavingSymbol(entering *Symbol) *Symbol {
	ratio := math.MaxFloat64
	found := invalidSymbol
	rows := s.rowMap
	for i := 0; i < rows.Size(); i++ {
		pair := rows.ItemAt(i)
		symbol := pair.First
		if symbol.Type() != SymbolExternal {
			row := pair.Second
			temp := row.CoefficientFor(entering)
			if temp < 0.0 {
				tempRatio := -row.Constant() / temp
				if tempRatio < ratio {
					ratio = tempRatio
					found = symbol
				}
			}
		}
	}
	return found
}

func (s *Solver) getMarkerLeavingSymbol(marker *Symbol) *Symbol {
	dmax := math.MaxFloat64
	r1 := dmax
	r2 := dmax
	first := invalidSymbol
	second := invalidSymbol
	third := invalidSymbol
	rows := s.rowMap

	for i := 0; i < rows.Size(); i++ {
		pair := rows.ItemAt(i)
		row := pair.Second
		c := row.CoefficientFor(marker)
		if c == 0.0 {
			continue
		}
		symbol := pair.First
		if symbol.Type() == SymbolExternal {
			third = symbol
		} else if c < 0.0 {
			r := -row.Constant() / c
			if r < r1 {
				r1 = r
				first = symbol
			}
		} else {
			r := row.Constant() / c
			if r < r2 {
				r2 = r
				second = symbol
			}
		}
	}

	if first != invalidSymbol {
		return first
	}
	if second != invalidSymbol {
		return second
	}
	return third
}

func (s *Solver) removeConstraintEffects(cn *Constraint, tag Tag) {
	if tag.Marker.Type() == SymbolError {
		s.removeMarkerEffects(tag.Marker, cn.Strength())
	}
	if tag.Other.Type() == SymbolError {
		s.removeMarkerEffects(tag.Other, cn.Strength())
	}
}

func (s *Solver) removeMarkerEffects(marker *Symbol, strength float64) {
	pair, ok := s.rowMap.Find(marker)
	if ok {
		s.objective.InsertRow(pair.Second, -strength)
	} else {
		s.objective.InsertSymbol(marker, -strength)
	}
}

func (s *Solver) anyPivotableSymbol(row *Row) *Symbol {
	cells := row.Cells()
	for i := 0; i < cells.Size(); i++ {
		pair := cells.ItemAt(i)
		t := pair.First.Type()
		if t == SymbolSlack || t == SymbolError {
			return pair.First
		}
	}
	return invalidSymbol
}
