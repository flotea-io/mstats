package config

var generalIp = "0.0.0.0"
var generalPort = "8088"
var generalPassword = "test_password"
var version = "v0.8"

func GetStringAddress() string {
	return generalIp + ":" + generalPort
}

func GetPassword() string {
	return generalPassword
}

func IsCurrentVersion(s string) bool {
	return version == s
}

func GetLicense() string {
	return license
}
