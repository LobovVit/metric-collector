package actions

import (
	"fmt"
	"strconv"
)

func badRequestErr(tp, value string) error {
	return fmt.Errorf("bad request metric type:\"%v\" with value:\"%v\"", tp, value)
}

func CheckAndSave(tp string, name string, value string) error {

	switch tp {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return badRequestErr(tp, value)
		}
		store.storage.SetGauge(name, v)
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return badRequestErr(tp, value)
		}
		store.storage.SetCounter(name, v)
	default:
		return badRequestErr(tp, value)
	}
	return nil
}
