package service

import (
	"crypto/sha256"
	"fmt"
)

// Checksum computes a SHA-256 digest of data and returns it as "sha256:<hex>".
func Checksum(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", sum)
}
