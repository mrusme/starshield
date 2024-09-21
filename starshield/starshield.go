package main

import (
	"log"
	"os"
	"strings"

	"github.com/mrusme/starshield/serialdata"
	"go.bug.st/serial"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No serial port given!")
	}

	sport := os.Args[1]

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	sportFound := false
	for _, port := range ports {
		if port == sport {
			sportFound = true
		}
	}
	if !sportFound {
		log.Fatal("Given serial port not found!")
	}

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(sport, mode)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var sdjs string

		for {
			buff := make([]byte, 128)
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				break
			}
			sdjs += string(buff[:n])

			if strings.Contains(string(buff[:n]), "\n") {
				break
			}
		}
		log.Printf("%v", sdjs)
		sd, err := serialdata.New([]byte(sdjs))
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%v\n", sd)
	}
}
