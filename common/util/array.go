package util

func Unique(slices []string) []string {
	result := make([]string, 0)
	m := make(map[string]struct{}, 0)
	for _, s := range slices {
		if _, ok := m[s]; !ok {
			result = append(result, s)
			m[s] = struct{}{}
		}
	}
	return result
}
