package net

import "testing"

func TestFileWriter(t *testing.T) {
	wifi := WiFi{SSID: "dfensenet", Wpa_Password: "gofish"}
	err := WiFiCredentials(wifi)
	if err != nil {
		t.Error(err)
	}
}
