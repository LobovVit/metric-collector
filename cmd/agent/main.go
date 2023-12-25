package main

import "github.com/LobovVit/metric-collector/internal/agent/skeduller"

const endPoint = "http://localhost:8080/update/"
const readTime int64 = 2
const sendTime int64 = 10

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	skeduller.StartTimer(readTime, sendTime, endPoint)
	return nil
}
