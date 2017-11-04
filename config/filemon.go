package config

import (
	"crypto/sha1"
	"encoding/base64"
)

// Get a SHA1 hash out of a string
func GetFileHash(content string) string {
	hasher := sha1.New()
	hasher.Write([]byte(content))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
