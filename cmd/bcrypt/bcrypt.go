package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Invalid command: %v\n", os.Args[1])
	}
}

func hash(password string) {
	fmt.Printf("Hashing the password: %s\n", password)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(hashedBytes))
}

func compare(password, hash string) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		fmt.Println("Passwords do not match")
		return
	}
	fmt.Println("Passwords match")
}
