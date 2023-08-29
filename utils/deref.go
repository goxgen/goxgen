package utils

func Deref[T any](s *T) T {
	if s != nil {
		return *s
	}

	return *new(T)
}
