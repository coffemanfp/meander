package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/coffemanfp/meander"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	apiKey := os.Getenv("GOOGLE_APIKEY")
	if apiKey == "" {
		log.Fatalln("google: not api key found")
	}

	meander.APIKey = apiKey

	http.HandleFunc("/journeys", getJourneys)

	http.HandleFunc("/recommendations", getRecommendations)

	fmt.Println("Starting on port 8080...")
	http.ListenAndServe(":8080", http.DefaultServeMux)
}
