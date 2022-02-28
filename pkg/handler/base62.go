package handler

import (
	"errors"
	"math"
	"strings"
)

var (
	Base         = 62
	CharacterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func Encode(num int) string {
	base62 := make([]byte, 0)
	for num > 0 {
		r := math.Mod(float64(num), float64(Base))
		num /= Base
		base62 = append([]byte{CharacterSet[int(r)]}, base62...)
	}
	return string(base62)
}

func Decode(s string) (int, error) {
	var r, pow int
	for i, v := range s {
		pow = len(s) - (i + 1)
		pos := strings.IndexRune(CharacterSet, v)
		if pos == -1 {
			return pos, errors.New("invalid character: " + string(v))
		}
		r += pos * int(math.Pow(float64(Base), float64(pow)))
	}
	return int(r), nil
}
