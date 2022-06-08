package settings

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type Config struct {
	Totename   string
	FeedPrefix string
	IOURL      string
	Secureid   string
	Username   string
	B          struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

// read config file
func GetConfigFromFile(filename string) (*Config, error) {

	data, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	log.Println(string(data))
	return getConfigFromArray(data)
}

// parse yaml from a []byte
func getConfigFromArray(data []byte) (*Config, error) {

	config := &Config{}
	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}
	return config, err
}

// read yaml config from file location provided
func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
