package alerts

// Compare evaluates value against threshold using the specified operator string.
// Supported operators:
//   - ">"
//   - ">="
//   - "<"
//   - "<="
//   - "==" or "="
//   - "!="
func Compare(value float64, operator string, threshold float64) bool {
	switch operator {
	case ">":
		return value > threshold
	case ">=":
		return value >= threshold
	case "<":
		return value < threshold
	case "<=":
		return value <= threshold
	case "==", "=":
		return value == threshold
	case "!=":
		return value != threshold
	default:
		return false
	}
}
