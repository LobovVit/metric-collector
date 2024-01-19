package actions

import (
	"fmt"
	"strconv"
)

type badRequestErr struct {
	tp    string
	value string
}

func (e badRequestErr) Error() string {
	return fmt.Sprintf("bad request metric type:\"%v\" with value:\"%v\"", e.tp, e.value)
}

func (r Repo) CheckAndSave(tp string, name string, value string) error {
	switch tp {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		r.storage.SetGauge(name, v)
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		r.storage.SetCounter(name, v)
	default:
		return badRequestErr{tp, value}
	}
	return nil
}
