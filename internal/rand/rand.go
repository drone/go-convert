package rand

import (
	"math/rand"
	"strings"
	"time"
)

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const alphaNumeric = alpha + "0123456789"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// Alpha simple & performant random string generator
func Alpha(n int) string {
	return stringFromCharset(n, alpha)
}

// Alphanumeric simple & performant random alphanumeric string generator
func Alphanumeric(n int) string {
	return stringFromCharset(n, alphaNumeric)
}

// stringFromCharset implementation taken directly from SO, see https://stackoverflow.com/a/31832326
func stringFromCharset(n int, charset string) string {
	if n <= 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(charset) {
			sb.WriteByte(charset[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
