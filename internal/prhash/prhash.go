package prhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

func Hash(s string, key string) string {
	if key != "" {

		vSHA256 := sha256.Sum256([]byte(s + key))
		variables.FShowLog(fmt.Sprintf("Агент увидел ключ %s, сформировал хэш ", key, hex.EncodeToString(vSHA256[:])))

		return hex.EncodeToString(vSHA256[:])
	} else {

		variables.FShowLog(fmt.Sprintf("Получили пустой ключ, хеш не проверяем"))

		return ""
	}
}
