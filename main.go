package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func getGymData() {
	studio := viper.Get("studio").(int)
	req, err := http.NewRequest("GET", "https://my.mcfit.com/nox/public/v1/studios/"+fmt.Sprint(studio)+"/utilization/v2/today", nil)
	if err != nil {
		log.Println("Error while building request: ")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Nox-Client-Type", "WEB")
	req.Header.Set("X-Nox-Web-Context", "v=1")
	req.Header.Set("X-Tenant", "rsg-group")
	req.Header.Set("Alt-Used", "my.mcfit.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "countryCode=de")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Te", "trailers")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error while doing request: ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading response body: ", err)
	}

	var data []ResponseItem

	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Error while unmarshaling JSON: ", err)
	}

	var percentage int
	for _, item := range data {
		if item.Current {
			percentage = item.Percentage
			break
		}
	}
	percentageMetric.Set(float64(percentage))
	log.Printf("Current utilization: %s", fmt.Sprint(percentage))

}

type ResponseItem struct {
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Current    bool   `json:"current"`
	Percentage int    `json:"percentage"`
}

var (
	percentageMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "current_percentage",
		Help: "Current utilization percentage",
	})
)

func main() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Could not read config file: %w", err))
	}

	prometheus.MustRegister(percentageMetric)

	go func() {
		for {
			getGymData()
			time.Sleep(15 * time.Minute)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
