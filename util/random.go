package util

import (
	"math/rand"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	//

}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	/* for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	} */

	for range n {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwnerOld generates a random owner name
func RandomOwnerOld() string {
	return RandomString(6)
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return gofakeit.Name()
}

// RandomMomey generates a random amount of money
func RandomMomey() int64 {
	return RandomInt(0, 1000)
}

// RandomOwnerName generates a random currncy code
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
