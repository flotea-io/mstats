/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package programs

import (
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/server/config"
	"mstats-new/server/connection"
	"os"

	"github.com/jinzhu/gorm"
)

func StartClaymore(machine string, cfg string) {
	var cfgBasic internal.ClayMoreBasicConfig
	var cfgExport internal.ClayMoreBasicConfig

	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		logger.Error("Problem with connect to database..")
		logger.Error(err.Error())
		os.Exit(1)
		return
	}

	data.Where("name = ?", "export").First(&cfgExport)
	data.Where("name = ?", cfg).First(&cfgBasic)

	if cfgBasic.Params == "" {
		data.Where("name = ?", "default").First(&cfgBasic)
	}

	connection.SendMessageToClient(machine, internal.CreatePacket("server", cfgExport.ToJson(), internal.RegisterExportsPacket, config.GetPassword()).ToJson())
	//time.Sleep(500000000 * time.Nanosecond)
	connection.SendMessageToClient(machine, internal.CreatePacket("server", cfgBasic.ToJson(), internal.RunClaymorePacket, config.GetPassword()).ToJson())
}

func StopClaymore(machine string) {
	connection.SendMessageToClient(machine, internal.CreatePacket("server", "", internal.StopClaymorePacket, config.GetPassword()).ToJson())
}
