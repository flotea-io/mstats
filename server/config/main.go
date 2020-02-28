package config

import (
	"flag"
	"mstats-new/logger"
)

func InitStartConfig() {
	logger.Info("I am checking configuration params..")
	portFlag := flag.String("generalPort", "8088", "set server port")
	passwordFlag := flag.String("generalPassword", "test_password", "set server password")
	webPortFlag := flag.String("webPort", "8099", "set web server port")
	licenseServerIPFlag := flag.String("licenseServerIP", "146.185.154.46", "set license server ip")
	licenseServerPortFlag := flag.String("licenseServerPort", "9940", "set license server port")
	databaseNameFlag := flag.String("databaseName", "bitforce", "set database name")
	databasePasswordFlag := flag.String("databasePassword", "", "set database password")
	databaseUserFlag := flag.String("databaseUser", "root", "set database user")
	databaseHostFlag := flag.String("databaseHost", "127.0.0.1", "set database host")
	flag.Parse()

	generalPort = *portFlag
	generalPassword = *passwordFlag
	webPort = *webPortFlag
	licenseServerIP = *licenseServerIPFlag
	licenseServerPort = *licenseServerPortFlag
	databaseName = *databaseNameFlag
	databaseHost = *databaseHostFlag
	databaseUser = *databaseUserFlag
	databasePassword = *databasePasswordFlag

	logger.Info("Configuration params loaded..")
	logger.Info("===========================")
	logger.Info("General server Port: " + generalPort)
	logger.Info("Password to server: " + generalPassword)

	if generalPassword == "test_password" {
		logger.Warning("Your generalPassword is not set! Please set your own generalPassword by param --generalPassword")
	}
	logger.Info("Web server port: " + webPort)
	logger.Info("License: " + license)
	logger.Info("License server: " + licenseServerIP + ":" + licenseServerPort)

	logger.Info("Database host " + databaseHost)
	logger.Info("Database user " + databaseUser)
	logger.Info("Database name " + databaseName)
	logger.Info("Server app version: " + version)
	logger.Info("===========================")

}
