package kiwi

// Identifiable is an interface for types that have an integer ID().
type Identifiable interface {
	ID() int
}

// Pair defines a generic pair of key and value.
type Pair[T Identifiable, U any] struct {
	First  T
	Second U
}

// IndexedMap is an associative array maintaining key order and fast lookup by ID().
type IndexedMap[T Identifiable, U any] struct {
	index map[int]int
	array []Pair[T, U]
}

// NewIndexedMap creates a new IndexedMap instance.
func NewIndexedMap[T Identifiable, U any]() *IndexedMap[T, U] {
	return &IndexedMap[T, U]{
		index: make(map[int]int),
		array: make([]Pair[T, U], 0),
	}
}

// Size returns the number of items in the map.
func (m *IndexedMap[T, U]) Size() int {
	return len(m.array)
}

// Empty returns true if the map is empty.
func (m *IndexedMap[T, U]) Empty() bool {
	return len(m.array) == 0
}

// ItemAt returns the pair at the given array index.
func (m *IndexedMap[T, U]) ItemAt(index int) Pair[T, U] {
	return m.array[index]
}

// ItemAtPtr returns a pointer to the pair at the given array index.
func (m *IndexedMap[T, U]) ItemAtPtr(index int) *Pair[T, U] {
	return &m.array[index]
}

// Array returns the slice of pairs.
func (m *IndexedMap[T, U]) Array() []Pair[T, U] {
	return m.array
}

// Contains returns true if the key exists in the map.
func (m *IndexedMap[T, U]) Contains(key T) bool {
	_, ok := m.index[key.ID()]
	return ok
}

// Find returns a pointer to the pair for key and true if found, nil and false otherwise.
func (m *IndexedMap[T, U]) Find(key T) (*Pair[T, U], bool) {
	idx, ok := m.index[key.ID()]
	if !ok {
		return nil, false
	}
	return &m.array[idx], true
}

// SetDefault returns the pair associated with key.
// If the key does not exist, a new pair is created using the factory function.
func (m *IndexedMap[T, U]) SetDefault(key T, factory func() U) *Pair[T, U] {
	idx, ok := m.index[key.ID()]
	if ok {
		return &m.array[idx]
	}
	p := Pair[T, U]{First: key, Second: factory()}
	m.index[key.ID()] = len(m.array)
	m.array = append(m.array, p)
	return &m.array[len(m.array)-1]
}

// Insert inserts or overwrites the pair for key.
func (m *IndexedMap[T, U]) Insert(key T, value U) *Pair[T, U] {
	idx, ok := m.index[key.ID()]
	if ok {
		m.array[idx] = Pair[T, U]{First: key, Second: value}
		return &m.array[idx]
	}
	p := Pair[T, U]{First: key, Second: value}
	m.index[key.ID()] = len(m.array)
	m.array = append(m.array, p)
	return &m.array[len(m.array)-1]
}

// Erase removes and returns the pair for key.
func (m *IndexedMap[T, U]) Erase(key T) (Pair[T, U], bool) {
	idx, ok := m.index[key.ID()]
	if !ok {
		var zero Pair[T, U]
		return zero, false
	}
	delete(m.index, key.ID())
	pair := m.array[idx]
	lastIdx := len(m.array) - 1
	last := m.array[lastIdx]
	m.array = m.array[:lastIdx]
	if idx != lastIdx {
		m.array[idx] = last
		m.index[last.First.ID()] = idx
	}
	return pair, true
}

// Copy creates a copy of the IndexedMap.
func (m *IndexedMap[T, U]) Copy(copyVal func(U) U) *IndexedMap[T, U] {
	cp := NewIndexedMap[T, U]()
	cp.array = make([]Pair[T, U], len(m.array))
	for i, p := range m.array {
		val := p.Second
		if copyVal != nil {
			val = copyVal(p.Second)
		}
		cp.array[i] = Pair[T, U]{First: p.First, Second: val}
		cp.index[p.First.ID()] = i
	}
	return cp
}
