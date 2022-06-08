package settings

import (
	"testing"
)

const (
	TOTENAME   = "t1001"
	FEEDPREFIX = "dfense/feeds/"
	IOURL      = "ssl://io.adafruit.com:8883"
	SECUREID   = "dfensenetAmerica"
)

var data = `totename: "t1001"
feedprefix: "dfense/feeds/"
iourl: "ssl://io.adafruit.com:8883"
secureid: "dfensenetAmerica"
username: "dfense"
`

// var data = `totename: "t1001"
// feedprefix: "dfense/feeds/"
// ##
// username: "johnnyfive"
// iourl: "ssl://io.adafruit.com:8883"
// secureid: "dfensenetAmerica"
// `

func TestConfigFromString(t *testing.T) {
	config, err := getConfigFromArray([]byte(data))
	if err != nil {
		t.Error(err)
	}
	if config != nil {
		if config.Totename != TOTENAME {
			t.Error("totename not a match")
		}
		if config.FeedPrefix != FEEDPREFIX {
			t.Error("feedprefix not a match")
		}
		if config.IOURL != IOURL {
			t.Error("iourl not a match")
		}
		if config.Secureid != SECUREID {
			t.Error("secureid not a match")
		}
	} else {
		t.Error("config is nil")
	}
}

func TestReader(t *testing.T) {

	config, err := GetConfigFromFile("/tmp/config.yaml")
	if err != nil {
		t.Errorf("failed reading config: %s", err)
	}

	if config != nil {
		if config.Totename != TOTENAME {
			t.Error("totname not a match")
		}
		if config.FeedPrefix != FEEDPREFIX {
			t.Error("feedprefix not a match")
		}
		if config.IOURL != IOURL {
			t.Error("iourl not a match")
		}
		if config.Secureid != SECUREID {
			t.Error("secureid not a match")
		}
	} else {
		t.Error("config is nil")
	}
}
