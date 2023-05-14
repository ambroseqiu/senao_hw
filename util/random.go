package util

import (
	"math/rand"
	"strings"
	"time"
)

var (
	alphabet      = "abcdefghijklmnopqrstuvwxyz"
	number        = "0123456789"
	upperAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

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

func RandomPassword(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	c := alphabet[rand.Intn(k)]
	sb.WriteByte(c)
	c = number[rand.Intn(10)]
	sb.WriteByte(c)
	c = upperAlphabet[rand.Intn(k)]
	sb.WriteByte(c)
	for i := 3; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
