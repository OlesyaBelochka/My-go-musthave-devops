package prhash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(s string, key string) string {
	if key != "" {
		vSHA256 := sha256.Sum256([]byte(s + key))

		return hex.EncodeToString(vSHA256[:])
	} else {
		return ""
	}
}
