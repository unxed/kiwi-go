// Package kiwi provides a fast, lightweight Cassowary constraint solving algorithm
// ported directly from LUME Kiwi (@lume/kiwi).
//
// Kiwi allows declaring variables, setting linear constraints (equalities and inequalities)
// with symbolic strengths, suggesting edit values, and automatically calculating optimal
// variable values.
//
// Thread Safety:
// Variable instances use sync.RWMutex and are safe for concurrent read and write access
// across multiple goroutines.
package kiwi
