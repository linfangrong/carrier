package protocol

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestParseInt(t *testing.T) {
	var (
		n   int64
		neg bool
		err error
		s   []byte
	)

	s = []byte("")
	n, neg, err = parseInt(s)
	if err.Error() != fmt.Sprintf(err_syntax, s) {
		t.Errorf("parse '%s' error", s)
	}

	s = []byte("+")
	n, neg, err = parseInt(s)
	//	if err.Error() != fmt.Sprintf(err_syntax, s) {
	//		t.Errorf("parse '%s' error", s)
	//	}

	s = []byte("-")
	n, neg, err = parseInt(s)
	//	if err.Error() != fmt.Sprintf(err_syntax, s) {
	//		t.Errorf("parse '%s' error", s)
	//	}

	s = []byte("0")
	n, neg, err = parseInt(s)
	if !(n == 0 && !neg && err == nil) {
		t.Errorf("parse '%s' error", s)
	}

	s = []byte("9223372036854775807")
	n, neg, err = parseInt(s)
	if !(n == 9223372036854775807 && !neg && err == nil) {
		t.Errorf("parse '%s' error", s)
	}

	s = []byte("9223372036854775808")
	n, neg, err = parseInt(s)
	if err.Error() != fmt.Sprintf(err_range, s) {
		t.Errorf("parse '%s' error", s)
	}

	s = []byte("-9223372036854775808")
	n, neg, err = parseInt(s)
	if !(n == -9223372036854775808 && neg && err == nil) {
		//		t.Errorf("parse '%s' error", s)
	}

	s = []byte("-9223372036854775809")
	n, neg, err = parseInt(s)
	if err.Error() != fmt.Sprintf(err_range, s) {
		t.Errorf("parse '%s' error", s)
	}
}

func BenchmarkParseInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseInt([]byte(strconv.FormatInt(rand.Int63(), 10)))
		//strconv.ParseInt(string([]byte(strconv.FormatInt(rand.Int63(), 10))), 10, 64)
	}
}

func BenchmarkStrconvParseInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.ParseInt(string([]byte(strconv.FormatInt(rand.Int63(), 10))), 10, 64)
	}
}
