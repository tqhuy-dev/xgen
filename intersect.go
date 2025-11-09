package main

// Contains returns true if an element is present in a collection.
func Contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}
	return false
}

// ContainsBy returns true if predicate function returns true for any element.
func ContainsBy[T any](collection []T, predicate func(item T) bool) bool {
	for _, item := range collection {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Every returns true if all elements of a subset are contained in a collection or if the subset is empty.
func Every[T comparable](collection []T, subset []T) bool {
	for _, item := range subset {
		if !Contains(collection, item) {
			return false
		}
	}
	return true
}

// EveryBy returns true if the predicate returns true for all elements in the collection or if the collection is empty.
func EveryBy[T any](collection []T, predicate func(item T) bool) bool {
	for _, item := range collection {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Some returns true if at least 1 element of a subset is contained in a collection.
// If the subset is empty Some returns false.
func Some[T comparable](collection []T, subset []T) bool {
	for _, item := range subset {
		if Contains(collection, item) {
			return true
		}
	}
	return false
}

// SomeBy returns true if the predicate returns true for any of the elements in the collection.
// If the collection is empty SomeBy returns false.
func SomeBy[T any](collection []T, predicate func(item T) bool) bool {
	for _, item := range collection {
		if predicate(item) {
			return true
		}
	}
	return false
}

// None returns true if no element of a subset is contained in a collection or if the subset is empty.
func None[T comparable](collection []T, subset []T) bool {
	for _, item := range subset {
		if Contains(collection, item) {
			return false
		}
	}
	return true
}

// NoneBy returns true if the predicate returns true for none of the elements in the collection or if the collection is empty.
func NoneBy[T any](collection []T, predicate func(item T) bool) bool {
	for _, item := range collection {
		if predicate(item) {
			return false
		}
	}
	return true
}

// Intersect returns the intersection between two collections.
// Returns elements that appear in both list1 and list2.
func Intersect[T comparable, Slice ~[]T](list1 Slice, list2 Slice) Slice {
	result := Slice{}
	seen := make(map[T]struct{}, len(list1))

	for _, item := range list1 {
		seen[item] = struct{}{}
	}

	for _, item := range list2 {
		if _, ok := seen[item]; ok {
			result = append(result, item)
			delete(seen, item) // Avoid duplicates in result
		}
	}

	return result
}

// Difference returns the difference between two collections.
// The first value is the collection of elements absent from list2.
// The second value is the collection of elements absent from list1.
func Difference[T comparable, Slice ~[]T](list1 Slice, list2 Slice) (Slice, Slice) {
	left := Slice{}
	right := Slice{}

	seenLeft := make(map[T]struct{}, len(list1))
	seenRight := make(map[T]struct{}, len(list2))

	for _, item := range list1 {
		seenLeft[item] = struct{}{}
	}

	for _, item := range list2 {
		seenRight[item] = struct{}{}
	}

	for _, item := range list1 {
		if _, ok := seenRight[item]; !ok {
			left = append(left, item)
		}
	}

	for _, item := range list2 {
		if _, ok := seenLeft[item]; !ok {
			right = append(right, item)
		}
	}

	return left, right
}

// Union returns all distinct elements from given collections.
// Result preserves the relative order of elements.
func Union[T comparable, Slice ~[]T](lists ...Slice) Slice {
	var capLen int
	for _, list := range lists {
		capLen += len(list)
	}

	result := make(Slice, 0, capLen)
	seen := make(map[T]struct{}, capLen)

	for _, list := range lists {
		for _, item := range list {
			if _, ok := seen[item]; !ok {
				seen[item] = struct{}{}
				result = append(result, item)
			}
		}
	}

	return result
}

// Without returns slice excluding all given values.
func Without[T comparable, Slice ~[]T](collection Slice, exclude ...T) Slice {
	excludeMap := make(map[T]struct{}, len(exclude))
	for _, item := range exclude {
		excludeMap[item] = struct{}{}
	}

	result := make(Slice, 0, len(collection))
	for _, item := range collection {
		if _, ok := excludeMap[item]; !ok {
			result = append(result, item)
		}
	}
	return result
}

// WithoutEmpty returns slice excluding empty values.
//
// Deprecated: Use Compact instead.
func WithoutEmpty[T comparable, Slice ~[]T](collection Slice) Slice {
	return Compact(collection)
}
