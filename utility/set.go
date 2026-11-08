package utility

/*
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 📦 PACKAGE        : utils                                                                       ║
║ 🎯 WHAT IT DOES   : Provides a generic, memory-efficient Set data structure with Union,          ║
║                     Difference, and Intersection operators for in-memory data processing.        ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
*/

// Set is a generic set implementation using an internal empty-struct map wrapper.
// An empty struct consumes zero bytes of memory, making this optimal for performance.
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates and populates a new Set with optional initial values.
func NewSet[T comparable](values ...T) *Set[T] {
	s := &Set[T]{
		items: make(map[T]struct{}),
	}
	for _, v := range values {
		s.Add(v)
	}
	return s
}

// Add inserts an element into the set.
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// Remove deletes an element from the set.
func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

// Contains checks if an element exists in the set.
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// Len returns the total number of unique elements in the set.
func (s *Set[T]) Len() int {
	return len(s.items)
}

// Union returns a new set combining all unique elements from both sets (A ∪ B).
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for item := range s.items {
		result.Add(item)
	}
	if other != nil {
		for item := range other.items {
			result.Add(item)
		}
	}
	return result
}

// Difference returns a new set containing elements present in set A that are NOT in set B (A \ B).
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for item := range s.items {
		if other == nil || !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Intersection returns a new set containing only elements common to both sets (A ∩ B).
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	if other == nil {
		return result
	}
	for item := range s.items {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Elements returns all set items as a standard Go slice for iteration.
func (s *Set[T]) Elements() []T {
	var list []T
	for item := range s.items {
		list = append(list, item)
	}
	return list
}
