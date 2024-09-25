package main

import (
	"github.com/ViniciusDSLima/golang01/cmd/api"
	"github.com/ViniciusDSLima/golang01/config"
	"github.com/ViniciusDSLima/golang01/db"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	cfg := config.Config{
		Host:     config.Env.Host,
		User:     config.Env.User,
		Password: config.Env.Password,
		DBName:   config.Env.DBName,
		Port:     config.Env.Port,
		SSLMode:  config.Env.SSLMode,
	}

	db, err := db.NewPostgresStorage(cfg)

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":9090", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
