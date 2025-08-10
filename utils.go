package main

func toUint8(v interface{}) uint8 {
	switch val := v.(type) {
	case int:
		if val < 0 || val > 255 {
			return 0
		}
		return uint8(val)
	case int32:
		if val < 0 || val > 255 {
			return 0
		}
		return uint8(val)
	case int64:
		if val < 0 || val > 255 {
			return 0
		}
		return uint8(val)
	case uint:
		if val > 255 {
			return 0
		}
		return uint8(val)
	case float32:
		i := int(val)
		if i < 0 || i > 255 {
			return 0
		}
		return uint8(i)
	case float64:
		i := int(val)
		if i < 0 || i > 255 {
			return 0
		}
		return uint8(i)
	default:
		return 0
	}
}
