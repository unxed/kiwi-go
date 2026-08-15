package kiwi

func nearZero(value float64) bool {
	eps := 1.0e-8
	if value < 0.0 {
		return -value < eps
	}
	return value < eps
}

// Row represents an internal equation row used by the solver.
type Row struct {
	cellMap  *IndexedMap[*Symbol, float64]
	constant float64
}

// NewRow creates a new Row instance with an optional initial constant value.
func NewRow(constant ...float64) *Row {
	c := 0.0
	if len(constant) > 0 {
		c = constant[0]
	}
	return &Row{
		cellMap:  NewIndexedMap[*Symbol, float64](),
		constant: c,
	}
}

// Cells returns the mapping of symbols to coefficients.
func (r *Row) Cells() *IndexedMap[*Symbol, float64] {
	return r.cellMap
}

// Constant returns the constant for the row.
func (r *Row) Constant() float64 {
	return r.constant
}

// IsConstant returns true if the row is a constant value.
func (r *Row) IsConstant() bool {
	return r.cellMap.Empty()
}

// AllDummies returns true if the Row contains only dummy symbols.
func (r *Row) AllDummies() bool {
	cells := r.cellMap
	for i := 0; i < cells.Size(); i++ {
		pair := cells.ItemAt(i)
		if pair.First.Type() != SymbolDummy {
			return false
		}
	}
	return true
}

// Copy creates a deep copy of the row.
func (r *Row) Copy() *Row {
	theCopy := NewRow(r.constant)
	theCopy.cellMap = r.cellMap.Copy(nil)
	return theCopy
}

// Add adds a value to the row constant and returns the new value.
func (r *Row) Add(value float64) float64 {
	r.constant += value
	return r.constant
}

// InsertSymbol inserts a symbol into the row with a given coefficient.
func (r *Row) InsertSymbol(symbol *Symbol, coefficient ...float64) {
	coeff := 1.0
	if len(coefficient) > 0 {
		coeff = coefficient[0]
	}
	factory := func() float64 { return 0.0 }
	pair := r.cellMap.SetDefault(symbol, factory)
	pair.Second += coeff
	if nearZero(pair.Second) {
		r.cellMap.Erase(symbol)
	}
}

// InsertRow inserts a row into this row multiplied by a given coefficient.
func (r *Row) InsertRow(other *Row, coefficient ...float64) {
	coeff := 1.0
	if len(coefficient) > 0 {
		coeff = coefficient[0]
	}
	r.constant += other.constant * coeff
	cells := other.cellMap
	for i := 0; i < cells.Size(); i++ {
		pair := cells.ItemAt(i)
		r.InsertSymbol(pair.First, pair.Second*coeff)
	}
}

// RemoveSymbol removes a symbol from the row.
func (r *Row) RemoveSymbol(symbol *Symbol) {
	r.cellMap.Erase(symbol)
}

// ReverseSign reverses the sign of the constant and all cell coefficients in the row.
func (r *Row) ReverseSign() {
	r.constant = -r.constant
	cells := r.cellMap
	for i := 0; i < cells.Size(); i++ {
		pair := cells.ItemAtPtr(i)
		pair.Second = -pair.Second
	}
}

// SolveFor solves the row for the given symbol.
func (r *Row) SolveFor(symbol *Symbol) {
	pair, ok := r.cellMap.Erase(symbol)
	if !ok {
		return
	}
	coeff := -1.0 / pair.Second
	r.constant *= coeff
	cells := r.cellMap
	for i := 0; i < cells.Size(); i++ {
		p := cells.ItemAtPtr(i)
		p.Second *= coeff
	}
}

// SolveForEx solves the row for the given LHS and RHS symbols.
func (r *Row) SolveForEx(lhs *Symbol, rhs *Symbol) {
	r.InsertSymbol(lhs, -1.0)
	r.SolveFor(rhs)
}

// CoefficientFor returns the coefficient for the given symbol.
func (r *Row) CoefficientFor(symbol *Symbol) float64 {
	pair, ok := r.cellMap.Find(symbol)
	if !ok {
		return 0.0
	}
	return pair.Second
}

// Substitute substitutes a symbol with the data from another row.
func (r *Row) Substitute(symbol *Symbol, row *Row) {
	pair, ok := r.cellMap.Erase(symbol)
	if ok {
		r.InsertRow(row, pair.Second)
	}
}
