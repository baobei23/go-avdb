package api

import "github.com/baobei23/go-avdb/internal/store"

type Application struct {
	Config Config
	Store  store.Storage
}

type Config struct {
	Port string
}
