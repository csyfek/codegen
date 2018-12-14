package util

func SetContainsString(set []string, s string) bool {
	for _, s_ := range set {
		if s == s_ {
			return true
		}
	}
	return false
}

