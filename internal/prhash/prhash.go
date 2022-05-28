package prhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Hash(s string, key string) string {
	if key != "" {

		vSHA256 := sha256.Sum256([]byte(s + key))
		fmt.Println("Агент увидел ключ ", key, "сформировал хэш ", hex.EncodeToString(vSHA256[:]))

		return hex.EncodeToString(vSHA256[:])
	} else {

		return ""
	}
}
