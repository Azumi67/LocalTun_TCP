package hashclient

import (
	"crypto/sha256"
)

func GenHash(challenge []byte, secretKey string) []byte {
	hashres := sha256.New()
	hashres.Write(challenge)
	hashres.Write([]byte(secretKey))
	return hashres.Sum(nil)
}
