package config

import (
	"flag"
	"mstats-new/logger"
	"os"
	"regexp"
)

var ip = "127.0.0.1"
var port = "8080"
var password = "test_password"
var clientName = "test_client"
var version = "v0.14"

func GetPassword() string {
	return password
}

func GetStringAddress() string {
	return ip + ":" + port
}

func GetClientName() string {
	return clientName
}

func InitStartConfig() {
	logger.Info("I am checking configuration params..")

	portFlag := flag.String("port", "8088", "set server port")
	passwordFlag := flag.String("password", "test_password", "set password to server")
	ipFlag := flag.String("ip", "127.0.0.1", "set server ip")
	clientNameFlag := flag.String("name", "test_client", "set name of client")
	licenseServerIPFlag := flag.String("licenseServerIP", "123.123.123.123", "set license server ip")
	licenseServerPortFlag := flag.String("licenseServerPort", "9940", "set license server port")
	flag.Parse()

	port = *portFlag
	password = *passwordFlag
	ip = *ipFlag
	clientName = *clientNameFlag
	licenseServerIP = *licenseServerIPFlag
	licenseServerPort = *licenseServerPortFlag

	logger.Info("Configuration params loaded..")
	logger.Info("===========================")
	logger.Info("Client name: " + clientName)
	logger.Info("Socket IP connect: " + ip)
	logger.Info("Socket port connect: " + port)
	logger.Info("Password to server: " + password)
	logger.Info("Client app version: " + version)
	logger.Info("===========================")

	validateName(clientName)

}

func IsCurrentVersion(ver string) bool {
	return version == ver
}

func validateName(name string) {

	re := regexp.MustCompile("^[a-zA-Z0-9_]{3,30}$").MatchString
	if !re(name) {
		logger.Error("Name can contain only a-z, A-Z, 0-9 and _ characters and must be 3-30 characters long")
		os.Exit(1)
	}

}
