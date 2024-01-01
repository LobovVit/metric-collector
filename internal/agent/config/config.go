package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Host           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

var instance *Config

func GetConfig() *Config {
	instance = &Config{}
	err := env.Parse(instance)
	if err != nil {
		log.Fatal(err)
	}

	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	reportInterval := flag.Int64("r", 10, "частота отправки метрик на сервер")
	pollInterval := flag.Int64("p", 2, "частота опроса метрик из пакета runtime")
	flag.Parse()

	if instance.Host == "" {
		instance.Host = *host
	}
	if instance.ReportInterval == 0 {
		instance.ReportInterval = *reportInterval
	}
	if instance.PollInterval == 0 {
		instance.PollInterval = *pollInterval
	}
	return instance
}
