package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomInt returns a random integer between min and max
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
