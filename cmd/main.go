package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/qlfzn/parcel2u/cmd/api"
	"github.com/qlfzn/parcel2u/internal/db"
)

func main() {
	godotenv.Load()

	addr, ok := os.LookupEnv("DB_ADDR")
	if !ok {
		log.Fatal("failed to load env")
	}

	db, err := db.New(addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("database connection established")

	app := &api.Application{
		Addr: ":8080",
	}

	mux := app.Mount()
	app.Run(mux)
}
