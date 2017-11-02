package protocol

import (
	"fmt"
)

var (
	err_syntax string = "parsing '%s': invalid syntax"
	err_range  string = "parsing '%s': value out of range"
)

func parseInt(s []byte) (n int64, neg bool, err error) {
	var (
		base    int = 10
		bitSize int = 64
		len     int = len(s)
		i       int = 0
		v       byte

		maxInt int64 = 1<<uint(bitSize-1) - 1
		minInt int64 = -1 << uint(bitSize-1)
		// Cutoff is the smallest number such that cutoff*base is overflows.
		maxCutoff int64 = maxInt/int64(base) + 1
		minCutoff int64 = minInt/int64(base) - 1
		// Pluson is check overflows.
		pluson int64
	)
	// Empty string bad.
	if len == 0 {
		return 0, false, fmt.Errorf(err_syntax, s)
	}

	// Pick off leading sign.
	switch s[0] {
	case '+':
		neg = false
		i = 1
	case '-':
		neg = true
		i = 1
	default:
		neg = false
		i = 0
	}
	if len == i {
		return 0, neg, fmt.Errorf(err_syntax, s)
	}

	for ; i < len; i++ {
		switch {
		case '0' <= s[i] && s[i] <= '9':
			v = s[i] - '0'
		case 'a' <= s[i] && s[i] <= 'z':
			v = s[i] - 'a' + 10
		case 'A' <= s[i] && s[i] <= 'Z':
			v = s[i] - 'A' + 10
		default:
			return 0, neg, fmt.Errorf(err_syntax, s)
		}
		if v >= byte(base) {
			return 0, neg, fmt.Errorf(err_syntax, s)
		}

		if n >= maxCutoff {
			// n*base overflows
			return maxInt, neg, fmt.Errorf(err_range, s)
		}
		if n <= minCutoff {
			// n*base overflows
			return minInt, neg, fmt.Errorf(err_range, s)
		}
		n *= int64(base)

		if !neg {
			pluson = n + int64(v)
			if pluson < n || pluson > maxInt {
				// n+v overflows
				return maxInt, neg, fmt.Errorf(err_range, s)
			}
		} else {
			pluson = n - int64(v)
			if pluson > n || pluson < minInt {
				// n-v overflows
				return minInt, neg, fmt.Errorf(err_range, s)
			}
		}
		n = pluson
	}
	return
}
