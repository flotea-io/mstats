/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package connection

import (
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/server/config"
	"time"

	"github.com/jinzhu/gorm"
)

func InitClientAutoreset() {
	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		logger.Error("Can't connect to database..")
		return
	}

	//wait for connects and initialize
	time.Sleep(1 * time.Minute)

	for {

		clients := make(map[string]internal.Miner)
		clients = GetConnectedClients()

		connected := make([]string, 0, len(clients))

		for name := range clients {
			connected = append(connected, name)
		}

		var hardResets []internal.HardReset
		var names []string

		fifteenMinutesAgo := time.Now().Add(-15 * time.Minute).Format("2006-01-02 15:04:05")

		db.Table("claymore_stats").Group("miner_name").Where("miner_name IN (?)", connected).Having("MAX(time) < (?)", fifteenMinutesAgo).Pluck("miner_name", &names)
		db.Table("pin_functions").Where("miner_name IN (?) AND function = ? AND auto_reset = ?", names, 1, 1).Select("id as 'machine_number'").Scan(&hardResets)

		for _, hardReset := range hardResets {
			db.Save(&hardReset)
			logger.Warning("Machine on pin " + hardReset.MachineNumber + " is zombie. Queued to hard reset..")
		}

		time.Sleep(15 * time.Minute)

	}
}
