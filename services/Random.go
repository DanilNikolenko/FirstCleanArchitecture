package services

import (
	"math/rand"
	"time"
)

const (
	MIN = 97
	MAX = 122
)

func GetApplicationRandomNAme() string {
	var name string
	for i := 0; i < 2; i++ {
		TempRand := Random(MIN, MAX)
		name += string(byte(TempRand))
	}

	return name
}

func Random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
