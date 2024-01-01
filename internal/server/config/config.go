package config

import "flag"

type Config struct {
	Host string //Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).

}

var instance *Config

func GetConfig() *Config {
	instance = &Config{}
	flag.StringVar(&instance.Host, "a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	// делаем разбор командной строки
	flag.Parse()
	return instance
}
