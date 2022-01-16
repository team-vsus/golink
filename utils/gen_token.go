package utils

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateNewToken(n int) string {
	t := make([]rune, n)
	for i := range t {
		t[i] = letters[rand.Intn(len(letters))]
	}
	return string(t)
}
