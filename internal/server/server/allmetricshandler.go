package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// allMetricsHandler - handler returning all server data in html format
func (a *Server) allMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	res, err := a.storage.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := createHTML(mapToMetric(res))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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
		return []byte("Err Parse" + err.Error()), fmt.Errorf("parse template: %w", err)
	}
	writer := new(strings.Builder)
	err = tmpl.Execute(writer, data)
	if err != nil {
		return []byte("Err Execute" + err.Error()), fmt.Errorf("execute template: %w", err)
	}
	return []byte(writer.String()), nil
}
