package crypto

import (
	"github.com/bootstrap-library/stock-plus/env"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

// Crypto type
// AES/ECB/PKCS5Padding/128bit/base64

var key string

func init() {
	key = env.String("PASSWORD")
	if key == "" {
		panic("env PASSWORD is required")
	}
	for len(key) < 16 {
		key += "\x00"
	}
	if len(key) > 16 {
		key = key[:16]
	}
}

func Encrypt(s string) string {
	return crypto.
		FromString(s).
		SetKey(key).
		Aes().
		ECB().
		PKCS5Padding().
		Encrypt().
		ToBase64String()
}

func Decrypt(s string) string {
	return crypto.
		FromBase64String(s).
		SetKey(key).
		Aes().
		ECB().
		PKCS5Padding().
		Decrypt().
		ToString()
}
