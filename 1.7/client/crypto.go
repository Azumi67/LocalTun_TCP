package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("invalid private key data")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing private key failed: %v", err)
		}
		privateKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not RSA private key")
		}
		return privateKey, nil
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}
}
