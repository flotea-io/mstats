/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

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
