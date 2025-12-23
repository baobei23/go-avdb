package main

import (
	"log"

	"github.com/baobei23/go-avdb/internal/api"
	"github.com/baobei23/go-avdb/internal/env"
	"github.com/baobei23/go-avdb/internal/store"
)

func main() {
	cfg := api.Config{
		Port: env.GetString("PORT", ":8080"),
	}
	store := store.NewStorage(nil)

	app := api.Application{
		Config: cfg,
		Store:  store,
	}

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
