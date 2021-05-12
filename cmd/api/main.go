package main

import (
	"log"

	"github.com/Metalisaac321/practice-concurrency/cmd/api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
