package fs

import (
	"fmt"
	//"crypto/aes"
	//"crypto/cipher"
	"crypto/rand"

	"github.com/syslab-wm/nestedaes"
)

const BlockSize = 64 * 1024
//const KeySize = nestedaes.KeySize
const KeySize = 32


// TODO: add additionaldata to nestedaes
func NestedSeal(
	key *[KeySize]byte,
	buf []byte) (header []byte, err error) {

	// TODO: don't hardcode iv size
	nonce := make([]byte, 16)
	rand.Read(nonce)

	//fmt.Printf("enc w/ k %d %x nonce %d %x\n", len(key), *key, len(nonce), nonce)
	cipherText, err := nestedaes.Encrypt(buf, key[:], nonce, nil)
	println("done")
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
	key *[KeySize]byte,
	header []byte,
	buf []byte) (err error) {

	workBuf := make([]byte, len(header) + len(buf))
	copy(workBuf, header)
	copy(workBuf[len(header):], buf)

	fmt.Printf("enc w/ k %d %x\n", len(key), *key)
	plainText, err := nestedaes.Decrypt(workBuf, key[:], nil)
	println("done")
	if err != nil {
		return
	}

	if len(buf) != len(plainText) {
		panic("Crypto buffer length mismatch, what is wrong with the lib")
	}

	copy(buf, plainText)

	return
}
