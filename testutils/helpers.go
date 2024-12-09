package testutils

func PointerTo[T any](v T) *T {
	return &v
}
