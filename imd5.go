package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
)

var IMD5Magic = [4]byte{'I', 'M', 'D', '5'}

type IMD5 struct {
	Magic    [4]byte
	Filesize uint32
	_        [8]byte
	MD5      [16]byte
}

func makeIMD5(file []byte) ([]byte, error) {
	imd5 := IMD5{
		Magic:    IMD5Magic,
		Filesize: uint32(len(file)),
		MD5:      md5.Sum(file),
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, imd5)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}
