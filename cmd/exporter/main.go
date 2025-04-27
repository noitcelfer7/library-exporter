package main

import (
	"encoding/json"
	"fmt"
	"os"

	"library_exporter/internal/exporter/config"
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

	rpc.Serve(&config)
}
