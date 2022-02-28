package handler

import (
	"math/rand"
	"time"
)

func GetSalt(saltLen int) string {
	rand.Seed(time.Now().UnixNano())
	base62 := Encode(time.Now().Nanosecond())
	chars := []rune(base62)
	saltBase62 := ""
	for index := 0; index < saltLen; index++ {
		saltBase62 += string(chars[rand.Intn(len(base62))])
	}
	return saltBase62
}
