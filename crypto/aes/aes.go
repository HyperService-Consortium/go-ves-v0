package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

// GCMCipher provides encrypto and decrypto service
type GCMCipher struct {
	aesgcm cipher.AEAD
	nonce  []byte
}

// NewGCMCipher create a middleware
func NewGCMCipher(key, nonce []byte, salt io.Reader) (*GCMCipher, error) {
	var c = make([]byte, 0)
	_, err := io.ReadFull(salt, c)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(append(c, key...))
	if err != nil {
		return nil, err
	}

	m := new(GCMCipher)
	m.aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	m.nonce = nonce
	return m, nil
}

// Seal text
func (gcp *GCMCipher) Seal(plainText []byte) (encryptedText []byte) {
	return gcp.aesgcm.Seal(nil, gcp.nonce, plainText, nil)
}

// Open encryptedText text
func (gcp *GCMCipher) Open(encryptedText []byte) ([]byte, error) {
	return gcp.aesgcm.Open(nil, gcp.nonce, encryptedText, nil)
}
