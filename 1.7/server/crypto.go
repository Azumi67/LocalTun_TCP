package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("invalid public key data")
	}

	switch block.Type {
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing PKIX public key failed: %v", err)
		}
		publicKey, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("not RSA public key")
		}
		return publicKey, nil
	case "RSA PUBLIC KEY":
		publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing public key failed: %v", err)
		}
		return publicKey, nil
	default:
		return nil, fmt.Errorf("unsupported public key type: %s", block.Type)
	}
}urn nil, fmt.Errorf("invalid public key data")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}
