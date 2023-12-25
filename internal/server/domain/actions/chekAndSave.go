package actions

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/storage"
	"net/http"
	"strconv"
	"strings"
)

var store = StorInterface(storage.GetStorage())

func CheckAndSave(url string) int {
	part := strings.Split(url, `/`)
	if len(part) != 5 && len(part) != 6 {
		return http.StatusNotFound
	}
	switch part[2] {
	case "gauge":
		v, err := strconv.ParseFloat(part[4], 64)
		if err != nil {
			return http.StatusBadRequest
		}
		store.SetGauge(part[3], v)
	case "counter":
		v, err := strconv.ParseInt(part[4], 10, 64)
		if err != nil {
			return http.StatusBadRequest
		}
		store.SetCounter(part[3], v)
	default:
		return http.StatusBadRequest
	}
	return http.StatusOK
}
