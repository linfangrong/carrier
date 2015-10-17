package protocol_test

import (
	"bufio"
	"bytes"
	"carrier/protocol"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&MessageTest{})

type MessageTest struct {
}

func (msg *MessageTest) SetUpTest(c *C) {

}

func (msg *MessageTest) TearDownTest(c *C) {

}

// SimpleStrings
func (msg *MessageTest) TestMessageSimpleStrings(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("+OK\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

// ErrorsStrings
func (msg *MessageTest) TestMessageErrorsStrings(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("-Error message\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

// Integers
func (msg *MessageTest) TestIntegers(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString(":1000\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

// BulkStrings
func (msg *MessageTest) TestBulkStrings(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("$6\r\nfoobar\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

func (msg *MessageTest) TestEmptyStrings(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("$0\r\n\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

func (msg *MessageTest) TestNullBulkStrings(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("$-1\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

// Arrsys
func (msg *MessageTest) TestEmptyArrays(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("*0\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

func (msg *MessageTest) TestBulkStringsArrays(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

func (msg *MessageTest) TestIntegersArrays(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("*3\r\n:1\r\n:2\r\n:3\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

func (msg *MessageTest) TestIntegersBulkStringArrays(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

func (msg *MessageTest) TestNullArrays(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("*-1\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

func (msg *MessageTest) TestMultiArrays(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}

func (msg *MessageTest) TestNullElementsArrays(c *C) {
	var (
		buf  *bytes.Buffer    = bytes.NewBufferString("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n")
		br   *bufio.Reader    = bufio.NewReader(buf)
		pmsg protocol.Message = protocol.NewMessage()
		err  error
	)
	err = pmsg.ReadOne(br)
	c.Assert(err, IsNil)
	//	c.Error(pmsg)
}
