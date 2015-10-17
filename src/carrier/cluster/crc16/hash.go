package crc16

import (
	"bytes"
)

const (
	leftBracket  byte = byte('{')
	rightBracket byte = byte('}')
)

func HashSlot(buf []byte) (hash uint16) {
	var (
		leftIndex  int
		rightIndex int
	)
	if leftIndex = bytes.IndexByte(buf, leftBracket); leftIndex != -1 {
		if rightIndex = bytes.IndexByte(buf[leftIndex+1:], rightBracket); rightIndex > 0 {
			buf = buf[leftIndex+1 : leftIndex+rightIndex+1]
		}
	}
	return crc16(buf) & 0x3FFF //16383
}
