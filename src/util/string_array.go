package util

func Unique(slice []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueSlice []string
	for _, entry := range slice {
		if _, value := uniqueMap[entry]; !value {
			uniqueMap[entry] = true
			uniqueSlice = append(uniqueSlice, entry)
		}
	}
	return uniqueSlice
}

func Filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func Contains(a []string, s string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}
	return false
}
