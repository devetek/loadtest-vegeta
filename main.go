package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	err := godotenv.Load(".env.example")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rate := vegeta.Rate{Freq: 300, Per: time.Second}
	duration := 10 * time.Second

	uaTarget := os.Getenv("LOADTEST_TARGET_UA")

	headers := http.Header{}
	headers.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	if uaTarget == "user" {
		headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	}

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: os.Getenv("LOADTEST_TARGET_METHOD"),
		URL:    os.Getenv("LOADTEST_TARGET_URL"),
		Header: headers,
	})

	attacker := vegeta.NewAttacker(vegeta.Timeout(5 * time.Second))

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	println("===== RESULT ======")

	println("HTTP Response Status Codes")
	for status, count := range metrics.StatusCodes {
		fmt.Printf("%s: %d\n", status, count)
	}

	println("")

	if len(metrics.Errors) > 0 {
		println("Error when attack:")
		for i, error := range metrics.Errors {
			println(i, error)
		}
		println("")
	}

	fmt.Printf("User Agent: %s\n", uaTarget)
	fmt.Printf("Total Request: %d\n", metrics.Requests)
	fmt.Printf("Rate: %.2f\n", metrics.Rate)
	fmt.Printf("Duration: %s\n", metrics.Duration)
	fmt.Printf("Success Rate: %.2f%%\n", metrics.Success)
	fmt.Printf("50th percentile: %s\n", metrics.Latencies.P50)
	fmt.Printf("90th percentile: %s\n", metrics.Latencies.P90)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}
