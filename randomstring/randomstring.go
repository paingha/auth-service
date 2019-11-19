package randomstring
//GenerateRandomString generates random string of set length
import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
 "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
 rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
 b := make([]byte, length)
 for i := range b {
   b[i] = charset[seededRand.Intn(len(charset))]
 }
 return string(b)
}

func GenerateRandomString(n int) string {
	return StringWithCharset(n, charset)
}
