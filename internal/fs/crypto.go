package fs

import (
	//"crypto/aes"
	//"crypto/cipher"
	"crypto/rand"

	"github.com/syslab-wm/nestedaes"
	//"github.com/syslab-wm/nestedaes/internal/aesx"
)

const BlockSize = 64 * 1024

// TODO: add additionaldata to nestedaes
func NestedSeal(
	key []byte,
	buf []byte) (header []byte, err error) {

	// TODO: don't hardcode iv size
	nonce := make([]byte, 16)
	rand.Read(nonce)

	cipherText, err := nestedaes.Encrypt(buf, key, nonce)
	if err != nil {
		return
	}

	header, payload, err := nestedaes.SplitHeaderPayload(cipherText)
	if err != nil {
		return
	}

	if len(payload) != len(buf) {
		panic("Crypto buffer length mismatch, what is wrong with the lib")
	}

	copy(buf, payload)

	return
}

func NestedOpen(
	key []byte,
	header []byte,
	buf []byte) (err error) {

	workBuf := make([]byte, len(header) + len(buf))
	copy(workBuf, header)
	copy(workBuf[len(header):], buf)

	plainText, err := nestedaes.Decrypt(workBuf, key)
	if err != nil {
		return
	}

	if len(buf) != len(plainText) {
		panic("Crypto buffer length mismatch, what is wrong with the lib")
	}

	copy(buf, plainText)

	return
}
