package helpers

// Pointer returns a pointer to the given value of any type T.
// This is a generic helper function useful for creating pointers to literals
// or values where taking the address directly is not possible.
//
// Example:
//
//	s := Pointer("hello")  // *string
//	i := Pointer(42)       // *int
//	b := Pointer(true)     // *bool
func Pointer[T any](v T) *T {
	return &v
}

// DerefPointer safely dereferences a pointer and returns its value.
// If the pointer is nil, it returns the provided defaultValue instead.
// This function is generic and works with any type T.
//
// Parameters:
//   - p: A pointer to a value of type T that may be nil
//   - defaultValue: The value to return if p is nil
//
// Returns:
//   - The dereferenced value if p is not nil, otherwise defaultValue
//
// Example:
//
//	var ptr *int
//	value := DerefPointer(ptr, 42) // returns 42
//
//	num := 10
//	value = DerefPointer(&num, 42) // returns 10
func DerefPointer[T any](p *T, defaultValue T) T {
	if p != nil {
		return *p
	}
	return defaultValue
}
