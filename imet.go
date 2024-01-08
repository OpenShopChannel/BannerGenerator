package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"unicode/utf16"
)

var IMETMagic = [4]byte{'I', 'M', 'E', 'T'}

type IMET struct {
	_         [64]byte
	Magic     [4]byte
	HashSize  uint32
	Version   uint32
	FileSizes FileSizes
	// WiiBrew says this is a flag. We can set it to 0, and it will have no affect as far as I know
	_     uint32
	Title Languages
	_     [588]byte
	MD5   [16]byte
}

type FileSizes struct {
	IconSize   uint32
	BannerSize uint32
	SoundSize  uint32
}

type Languages struct {
	Japanese           [42]uint16
	English            [42]uint16
	German             [42]uint16
	French             [42]uint16
	Spanish            [42]uint16
	Italian            [42]uint16
	Dutch              [42]uint16
	SimplifiedChinese  [42]uint16
	TraditionalChinese [42]uint16
	Korean             [42]uint16
}

func (b *Banner) makeIMET(name string) {
	b.IMET = IMET{
		Magic:    IMETMagic,
		HashSize: 0x00000600,
		Version:  3,
		Title:    Languages{},
		MD5:      [16]byte{},
	}

	makeLanguages(&b.IMET, name)
}

func (b *Banner) CalculateIMETMD5() error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, &b.IMET)
	if err != nil {
		return err
	}

	b.IMET.MD5 = md5.Sum(buf.Bytes())
	return nil
}

func makeLanguages(imet *IMET, name string) {
	copy(imet.Title.Japanese[:], utf16.Encode([]rune(name)))
	copy(imet.Title.English[:], utf16.Encode([]rune(name)))
	copy(imet.Title.German[:], utf16.Encode([]rune(name)))
	copy(imet.Title.French[:], utf16.Encode([]rune(name)))
	copy(imet.Title.Spanish[:], utf16.Encode([]rune(name)))
	copy(imet.Title.Italian[:], utf16.Encode([]rune(name)))
	copy(imet.Title.Dutch[:], utf16.Encode([]rune(name)))
	copy(imet.Title.SimplifiedChinese[:], utf16.Encode([]rune(name)))
	copy(imet.Title.TraditionalChinese[:], utf16.Encode([]rune(name)))
	copy(imet.Title.Korean[:], utf16.Encode([]rune(name)))
}
