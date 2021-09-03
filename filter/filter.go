package filter

// StringSlice scan elements in the source, if condition true, put it into impurity slice or pure slice
// return impurity []string, pure []string
func StringSlice(source []string, condition func(val string) bool) ([]string, []string) {
	impurity := make([]string, 0)
	pure := make([]string, 0)

	if condition == nil {
		pure = source[:]

		return impurity, pure
	}

	for _, val := range source {
		if condition(val) {
			impurity = append(impurity, val)

			continue
		}

		pure = append(pure, val)
	}

	return impurity, pure
}
