package main

import (
	"fmt"
	"github.com/northwindman/testREST-autentification/internal/lib/random"
)

func main() {

	for i := 20; i <= 25; i++ {
		token, err := random.NewSecret(i)
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
	}

}
