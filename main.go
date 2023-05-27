package main

import (
	"flag"
	"fmt"
	"log"

	dotenv "github.com/joho/godotenv"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	return acc

}

func seedAccounts(store Storage) {
	seedAccount(store, "John", "Doe", "admin123")
	seedAccount(store, "Jane", "Doe", "admin123")
	seedAccount(store, "John", "Smith", "admin123")
}

func main() {

	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

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

	if *seed {
		fmt.Println("Seeding the database...")
		seedAccounts(store)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
