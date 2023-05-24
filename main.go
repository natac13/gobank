package main

import (
	"fmt"
	"log"

	dotenv "github.com/joho/godotenv"
)

func main() {

	envErr := dotenv.Load(".env")
	if envErr != nil {
		log.Fatal(envErr)
	}

	store, err := NewPostgresStory()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Store: %+v\n", store)

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
