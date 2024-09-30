package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/mrusme/starshield/reader"
	"github.com/mrusme/starshield/serialdata"
)

var STATE serialdata.SerialData

type LocationResponse struct {
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lng"`
	} `json:"location"`
	Accuracy float64 `json:"accuracy"`
}

func handler(sd *serialdata.SerialData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

			w.Header().Set("Content-Type", "application/json")

			jsonResponse := sd.ToJSON()
			w.Write(jsonResponse)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func handlerLocation(sd *serialdata.SerialData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			w.Header().Set("Content-Type", "application/json")

			lat, lon, _, acc := sd.GetLatLonAltAcc()
			locResp := LocationResponse{}
			locResp.Location.Latitude = (float64(lat) * math.Pow(10, -7))
			locResp.Location.Longitude = (float64(lon) * math.Pow(10, -7))
			locResp.Accuracy = (float64(acc) / 1000.0)

			jsonResponse, err := json.Marshal(locResp)
			if err != nil {
				http.Error(w, "JSON error", http.StatusInternalServerError)
			}
			w.Write(jsonResponse)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func httpServer(sd *serialdata.SerialData) {
	http.HandleFunc("/", handler(sd))
	http.HandleFunc("/location", handlerLocation(sd))
	log.Fatal(http.ListenAndServe("127.0.0.1:3232", nil))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No serial port given!")
	}

	sport := os.Args[1]

	STATE = serialdata.SerialData{}

	go reader.Reader(sport, &STATE)
	go httpServer(&STATE)

	for {
		time.Sleep(time.Second * 5)
	}
}
