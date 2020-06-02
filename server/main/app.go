/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package main

import (
	"mstats-new/logger"
	"mstats-new/server/config"
	"mstats-new/server/connection"
	"mstats-new/server/database"
	"mstats-new/server/license"
	"mstats-new/server/updater"
	"mstats-new/server/web"
	"time"
)

func main() {
	logger.InitLoggerSystem()
	updater.Init()
	config.InitStartConfig()
	database.InitDatabase()
	go license.InitLicense(1*time.Hour, "_server")
	go connection.InitSocketServer()
	database.CleanDatabase()
	go connection.InitClientAutoreset()
	web.InitWebServer()
	for {
		logger.Info("I am holding connection..")
		time.Sleep(10 * time.Second)
	}

}
