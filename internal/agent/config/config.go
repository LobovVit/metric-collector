package config

import "flag"

type Config struct {
	Host           string //Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	ReportInterval int64  //Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	PollInterval   int64  //Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
}

var instance *Config

func GetConfig() *Config {
	instance = &Config{}
	flag.StringVar(&instance.Host, "a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	flag.Int64Var(&instance.ReportInterval, "r", 10, "частота отправки метрик на сервер")
	flag.Int64Var(&instance.PollInterval, "p", 2, "частота опроса метрик из пакета runtime")
	// делаем разбор командной строки
	flag.Parse()
	return instance
}
