package file

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

// Hash 计算哈希
func Hash(reader io.Reader, algorithm HashAlgorithm) (string, error) {

	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	if _, err := io.Copy(h, reader); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
