package reader

import (
	"log"
	"strings"
	"time"

	"github.com/mrusme/starshield/serialdata"
	"go.bug.st/serial"
)

func Reader(sport string, state *serialdata.SerialData) {
	var port serial.Port
	var err error
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	for {
		port, err = serial.Open(sport, mode)
		if err == nil {
			break
		}
		time.Sleep(30 * time.Second)
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
		// log.Printf("%v", sdjs)
		sd, err := serialdata.New([]byte(sdjs))
		if err != nil {
			// log.Println(err)
			continue
		}
		// log.Printf("%v\n", sd)

		state.UpdateFrom(sd)
	}
}
