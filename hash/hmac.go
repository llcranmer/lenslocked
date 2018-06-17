package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// HMAC is a wrapper around the hash.Hash type
type HMAC struct {
	hmac hash.Hash
}

// NewHMAC creates new hash and returns new HMAC type with set hash value
func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{h}
}

// Hash will hash provided input string with key that was provided when HMAC
// struct was created
// and will return url save hash
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
