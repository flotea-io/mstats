package config

import "strconv"

var databaseName = "bitforce"
var databaseHost = "localhost"
var databasePort = 3306
var databaseUser = "root"
var databasePassword = ""

func GetDatabaseDSN() string {
	return databaseUser + ":" + databasePassword + "@tcp(" + databaseHost + ":" + strconv.Itoa(databasePort) + ")/" + databaseName
}
