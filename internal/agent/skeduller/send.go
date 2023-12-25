package skeduller

import (
	"fmt"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"net/http"
	"os"
	"strconv"
)

func sendRequest(metrics *metrics.Metrics, endPoint string) {
	for k, v := range metrics.Gauge {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%vgauge/%v/%v", endPoint, k, strconv.FormatFloat(v, 'f', 10, 64)), nil)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "text/plain")
		/*resp*/ _, err = client.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}
		//fmt.Println(resp.Status)
	}
	for k, v := range metrics.Counter {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%vcounter/%v/%v", endPoint, k, strconv.FormatInt(v, 10)), nil)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)

		}
		req.Header.Set("Content-Type", "text/plain")
		/*resp*/ _, err = client.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}
		//fmt.Println(resp.Status)
		//fmt.Println(req.URL)
	}
}
