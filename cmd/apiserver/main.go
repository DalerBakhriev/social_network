package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/DalerBakhriev/social_network/internal/app/apiserver"
)

func main() {

	var configPath string
	flag.StringVar(&configPath, "config-path", "./configs/apiserver.toml", "path to config file")
	flag.Parse()

	config := apiserver.NewConfig()

	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatalf("Failed to parse config file %s: %v", configPath, err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
