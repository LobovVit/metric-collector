package actions

import (
	"errors"
	"strconv"
)

func CheckAndSave(tp string, name string, value string) error {

	switch tp {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("Bad Request")
		}
		store.SetGauge(name, v)
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errors.New("Bad Request")
		}
		store.SetCounter(name, v)
	default:
		return errors.New("Bad Request")
	}
	return nil
}
