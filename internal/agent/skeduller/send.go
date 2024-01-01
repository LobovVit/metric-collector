package skeduller

import (
	"fmt"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
)

var client = resty.New()

func sendRequest(metrics *metrics.Metrics, endPoint string) {
	for k, v := range metrics.Gauge {
		//client := &http.Client{}
		//req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%vgauge/%v/%v", endPoint, k, strconv.FormatFloat(v, 'f', 10, 64)), nil)
		_, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(fmt.Sprintf("%vgauge/%v/%v", endPoint, k, strconv.FormatFloat(v, 'f', 10, 64)))
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}
		//req.Header.Set("Content-Type", "text/plain")
		//resp, err := client.Do(req)
		//if err != nil {
		//	fmt.Printf("client: error making http request: %s\n", err)
		//	os.Exit(1)
		//}
		//defer resp.Body.Close()
		//fmt.Println(resp.Status)
	}
	for k, v := range metrics.Counter {
		//client := &http.Client{}
		//req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%vcounter/%v/%v", endPoint, k, strconv.FormatInt(v, 10)), nil)
		//if err != nil {
		//	fmt.Printf("client: error making http request: %s\n", err)
		//	os.Exit(1)
		//
		//}
		//req.Header.Set("Content-Type", "text/plain")
		//resp, err := client.Do(req)
		//if err != nil {
		//	fmt.Printf("client: error making http request: %s\n", err)
		//	os.Exit(1)
		//}
		//defer resp.Body.Close()
		//fmt.Println(resp.Status)
		//fmt.Println(req.URL)
		_, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(fmt.Sprintf("%vcounter/%v/%v", endPoint, k, strconv.FormatInt(v, 10)))
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}
	}
}
