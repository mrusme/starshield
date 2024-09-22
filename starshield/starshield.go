package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mrusme/starshield/reader"
	"github.com/mrusme/starshield/serialdata"
	"go.bug.st/serial"
)

var STATE serialdata.SerialData

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
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

func httpServer(sd *serialdata.SerialData) {
	http.HandleFunc("/", handler(sd))
	log.Fatal(http.ListenAndServe("127.0.0.1:3232", nil))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No serial port given!")
	}

	sport := os.Args[1]

	for {
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}

		if len(ports) > 0 {
			sportFound := false
			for _, port := range ports {
				if port == sport {
					sportFound = true
				}
			}
			if sportFound {
				break
			}
		}

		time.Sleep(10 * time.Second)
	}

	STATE = serialdata.SerialData{}

	go reader.Reader(sport, &STATE)
	go httpServer(&STATE)

	for {
		STATE.Print()
		time.Sleep(time.Second * 5)
	}
}
