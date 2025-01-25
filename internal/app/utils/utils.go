package utils

import (
	"math/rand"
	"time"

)

func GenerateRandomString(length int) string {
    characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    rand.Seed(time.Now().UnixNano()) 

    result := make([]byte, length)
    for i := 0; i < length; i++ {
        result[i] = characters[rand.Intn(len(characters))]
    }
    return string(result)
}


