package helper

func DefaultValue[T any](new, old T) T {
	if (any(new) == any("")) || (any(new) == any(0)) {
		return old
	}
	return new
}
