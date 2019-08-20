package protocol

// TODO 自己实现 strconv.ParseInt
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

const (
	SimpleStringsType uint8 = iota
	ErrorsStringsType
	IntegersType
	BulkStringsType
	ArraysType
	InvalidType
)

type Message interface {
	GetRawBytes() ([]byte, error)
	ReadOne(br *bufio.Reader) error
	WriteOne(bw *bufio.Writer) error
	GetProtocolType() uint8
	GetBytesValue() []byte
	GetIntegersValue() int64
	GetArraysValue() []Message

	AppendIntegersValue(Message) Message
	AppendArraysValue(Message) Message
}

const (
	simpleStringsPrefix byte = byte('+')
	errorsStringsPrefix byte = byte('-')
	integersPrefix      byte = byte(':')
	bulkStringsPrefix   byte = byte('$')
	arraysPrefix        byte = byte('*')
	crSuffix            byte = byte('\r')
	lfSuffix            byte = byte('\n')
)

var (
	crlfSuffix []byte = []byte{crSuffix, lfSuffix}
)

var (
	parseCRLFError error = fmt.Errorf("parse crlf error")
)

// protocolType 协议类型
// SimpleStringsType|ErrorsStringsType, bytesValue is used
// IntegersType, integersValue is used
// BulkStringsType, bytesValue and integersValue is used
// ArraysType, integersValue and arraysValue is used
type message struct {
	protocolType   uint8
	bytesValue     []byte
	integersValue  int64
	arraysValue    []Message
	rawBytesBuffer *bytes.Buffer
}

func NewMessage() Message {
	return &message{
		rawBytesBuffer: new(bytes.Buffer),
	}
}

func (msg *message) readSimpleStrings(br *bufio.Reader) (err error) {
	var bytesValue []byte
	if bytesValue, err = br.ReadBytes(lfSuffix); err != nil {
		return
	}
	if !bytes.HasSuffix(bytesValue, crlfSuffix) {
		return parseCRLFError
	}
	msg.protocolType = SimpleStringsType
	msg.bytesValue = bytesValue[1 : len(bytesValue)-2]
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bytesValue); err != nil {
		return
	}
	return
}

func (msg *message) readErrorsStrings(br *bufio.Reader) (err error) {
	var bytesValue []byte
	if bytesValue, err = br.ReadBytes(lfSuffix); err != nil {
		return
	}
	if !bytes.HasSuffix(bytesValue, crlfSuffix) {
		return parseCRLFError
	}
	msg.protocolType = ErrorsStringsType
	msg.bytesValue = bytesValue[1 : len(bytesValue)-2]
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bytesValue); err != nil {
		return
	}
	return
}

func (msg *message) readIntegers(br *bufio.Reader) (err error) {
	var (
		bytesValue    []byte
		integersValue int64
	)
	if bytesValue, err = br.ReadBytes(lfSuffix); err != nil {
		return
	}
	if !bytes.HasSuffix(bytesValue, crlfSuffix) {
		return parseCRLFError
	}
	if integersValue, err = strconv.ParseInt(string(bytesValue[1:len(bytesValue)-2]), 10, 64); err != nil {
		return
	}
	msg.protocolType = IntegersType
	msg.integersValue = integersValue
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bytesValue); err != nil {
		return
	}
	return
}

func (msg *message) readBulkStrings(br *bufio.Reader) (err error) {
	var (
		bytesValue    []byte
		integersValue int64
	)
	if bytesValue, err = br.ReadBytes(lfSuffix); err != nil {
		return
	}
	if !bytes.HasSuffix(bytesValue, crlfSuffix) {
		return parseCRLFError
	}
	if integersValue, err = strconv.ParseInt(string(bytesValue[1:len(bytesValue)-2]), 10, 64); err != nil {
		return
	}
	msg.protocolType = BulkStringsType
	msg.integersValue = integersValue
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bytesValue); err != nil {
		return
	}
	if integersValue < 0 {
		return
	}
	var (
		bufferBytesValue []byte
		readNBytes       int
	)
	bytesValue = make([]byte, integersValue+2)
	for bufferBytesValue = bytesValue; len(bufferBytesValue) > 0; bufferBytesValue = bufferBytesValue[readNBytes:] {
		if readNBytes, err = br.Read(bufferBytesValue); err != nil {
			return
		}
	}
	if !bytes.HasSuffix(bytesValue, crlfSuffix) {
		return parseCRLFError
	}
	msg.bytesValue = bytesValue[:integersValue]
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bytesValue); err != nil {
		return
	}
	return
}

