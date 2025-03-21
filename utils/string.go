package utils

import "math/rand"

func GenerateRandomString(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result += string(rune(rand.Intn(26) + 65))
	}
	return result
}
