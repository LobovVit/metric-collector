package handlers

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"html/template"
	"net/http"
	"strings"
)

func allMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	res, err := actions.GetAll()
	data, errData := createHTML(mapToMetric(res))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		if errData != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}

	}
}

type Metric struct {
	Tp  string
	Key string
	Val string
}

type PageData struct {
	Title string
	Data  Metric
}

func mapToMetric(m map[string]map[string]string) []Metric {
	ret := []Metric{}
	for k, v := range m {
		for k2, v2 := range v {
			item := Metric{Tp: k, Key: k2, Val: v2}
			ret = append(ret, item)
		}
	}
	return ret
}

func createHTML(data []Metric) ([]byte, error) {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Metrics</title>
	</head>
	<body>
		<table border="1">
		<tr>
			<th>type</th>
			<th>name</th>
			<th>value</th>
		</tr>
		{{range . }}<tr><td>{{ .Tp }}</td><td>{{ .Key }}</td><td>{{ .Val }}</td></tr>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`
	tmpl, err := template.New("AllMetrics").Parse(tpl)
	if err != nil {
		return []byte("Err Parse" + err.Error()), err
	}
	reader := new(strings.Builder)
	err = tmpl.Execute(reader, data)
	if err != nil {
		return []byte("Err Execute" + err.Error()), err
	}
	return []byte(reader.String()), nil
}
