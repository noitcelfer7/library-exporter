package main

import (
	"encoding/json"
	"fmt"
	"os"

	"library_exporter/internal/exporter/config"
	postgresql "library_exporter/internal/exporter/database"
	"library_exporter/internal/exporter/http"
	"library_exporter/internal/exporter/rpc"
)

func main() {
	data, err := os.ReadFile("config.json")

	if err != nil {
		panic(fmt.Sprintf("os.ReadFile Error: %v", err))
	}

	var config = config.Config{}

	err = json.Unmarshal(data, &config)

	if err != nil {
		panic(fmt.Sprintf("json.Unmarshal Error: %v", err))
	}

	db, err := postgresql.NewDB(config.Postgresql.Url)

	if err != nil {
		panic(fmt.Sprintf("postgresql.NewDB Error: %v", err))
	}

	defer db.Close()

	go http.Serve(&config, db)

	rpc.Serve(&config, db)
}
