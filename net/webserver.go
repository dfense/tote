package net

import "os"

type WiFi struct {
	SSID         string
	Wpa_Password string
}

// pass in creds to create a new file and put in /boot
func WiFiCredentials(wifi WiFi) error {

	f, err := os.Open("/tmp/dat")
	if err != nil {
		return err
	}

	f.WriteString("hello, is it me you're looking for")
	return nil
}

// set credentials and create a new file
// @POST
// func setWiFiHandler(http http.HttpHandler) {
// 	// parse json out
// 	// call WiFiCredentials
// }

// return all ifconifg structs
// @GET
func GetNetworkConfig() {

}

func StartServer(port int) {

	//serveHTTP

}
