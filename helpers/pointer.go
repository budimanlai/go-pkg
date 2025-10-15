package helpers

// Pointer returns a pointer to the given value of any type.
func Pointer[T any](v T) *T {
	return &v
}

// DerefPointer dereferences a pointer to a value of any type, returning a default value if the pointer is nil.
func DerefPointer[T any](p *T, defaultValue T) T {
	if p != nil {
		return *p
	}
	return defaultValue
}
