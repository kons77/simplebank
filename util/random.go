package util

import (
	crand "crypto/rand"
	"math/rand"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// var rng = rand.New(rand.NewSource(time.Now().UnixNano()))  -  deprecated

func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		panic("cannot generate random bytes: " + err.Error())
	}
	return b
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
func RandomUsername() string {
	// return RandomString(6)
	return gofakeit.Username()

}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return gofakeit.Name()
}

// RandomEmail generates a random owner name
func RandomEmail() string {
	// return fmt.Sprintf("%s@email.com", RandomString(6))
	return gofakeit.Email()
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomOwnerName generates a random currncy code
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
