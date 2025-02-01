package util

import (
	"fmt"
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

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int {
	return (RandomInt(1, 1000))
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	// Simple set of characters for username generation
	const letters = "abcdefghijklmnopqrstuvwxyz"
	const numbers = "0123456789"

	// Basic domains
	domains := []string{"gmail.com", "yahoo.com", "hotmail.com"}

	// Generate random username length between 5 and 10
	length := rand.Intn(6) + 5

	// Build random username
	username := make([]byte, length)
	for i := 0; i < length; i++ {
		if rand.Float32() < 0.7 { // 70% chance of letter
			username[i] = letters[rand.Intn(len(letters))]
		} else { // 30% chance of number
			username[i] = numbers[rand.Intn(len(numbers))]
		}
	}

	// Pick random domain
	domain := domains[rand.Intn(len(domains))]

	return fmt.Sprintf("%s@%s", string(username), domain)
}
