package main

import (
	"github.com/qlfzn/parcel2u/cmd/api"
)

func main() {
	app := &api.Application{
		Addr: ":8080",
	}

	mux := app.Mount()
	app.Run(mux)
}
