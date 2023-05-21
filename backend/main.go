package main

import (
	"flag"
	"log"

	"github.com/eqimd/transbyte-site/internal/apiserver"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "config", "configs/config.toml", "path to server config")

	flag.Parse()

	cfg, err := apiserver.NewConfig(configPath)
	if err != nil {
		log.Fatal("Could not parse toml config:", err)
	}

	server := apiserver.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
