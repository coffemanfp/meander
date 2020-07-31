package main

import (
	"encoding/json"
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

	http.HandleFunc("/journeys", func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	})

	http.ListenAndServe(":8080", http.DefaultServeMux)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {

	publicData := make([]interface{}, len(data))

	for i, d := range data {
		publicData[i] = meander.Public(d)
	}

	return json.NewEncoder(w).Encode(publicData)
}
