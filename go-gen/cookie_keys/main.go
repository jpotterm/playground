package main

import (
	"crypto/rand"
	"fmt"
)

func main() {
	cookieHashKey := make([]byte, 64)
	_, err := rand.Read(cookieHashKey)
	if err != nil {
		panic(err)
	}

	cookieBlockKey := make([]byte, 32)
	_, err = rand.Read(cookieBlockKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(cookieHashKey)
	fmt.Println(cookieBlockKey)
}
