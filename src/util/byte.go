package util

func ToUpper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c & 0x5F
	}
	return c
}

func ToLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c | 0x20
	}
	return c
}
