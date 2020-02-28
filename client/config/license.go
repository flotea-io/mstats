package config

var licenseServerIP = ""
var licenseServerPort = ""

func GetLicenseAddress() string {
	return licenseServerIP + ":" + licenseServerPort
}
