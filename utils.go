package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	_ "embed"
	"encoding/binary"
	"github.com/wii-tools/wadlib"
	"os"
	"strconv"
)

// tmdTemplate is a TMD that has most fields needed by homebrew
//
//go:embed templates/tmd
var tmdTemplate []byte

// contentAesKey is the AES key that is used to encrypt title contents.
var contentAesKey = [16]byte{0x72, 0x95, 0xDB, 0xC0, 0x47, 0x3C, 0x90, 0x0B, 0xB5, 0x94, 0x19, 0x9C, 0xB5, 0xBC, 0xD3, 0xDC}

// encryptContents encrypts both the banner and zip.
func (a *App) encryptContents(decrypted []byte, zipPath string) ([]byte, []byte, []byte, error) {
	// Get integer version of Title ID and pass it to the ticket to create the AES key
	titleId, err := strconv.ParseInt(a.Shop.TitleID, 16, 64)
	if err != nil {
		return nil, nil, nil, err
	}

	ticket := wadlib.Ticket{TitleID: uint64(titleId)}
	ticket.UpdateTitleKey(contentAesKey)
	titleKey := ticket.GetTitleKey()

	// Create TMD and calculate SHA1 sum
	tmd, err := a.createTmd(uint64(titleId))
	if err != nil {
		return nil, nil, nil, err
	}

	tmd.Contents[0].Size = uint64(len(decrypted))
	tmd.Contents[0].Hash = sha1.Sum(decrypted)

	// Open the zip file to encrypt
	zipFile, err := os.ReadFile(zipPath + a.Slug + ".zip")
	if err != nil {
		return nil, nil, nil, err
	}

	tmd.Contents[3].Size = uint64(len(zipFile))
	tmd.Contents[3].Hash = sha1.Sum(zipFile)

	encBanner, err := encryptContent(titleKey, decrypted, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	encZip, err := encryptContent(titleKey, zipFile, 3)
	if err != nil {
		return nil, nil, nil, err
	}

	// Write TMD to byte slice
	wad := wadlib.WAD{}
	wad.TMD = *tmd

	tmdBytes, err := wad.GetTMD()
	if err != nil {
		return nil, nil, nil, err
	}

	// Wii Shop Channel expects a TMD with a certificate chain appended to it.
	tmdBytes = append(tmdBytes, tmdTemplate[len(tmdBytes):]...)

	return encBanner, encZip, tmdBytes, nil
}

// encryptContent encrypts a singular content.
func encryptContent(key [16]byte, decrypted []byte, cid int) ([]byte, error) {
	// The IV we'll use will be the two bytes sourced from the content's index,
	// padded with 14 null bytes.
	var indexBytes [2]byte
	binary.BigEndian.PutUint16(indexBytes[:], uint16(cid))

	iv := make([]byte, 16)
	iv[0] = indexBytes[0]
	iv[1] = indexBytes[1]

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)

	// Pad the file to 16 bytes for encryption
	paddedSize := len(decrypted)
	leftover := paddedSize % 16
	if leftover != 0 {
		paddedSize += 16 - leftover
	}

	decryptedData := make([]byte, paddedSize)
	copy(decryptedData, decrypted)

	encryptedData := make([]byte, len(decryptedData))

	blockMode.CryptBlocks(encryptedData, decryptedData)
	return encryptedData, err
}

func (a *App) createTmd(titleId uint64) (*wadlib.TMD, error) {
	wad := wadlib.WAD{}
	err := wad.LoadTMD(tmdTemplate)
	if err != nil {
		return nil, err
	}

	tmd := wad.TMD
	tmd.TitleID = titleId
	tmd.TitleVersion = uint16(a.Shop.Version)

	return &tmd, nil
}
