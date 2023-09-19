package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"

	"github.com/googlecloudplatform/gcsfuse/internal/logger"
)

const bufSize = 64

func newAESCTRStream(key, iv []byte) (streamer cipher.Stream, err error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	streamer = cipher.NewCTR(blockCipher, iv)
	return
}

func AESCTRStream(key []byte, iv []byte, off int64, dst, src []byte) {
	var tot int64

	streamer, err := newAESCTRStream(key, iv)
	if err != nil {
		logger.Fatal("can't create AES-CTR stream: %v", err)
	}

	// produce the stream up until the point of `off`
	ciphertext := make([]byte, bufSize)
	zeros := make([]byte, bufSize)

	for tot < off {
		left := off - tot
		n := min(left, bufSize)
		streamer.XORKeyStream(ciphertext, zeros[:n])
		tot += n
	}

	// encrypt/decrypt the requested buffer
	streamer.XORKeyStream(dst, src)
}

func ReadKeyFile(path string) (key []byte, err error) {
	key, err = os.ReadFile(path)
	if err != nil {
		return
	}

	keySize := len(key)
	if keySize != 16 && keySize != 24 && keySize != 32 {
		err = fmt.Errorf("AES CTR key must be 16, 24, or 32-bytes long; got %d", keySize)
		return
	}

	return
}

func ReadIVFile(path string) (iv []byte, err error) {
	iv, err = os.ReadFile(path)
	if err != nil {
		return
	}

	if len(iv) != aes.BlockSize {
		err = fmt.Errorf("IV must be %d bytes; got %d", aes.BlockSize, len(iv))
	}

	return
}
