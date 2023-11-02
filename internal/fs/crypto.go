package fs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

const BlockSize = 64 * 1024

func GcmSeal(
	key []byte,
	buf []byte,
	additionalData []byte) (nonce []byte, tag []byte, err error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcmCipher, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return
	}

	nonce = make([]byte, gcmCipher.NonceSize())
	rand.Read(nonce)

	workBuf := make([]byte, len(buf) + gcmCipher.Overhead())
	cipherText := gcmCipher.Seal(workBuf[:0], nonce, buf, additionalData)
	tag = cipherText[len(buf):]

	copy(buf, cipherText)

	return
}

func GcmOpen(
	key []byte,
	nonce []byte,
	buf []byte,
	additionalData []byte,
	tag []byte) (err error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcmCipher, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return
	}

	workBuf := make([]byte, len(buf) + len(tag))
	copy(workBuf, buf)
	copy(workBuf[len(buf):], tag)

	plainText, err := gcmCipher.Open(buf[:0], nonce, workBuf, additionalData)

	if len(buf) != len(plainText) {
		panic("Crypto buffer length mismatch, this should not happen")
	}

	if &buf[0] != &plainText[0] {
		panic("Crypto reallocated when it shouldn't have, this should not happen")
	}

	return
}
