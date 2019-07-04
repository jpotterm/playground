package main

import (
	"crypto/rand"
	"fmt"
)

func main() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	fmt.Print("[32]byte{")

	for i, keyByte := range key {
		if i != 0 {
			fmt.Print(", ")
		}

		fmt.Print(keyByte)
	}

	fmt.Println("}")
}
