package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Host string `env:"ADDRESS"`
}

var instance *Config

func GetConfig() *Config {
	instance = &Config{}
	err := env.Parse(instance)
	if err != nil {
		log.Fatal(err)
	}
	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	flag.Parse()

	if instance.Host == "" {
		instance.Host = *host
	}
	return instance
}
