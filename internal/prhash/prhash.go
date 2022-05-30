package prhash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/OlesyaBelochka/My-go-musthave-devops/internal/variables"
)

func Hash(s string, key string) string {
	if key != "" {
		h := hmac.New(sha256.New, []byte(key))
		h.Write([]byte(s))

		variables.FShowLog(fmt.Sprintf("увидели ключ %s, сформировали хэш %s", key, hex.EncodeToString(h.Sum(nil))))

		return hex.EncodeToString(h.Sum(nil))
	} else {
		variables.FShowLog("Получили пустой ключ, хеш не проверяем")
		return ""
	}
}
