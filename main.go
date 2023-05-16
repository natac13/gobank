package main

import (
	"fmt"
	"log"
)

func main() {
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
