package rand

import (
	"math/rand"
	"time"
)

func init() {
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// String return random string
func String(length int) string {
	return stringWithCharset(length, charset)
}
