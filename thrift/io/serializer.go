package io

import (
	"bytes"
	"encoding/binary"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	// Signature pinpoint singnature
	Signature int8 = -17
	// DtoVersion dto version
	DtoVersion uint8 = 16
)

// EncodeTstruct ...
func EncodeTstruct(ttype uint16, tstruct thrift.TStruct) ([]byte, error) {
	// buffer := bytes.NewBuffer(make([]byte, 4))

	var buffer bytes.Buffer
	var err error
	// signature >b 1
	err = binary.Write(&buffer, binary.BigEndian, Signature)
	if err != nil {
		return nil, err
	}
	// version >B 1
	err = binary.Write(&buffer, binary.BigEndian, DtoVersion)
	if err != nil {
		return nil, err
	}
	// type >H 2
	err = binary.Write(&buffer, binary.BigEndian, ttype)
	if err != nil {
		return nil, err
	}
	// thrift
	tmem := thrift.TMemoryBuffer{Buffer: &buffer}
	err = tstruct.Write(thrift.NewTCompactProtocol(&tmem))
	if err != nil {
		return nil, err
	}

	/*
		_, err = buffer.Write(tmem.Bytes())
		if err != nil {
			return nil, err
		}
	*/

	return buffer.Bytes(), nil
}

// EncodeTCPTTstuctRequest ...
func EncodeTCPTTstuctRequest(messageID uint32, ttype uint16, tstruct thrift.TStruct) ([]byte, error) {
	tstructData, err := EncodeTstruct(ttype, tstruct)
	if err != nil {
		return nil, err
	}

	// buffer := bytes.NewBuffer(make([]byte, 10))
	var buffer bytes.Buffer
	// request type >H 2
	err = binary.Write(&buffer, binary.BigEndian, RequestTypeAppRequest)
	if err != nil {
		return nil, err
	}
	// message id   >I 4
	err = binary.Write(&buffer, binary.BigEndian, messageID)
	if err != nil {
		return nil, err
	}
	// data length  >I 4
	var length uint32 = uint32(len(tstructData))
	err = binary.Write(&buffer, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	// tstruct data
	_, err = buffer.Write(tstructData)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
