package nonceserver

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var nonceStorage = make(map[string]time.Time)
var nonceStorageDuplex = sync.Mutex{}
var nonceExpiry = 5 * time.Minute

func CalHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Cleanup() {
	nonceStorageDuplex.Lock()
	defer nonceStorageDuplex.Unlock()

	for nonce, timestamp := range nonceStorage {
		if time.Since(timestamp) > nonceExpiry {
			delete(nonceStorage, nonce)
		}
	}
}

func NonceDel() {
	for {
		time.Sleep(1 * time.Minute)
		Cleanup()
	}
}

func UniqueNonce() (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	nonceStr := fmt.Sprintf("%x", nonce)
	nonceStorageDuplex.Lock()
	nonceStorage[nonceStr] = time.Now()
	nonceStorageDuplex.Unlock()

	return nonceStr, nil
}
