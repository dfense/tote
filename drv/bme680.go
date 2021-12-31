package drv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type i2cAddress string

const (
	PRIMARY   i2cAddress = "BME680_I2C_ADDR_PRIMARY"
	SECONDARY            = "BME680_I2C_ADDR_SECONDARY"
)

type BME680 struct {
	IAQ          float64    `json:"IAQ"`
	IAQaccuracy  int        `json:"IAQ_Accuracy"`
	TemperatureC float64    `json:"Temperature_C"`
	RH           float64    `json:"Humidity_%"`
	Pressure     float64    `json:"Pressure_hPa"`
	Altitude     float64    `json:"Altitude_m"`
	RawTempC     float64    `json:"Raw_Temperature_C"`
	RawRH        float64    `json:"Raw_Humidity_%"`
	VOC          float64    `json:"VOC_Ohm"`
	SIAQ         float64    `json:"Static_IAQ_Index"`
	Co2PPM       float64    `json:"CO2_Equivalent_PPM"`
	BVOCe        float64    `json:"Breath_VOC_Equivalent"`
	Status       int        `json:"Stabilization_Status"`
	I2c          i2cAddress `json:"-"`
	pid          int        `json:"-"`
}

func NewBME680(i2cLocal i2cAddress) *BME680 {
	return &BME680{I2c: i2cLocal}
}

// Start runs the bsec_bme680 binary and streams (emit) the stdout over the supplied channel
// Stop will kill the process and will send (emit) a "EOF" string over the channel
func (bme *BME680) Start(outC chan<- interface{}) {

	fmt.Printf("bme : %s\n", bme.I2c)
	// The command you want to run along with the argument
	cmd := exec.Command("./bsec_bme680", string(bme.I2c))
	// ./bsec_bme680 <BME680_I2C_ADDR_PRIMARY> or <BME680_I2C_ADDR_SECONDARY>

	if cmd.Process != nil {
		fmt.Println("========== process not nil =====")
	}
	// Get a pipe to read from standard out
	r, _ := cmd.StdoutPipe()

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {

		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "{") {
				err := json.Unmarshal([]byte(line), bme)
				if err != nil {
					log.Printf("Error json: %s\n", err)
				}
				outC <- bme
			} else {
				log.Println("no prefix match")
				log.Printf("-< %s\n", line)
			}
		}

		outC <- "EOF" + bme.I2c

	}()

	// Start the command and check for errors
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
	if cmd.Process != nil {
		bme.pid = cmd.Process.Pid
		log.Printf("process id ====> %d\n", bme.pid)
	}

}

func (bme BME680) Stop() {
	log.Println("stopping bme")
	if bme.pid > 0 {

		fmt.Printf("PID: %d will be killed.\n", bme.pid)
		proc, err := os.FindProcess(bme.pid)
		if err != nil {
			log.Println(err)
		}
		// Kill the process
		err = proc.Kill()
		if err != nil {
			log.Printf("error killing : %s\n", err)
		}

		log.Println("kill zone complete")

	}
}
