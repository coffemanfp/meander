package meander

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// APIKey Google API Key.
var APIKey string

// Query is the user representation of a query to build recommendations.
type Query struct {
	Lat          float64
	Lng          float64
	Journey      []string
	Radius       int
	CostRangeStr string
}

func (q *Query) find(types string) (gR *googleResponse, err error) {
	urlString := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"

	vals := make(url.Values)

	vals.Set("location", fmt.Sprintf("%g,%g", q.Lat, q.Lng))
	vals.Set("radius", fmt.Sprintf("%d", q.Radius))
	vals.Set("types", types)
	vals.Set("key", APIKey)

	if q.CostRangeStr == "" {
		r := ParseCostRange(q.CostRangeStr)

		vals.Set("minprice", fmt.Sprintf("%d", int(r.From)-1))
		vals.Set("maxprice", fmt.Sprintf("%d", int(r.To)-1))
	}

	res, err := http.Get(urlString + "?" + vals.Encode())
	if err != nil {
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&gR)
	return
}

// Run runs the query concurrently, and returns the results.
func (q *Query) Run() (results []interface{}) {
	rand.Seed(time.Now().UnixNano())

	var w sync.WaitGroup
	var l sync.Mutex

	places := make([]interface{}, len(q.Journey))

	for i, typeJourney := range q.Journey {
		w.Add(1)

		go func(types string, i int) {
			defer w.Done()

			response, err := q.find(types)
			if err != nil {
				log.Println("Failed to find places:", err)
				return
			}

			if len(response.Results) == 0 {
				log.Println("No places found for", types)
				return
			}

			for _, result := range response.Results {
				for _, photo := range result.Photos {
					photoURL := fmt.Sprintf(
						"https://maps.googleapis.com/maps/api/place/photo?maxwidth=1000&photoreference=%s&key=%s",
						photo.PhotoRef,
						APIKey,
					)

					photo.URL = photoURL
				}
			}

			randI := rand.Intn(len(response.Results))

			l.Lock()
			places[i] = response.Results[randI]
			l.Unlock()
		}(typeJourney, i)
	}

	w.Wait()
	return
}

// Place is the Google Place model representation.
type Place struct {
	*googleGeometry `json:"geometry"`
	Name            string         `json:"name"`
	Icon            string         `json:"icon"`
	Photos          []*googlePhoto `json:"photos"`
	Vicinity        string         `json:"vicinity"`
}

type googleResponse struct {
	Results []*Place `json:"results"`
}

type googleGeometry struct {
	*googleLocation `json:"location"`
}

type googleLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type googlePhoto struct {
	PhotoRef string `json:"photoRef"`
	URL      string `json:"url"`
}

// Public returns the public version of Place.
func (p *Place) Public() interface{} {
	return map[string]interface{}{
		"name":     p.Name,
		"icon":     p.Icon,
		"photos":   p.Photos,
		"vicinity": p.Vicinity,
		"lat":      p.Lat,
		"lng":      p.Lng,
	}
}
