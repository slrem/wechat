package trader

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func sha(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
func GetRandStr(a int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < a; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
