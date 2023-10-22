package fs

import (
	//"fmt"
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
)

const BlockSize = 64 * 1024

// Increments a big-endian byte array by an integer. Good for ctr cipher things.
// Go might have a builtin way to do this but idk what it is
func incCtr(ctr []byte, inc int64) {
	// Stole this loop from some C code I wrote a few months ago
	for i := len(ctr) - 1; i >= 0; i-- {
		a := ctr[i]
		ctr[i] += uint8(inc)
		// Did we overflow?
		carry := (ctr[i] < uint8(inc) + a)
		inc >>= 8
		if carry {
			inc++
		}
		if inc == 0 {
			break
		}
	}
}

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

	cipherText := gcmCipher.Seal(buf[:0], nonce, buf, additionalData)
	tag = cipherText[len(buf):]

	// If the underlying array didn't have enough space for the tag, cipher lib might have quietly reallocated
	if &buf[0] != &cipherText[0] {
		copy(buf, cipherText)
	}

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

	workbuf := make([]byte, len(buf) + len(tag))
	copy(workbuf, buf)
	copy(workbuf[len(buf):], tag)

	plainText, err := gcmCipher.Open(workbuf[:0], nonce, workbuf, additionalData)

	copy(buf, plainText)

	return
}

func BadCtrCrypt(
	key []byte,
	offset int64,
	buf []byte) (err error) {
	var blockCipher cipher.Block
	blockCipher, err = aes.NewCipher(key);
	if err != nil {
		return
	}

	// Using constant IV of 0. This is bad bad bad, very insecure, do not do this. demonstration purposes only.
	var iv [16]byte
	// First, advance the IV (counter) appropriately for the file offset
	//fmt.Println(iv)
	incCtr(iv[:], offset / 16)
	//fmt.Println(offset)
	//fmt.Println(iv)
	ctrCipher := cipher.NewCTR(blockCipher, iv[:])
	// If the offset isn't aligned to the blocksize, crypt some junk bytes to advance the XOR stream appropriately
	if offset % 16 != 0 {
		var junk [16]byte
		junkSlice := junk[:offset % 16]
		ctrCipher.XORKeyStream(junkSlice, junkSlice)
	}

	// Crypt the data
	ctrCipher.XORKeyStream(buf, buf)
	return
}
