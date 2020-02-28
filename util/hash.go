package util

import (
	"crypto/sha512"
	"encoding/hex"
)

func Hash(s string) string {
	h := sha512.New()
	h.Write([]byte(s))
	pass := hex.EncodeToString(h.Sum(nil))

	return pass
}
