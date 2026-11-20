package helpers

func GetOrDefault[T any](value *T, d T) T {
	if value == nil {
		return d
	}

	return *value
}
