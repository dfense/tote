package drv

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestBMEjsonMarshal(t *testing.T) {

	testString := `{
		"IAQ": 25.00,
		"IAQ_Accuracy": 0,
		"Temperature_C": 24.56,
		"Humidity_%": 56.30,
		"Pressure_hPa": 993.50,
		"Altitude_m": 204.62,
		"Raw_Temperature_C": 24.56,
		"Raw_Humidity_%": 56.30,
		"VOC_Ohm": 37700.00,
		"Static_IAQ_Index": 25.00,
		"CO2_Equivalent_PPM": 500.00,
		"Breath_VOC_Equivalent": 0.50,
		"Stabilization_Status": 0
	}`

	bme := NewBME680(PRIMARY)
	err := json.Unmarshal([]byte(testString), &bme)
	if err != nil {
		t.Error(err)
	}

	// validate some fields below
	fmt.Printf("%+v\n", bme)
}

func TestBMEjsonUnmarshal(t *testing.T) {
	bme := NewBME680(PRIMARY)
	bme.IAQ = 45.90
	bme.IAQaccuracy = 1

	jsonBytes, err := json.Marshal(bme)
	if err != nil {
		t.Error(err)
	}

	// validate some fields
	fmt.Println(string(jsonBytes))
}

func TestHasPrefix(t *testing.T) {
	s := `{"IAQ": 95.55,"IAQ_Accuracy": 0,"Temperature_C": 23.91,"Humidity_%": 55.77,"Pressure_hPa": 996.05,"Altitude_m": 183.09,"Raw_Temperature_C": 24.02,"Raw_Humidity_%": 55.55,"VOC_Ohm": 29504.00,"Static_IAQ_Index": 55.82,"CO2_Equivalent_PPM": 623.27,"Breath_VOC_Equivalent": 0.80,"Stabilization_Status": 0}`
	log.Println(strings.HasPrefix(s, `{`))
}
