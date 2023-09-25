package tool

func ContainStr(s []string, sub string) bool {
	for _, v := range s {
		if v == sub {
			return true
		}
	}
	return false
}
