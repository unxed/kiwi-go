package kiwi_test

import (
	"testing"

	"github.com/lume/kiwi-go"
)

type dummyItem struct {
	id int
}

func (d *dummyItem) ID() int {
	return d.id
}

func TestIndexedMap(t *testing.T) {
	m := kiwi.NewIndexedMap[*dummyItem, string]()
	if !m.Empty() || m.Size() != 0 {
		t.Errorf("expected empty map")
	}

	i1 := &dummyItem{id: 1}
	i2 := &dummyItem{id: 2}

	m.Insert(i1, "first")
	m.Insert(i2, "second")

	if m.Size() != 2 {
		t.Errorf("expected size 2, got %d", m.Size())
	}

	if !m.Contains(i1) || !m.Contains(i2) {
		t.Errorf("expected map to contain i1 and i2")
	}

	p, ok := m.Find(i1)
	if !ok || p.Second != "first" {
		t.Errorf("expected 'first', got %v, ok=%v", p, ok)
	}

	if m.ItemAt(0).Second != "first" {
		t.Errorf("expected ItemAt(0) to be 'first'")
	}
	if m.ItemAtPtr(0).Second != "first" {
		t.Errorf("expected ItemAtPtr(0) to be 'first'")
	}

	// Overwrite i1
	m.Insert(i1, "updated_first")
	p, ok = m.Find(i1)
	if !ok || p.Second != "updated_first" {
		t.Errorf("expected 'updated_first', got %v", p)
	}

	// SetDefault
	i3 := &dummyItem{id: 3}
	p3 := m.SetDefault(i3, func() string { return "third" })
	if p3.Second != "third" || m.Size() != 3 {
		t.Errorf("expected 'third' inserted, size 3")
	}

	// Erase
	erased, ok := m.Erase(i1)
	if !ok || erased.Second != "updated_first" {
		t.Errorf("expected erase i1, got %v, %v", erased, ok)
	}
	if m.Contains(i1) || m.Size() != 2 {
		t.Errorf("expected i1 removed, size 2")
	}

	// Non-existent erase
	_, ok = m.Erase(i1)
	if ok {
		t.Errorf("expected false for erasing non-existent key")
	}

	// Copy
	cp := m.Copy(nil)
	if cp.Size() != 2 || !cp.Contains(i2) || !cp.Contains(i3) {
		t.Errorf("expected copy to have 2 items")
	}
}
