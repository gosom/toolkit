package cryptoext

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	"toolkit/pkg/errorsext"
)

const (
	keySize = 32
)

var (
	ErrorInvalidKeySize = errors.New("invalid key size")
	ErrInvalidCipher    = errors.New("invalid cipher")
)

func Encrypt(key string, data []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, errorsext.WithStack(ErrorInvalidKeySize)
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errorsext.WithStack(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, errorsext.WithStack(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errorsext.WithStack(err)
	}

	encrypted := gcm.Seal(nonce, nonce, data, nil)

	return encrypted, nil
}

func Decrypt(key string, data []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, errorsext.WithStack(ErrorInvalidKeySize)
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errorsext.WithStack(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, errorsext.WithStack(err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errorsext.WithStack(ErrInvalidCipher)
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func Hash(data []byte) string {
	h := sha256.New()

	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRandomString(size int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", errorsext.WithStack(err)
	}

	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}

	return string(b), nil
}
