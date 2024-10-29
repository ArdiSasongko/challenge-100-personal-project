package utils

func DefaultValue[T any](old, new any) T {
	switch v := any(new).(type) {
	case nil:
		return old.(T)
	case int:
		if v == 0 {
			return old.(T)
		}
	case string:
		if v == "" {
			return old.(T)
		}
	}
	return new.(T)
}
