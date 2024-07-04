package utils

import (
	"math/rand"
	"sync"
	"time"
)

const (
	lowercaseChars   = "abcdefghijklmnopqrstuvwxyz"
	number           = "0123456789"
	stringNumCharset = lowercaseChars + number
)

var (
	r    *rand.Rand
	once sync.Once
)

func initRand() {
	seed := rand.NewSource(time.Now().UnixNano())
	r = rand.New(seed)
}

func GenerateRandomString(length int) string {
	once.Do(initRand)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = lowercaseChars[r.Intn(len(lowercaseChars))]
	}

	return string(result)
}

func GenerateRandomStringWithNumber(length int) string {
	once.Do(initRand)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = stringNumCharset[r.Intn(len(stringNumCharset))]
	}

	return string(result)
}
