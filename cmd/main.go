package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

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

	http.HandleFunc("/recommendations", func(w http.ResponseWriter, r *http.Request) {
		journey := strings.Split(r.URL.Query().Get("journey"), "|")
		if len(journey) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no journey data found."))
			return
		}

		q := &meander.Query{
			Journey: journey,
		}

		var err error

		q.Lat, err = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid latitude"))
			return
		}

		q.Lng, err = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid longitude"))
			return
		}

		q.Radius, err = strconv.Atoi(r.URL.Query().Get("radius"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid radius"))
			return
		}

		q.CostRangeStr = r.URL.Query().Get("cost")

		places := q.Run()

		respond(w, r, places)
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
