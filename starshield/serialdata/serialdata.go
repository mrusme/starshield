package serialdata

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

type SerialData struct {
	Error               string  `json:"err,omitempty"`
	ErrorCode           int     `json:"errc,omitempty"`
	IAQ                 float64 `json:"iaq,omitempty"`
	IAQAccuracy         int     `json:"iaqa,omitempty"`
	CO2Equivalent       float64 `json:"co2e,omitempty"`
	Temperature         float64 `json:"temp,omitempty"`
	TemperatureComp     float64 `json:"temc,omitempty"`
	Pressure            float64 `json:"pres,omitempty"`
	Humidity            float64 `json:"humy,omitempty"`
	HumidityComp        float64 `json:"humc,omitempty"`
	GasResistance       float64 `json:"gasr,omitempty"`
	GasPercentage       float64 `json:"gasp,omitempty"`
	StanilizationStatus int     `json:"stbs,omitempty"`
	RunInStatus         int     `json:"runs,omitempty"`
	VisibleLight        int     `json:"visb,omitempty"`
	Infrared            int     `json:"infr,omitempty"`
	UVIndex             float64 `json:"ulvi,omitempty"`

	mtx sync.Mutex
}

func (sd *SerialData) UpdateFrom(usd *SerialData) {
	sd.mtx.Lock()

	aVal := reflect.ValueOf(sd).Elem()
	aTyp := aVal.Type()
	bVal := reflect.ValueOf(*usd)
	for i := 0; i < aVal.NumField(); i++ {
		if aTyp.Field(i).Name == "mtx" {
			continue
		}
		aVal.Field(i).Set(bVal.Field(i))
	}

	sd.mtx.Unlock()
}

func (sd *SerialData) Print() {
	sd.mtx.Lock()
	fmt.Printf("%v\n", sd)
	sd.mtx.Unlock()
}

func (sd *SerialData) ToJSON() []byte {
	sd.mtx.Lock()
	jsonResponse, err := json.Marshal(sd)
	sd.mtx.Unlock()
	if err != nil {
		return []byte("")
	}
	return jsonResponse
}

func (sd *SerialData) FromJSON(b []byte) error {
	if err := json.Unmarshal(b, sd); err != nil {
		return err
	}

	return nil
}

func New(b []byte) (*SerialData, error) {
	sd := new(SerialData)
	if err := sd.FromJSON(b); err != nil {
		return nil, err
	}

	return sd, nil
}
