package utils

func StripZeroPadding(hexString string) string {
	result := ""
	strip := true
	for _, c := range hexString {
		if c == '0' && strip {
			continue
		}
		strip = false
		result += string(c)
	}

	if len(result) == 0 {
		result = "0"
	}
	return result
}
