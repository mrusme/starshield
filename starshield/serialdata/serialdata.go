package serialdata

import "encoding/json"

type SerialData struct {
	Error               string  `json:"err,omitempty"`
	ErrorCode           int     `json:"errc,omitempty"`
	IAQ                 int     `json:"iaq,omitempty"`
	IAQAccuracy         int     `json:"iaqa,omitempty"`
	Temperature         float64 `json:"temp,omitempty"`
	Pressure            float64 `json:"pres,omitempty"`
	Humidity            float64 `json:"humy,omitempty"`
	GasResistance       float64 `json:"gasr,omitempty"`
	StanilizationStatus int     `json:"stbs,omitempty"`
	RunInStatus         int     `json:"runs,omitempty"`
	VisibleLight        int     `json:"visb,omitempty"`
	Infrared            int     `json:"infr,omitempty"`
	UVIndex             float64 `json:"ulvi,omitempty"`
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
