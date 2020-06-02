/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package main

import (
	"mstats-new/client/config"
	"mstats-new/client/connection"
	"mstats-new/client/dependencies"
	"mstats-new/client/files"
	"mstats-new/client/hardware"
	"mstats-new/client/license"
	"mstats-new/client/updater"
	"mstats-new/logger"
	"time"
)

func main() {

	logger.InitLoggerSystem()
	updater.Init()
	files.InitFilesCheck()
	dependencies.InitDependenciesCheck()
	config.InitStartConfig()
	go license.InitLicense(24 * time.Hour)
	connection.InitClient()
	hardware.SetUpMachine()

	go connection.CheckIfMinerWorks()

	for {
		logger.Info("Holding connection..")
		time.Sleep(5 * time.Second)
	}

}
