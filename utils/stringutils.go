package utils

func StringsReverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func StringsInArray(s []string, e string) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}

	return false, -1
}

func StringsExcept(ss1 []string, ss2 []string) []string {
	if len(ss1) == 0 {
		return make([]string, 0)
	}

	if len(ss2) == 0 {
		return ss1
	}

	s2map := make(map[string]bool)
	for _, s2 := range ss2 {
		s2map[s2] = true
	}

	se := make([]string, 0, len(ss1))
	for _, s1 := range ss1 {
		if _, found := s2map[s1]; !found {
			se = append(se, s1)
		}
	}
	return se
}

func StringsDistinct(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}

	return u
}
