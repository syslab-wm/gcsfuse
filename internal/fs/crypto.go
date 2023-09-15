package fs

import (
	//"fmt"
	"crypto/aes"
	"crypto/cipher"
)

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