func (msg *message) readArrays(br *bufio.Reader) (err error) {
	var (
		bytesValue    []byte
		integersValue int64
	)
	if bytesValue, err = br.ReadBytes(lfSuffix); err != nil {
		return
	}
	if !bytes.HasSuffix(bytesValue, crlfSuffix) {
		return parseCRLFError
	}
	if integersValue, err = strconv.ParseInt(string(bytesValue[1:len(bytesValue)-2]), 10, 64); err != nil {
		return
	}
	msg.protocolType = ArraysType
	msg.integersValue = integersValue
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bytesValue); err != nil {
		return
	}
	if integersValue < 0 {
		return
	}
	var (
		elementKey   int64
		elementValue Message
	)
	msg.arraysValue = make([]Message, integersValue)
	for elementKey = 0; elementKey < integersValue; elementKey++ {
		elementValue = NewMessage()
		if err = elementValue.ReadOne(br); err != nil {
			return
		}
		if bytesValue, err = elementValue.GetRawBytes(); err != nil {
			return
		}
		if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bytesValue); err != nil {
			return
		}
		msg.arraysValue[elementKey] = elementValue
	}
	return
}

func (msg *message) ReadOne(br *bufio.Reader) (err error) {
	var prefix []byte
	if prefix, err = br.Peek(1); err != nil {
		return
	}
	switch prefix[0] {
	case simpleStringsPrefix:
		return msg.readSimpleStrings(br)
	case errorsStringsPrefix:
		return msg.readErrorsStrings(br)
	case integersPrefix:
		return msg.readIntegers(br)
	case bulkStringsPrefix:
		return msg.readBulkStrings(br)
	case arraysPrefix:
		return msg.readArrays(br)
	default:
		return fmt.Errorf("illegal first byte: %s", prefix)
	}
}

func (msg *message) writeSimpleStringsRawBytes() (err error) {
	msg.rawBytesBuffer.Reset()
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, simpleStringsPrefix); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, msg.bytesValue); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, crlfSuffix); err != nil {
		return
	}
	return
}

func (msg *message) writeErrorsStringsRawBytes() (err error) {
	msg.rawBytesBuffer.Reset()
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, errorsStringsPrefix); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, msg.bytesValue); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, crlfSuffix); err != nil {
		return
	}
	return
}

func (msg *message) writeIntegersRawBytes() (err error) {
	msg.rawBytesBuffer.Reset()
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, integersPrefix); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, []byte(strconv.FormatInt(msg.integersValue, 10))); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, crlfSuffix); err != nil {
		return
	}
	return
}

func (msg *message) writeBulkStringsRawBytes() (err error) {
	msg.rawBytesBuffer.Reset()
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bulkStringsPrefix); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, []byte(strconv.FormatInt(msg.integersValue, 10))); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, crlfSuffix); err != nil {
		return
	}
	if msg.integersValue < 0 {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, msg.bytesValue); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, crlfSuffix); err != nil {
		return
	}
	return
}

func (msg *message) writeArraysRawBytes() (err error) {
	msg.rawBytesBuffer.Reset()
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, arraysPrefix); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, []byte(strconv.FormatInt(msg.integersValue, 10))); err != nil {
		return
	}
	if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, crlfSuffix); err != nil {
		return
	}
	var (
		elementValue Message
		bytesValue   []byte
	)
	for _, elementValue = range msg.arraysValue {
		if bytesValue, err = elementValue.GetRawBytes(); err != nil {
			return
		}
		if err = binary.Write(msg.rawBytesBuffer, binary.BigEndian, bytesValue); err != nil {
			return
		}
	}
	return
}

func (msg *message) GetRawBytes() (rawBytes []byte, err error) {
	if msg.rawBytesBuffer.Len() > 0 {
		return msg.rawBytesBuffer.Bytes(), nil
	}
	switch msg.protocolType {
	case SimpleStringsType:
		err = msg.writeSimpleStringsRawBytes()
	case ErrorsStringsType:
		err = msg.writeErrorsStringsRawBytes()
	case IntegersType:
		err = msg.writeIntegersRawBytes()
	case BulkStringsType:
		err = msg.writeBulkStringsRawBytes()
	case ArraysType:
		err = msg.writeArraysRawBytes()
	default:
		err = fmt.Errorf("illegal first byte: %s", msg.protocolType)
	}
	if err != nil {
		return
	}
	return msg.rawBytesBuffer.Bytes(), nil
}

func (msg *message) WriteOne(bw *bufio.Writer) (err error) {
	var rawBytes []byte
	if rawBytes, err = msg.GetRawBytes(); err != nil {
		return
	}
	if _, err = bw.Write(rawBytes); err != nil {
		return
	}
	return bw.Flush()
}

func (msg *message) GetProtocolType() uint8 {
	return msg.protocolType
}

func (msg *message) GetBytesValue() []byte {
	return msg.bytesValue
}

func (msg *message) GetIntegersValue() int64 {
	return msg.integersValue
}

func (msg *message) GetArraysValue() []Message {
	return msg.arraysValue
}

func (msg *message) AppendIntegersValue(m Message) Message {
	msg.protocolType = IntegersType
	msg.integersValue += m.(*message).integersValue
	return msg
}

func (msg *message) AppendArraysValue(m Message) Message {
	msg.protocolType = ArraysType
	msg.integersValue++
	msg.arraysValue = append(msg.arraysValue, m.(*message))
	return msg
}
