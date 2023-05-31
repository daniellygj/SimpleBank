package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// this function will be called automatically when the package is first used
func init() {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(stringSize int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < stringSize; i++ {
		letter := alphabet[rand.Intn(k)]
		sb.WriteByte(letter)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
