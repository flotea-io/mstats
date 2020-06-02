/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package database

import (
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/server/client"
	"mstats-new/server/config"
	"mstats-new/server/mail"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDatabase() {
	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		logger.Error("Problem with connect to database..")
		logger.Error(err.Error())
		os.Exit(1)
		return
	}

	logger.Info("I am checking tables..")

	if !db.HasTable(&internal.ClayMoreBasicConfig{}) {
		logger.Info("Table with claymore config not found.. Creating..")
		db.CreateTable(&internal.ClayMoreBasicConfig{})
		basicConfig := internal.ClayMoreBasicConfig{ID: 0, Name: "default", Params: "-epool eu1.ethermine.org:4444", Currency: "ETH"}
		basicETCConfig := internal.ClayMoreBasicConfig{ID: 1, Name: "default_etc", Params: "-epool eu1-etc.ethermine.org:4444", Currency: "ETH"}
		exportConfig := internal.ClayMoreBasicConfig{ID: 2, Name: "export", Params: "export GPU_MAX_HEAP_SIZE=100|export GPU_USE_SYNC_OBJECTS=1|export GPU_MAX_ALLOC_PERCENT=100|export GPU_SINGLE_ALLOC_PERCENT=100", Currency: ""}
		db.NewRecord(basicConfig)
		db.Create(&basicConfig)
		db.NewRecord(basicETCConfig)
		db.Create(&basicETCConfig)
		db.NewRecord(exportConfig)
		db.Create(&exportConfig)
		logger.Info("Successfully created default claymore config..")
	}

	if !db.HasTable(&internal.ClaymoreStat{}) {
		logger.Info("Table with claymore stat not found.. Creating..")
		db.CreateTable(&internal.ClaymoreStat{})
	}

	if !db.HasTable(&internal.HardwareStat{}) {
		logger.Info("Table with hardware stat not found.. Creating..")
		db.CreateTable(&internal.HardwareStat{})
	}

	if !db.HasTable(&client.CustomClient{}) {
		logger.Info("Table with hardware custom client not found.. Creating..")
		db.CreateTable(&client.CustomClient{})
	}

	if !db.HasTable(&mail.MailGunConfig{}) {
		logger.Info("Table with mailgun config not found.. Creating..")
		db.CreateTable(&mail.MailGunConfig{})
	}

	if !db.HasTable(&mail.Recipient{}) {
		logger.Info("Table with recipients not found.. Creating..")
		db.CreateTable(&mail.Recipient{})
	}

	if !db.HasTable(&mail.EmailHistory{}) {
		logger.Info("Table with email history not found.. Creating")
		db.CreateTable(&mail.EmailHistory{})
	}

	if !db.HasTable(&internal.HardReset{}) {
		logger.Info("Table with machines to hard reset not found.. Creating")
		db.CreateTable(&internal.HardReset{})
	}

	if !db.HasTable(&internal.TempInfo{}) {
		logger.Info("Table with temperature info not found.. Creating")
		db.CreateTable(&internal.TempInfo{})
	}

	if !db.HasTable(&internal.Reboot{}) {
		logger.Info("Table with reboot info not found.. Creating")
		db.CreateTable(&internal.Reboot{})
	}

	if !db.HasTable(&internal.Restart{}) {
		logger.Info("Table with restart info not found.. Creating")
		db.CreateTable(&internal.Restart{})
	}

	if !db.HasTable(&internal.Shutdown{}) {
		logger.Info("Table with shutdown info not found.. Creating")
		db.CreateTable(&internal.Shutdown{})
	}

	if !db.HasTable(&internal.PinFunction{}) {
		logger.Info("Table with arduino pins function not found.. Creating")
		db.CreateTable(&internal.PinFunction{})
	}

	if !db.HasTable(&internal.HardShutdown{}) {
		logger.Info("Table with machines to hard shutdown not found.. Creating")
		db.CreateTable(&internal.HardShutdown{})
	}

	if !db.HasTable(&internal.Settings{}) {
		logger.Info("Table with settings not found.. Creating")
		db.CreateTable(&internal.Settings{})
	}

	if !db.HasTable(&internal.Wallet{}) {
		logger.Info("Table with Wallet not found.. Creating")
		db.CreateTable(&internal.Wallet{})
	}

	db.AutoMigrate(&internal.PinFunction{})
	db.AutoMigrate(&internal.Wallet{})
	db.AutoMigrate(&client.CustomClient{})
	db.AutoMigrate(&internal.ClayMoreBasicConfig{})

	logger.Info("Successfully connected to database..")
}

func CleanDatabase() {
	timeToDelete := time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02 15:04:05")

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		logger.Error("Problem with connect to database..")
		logger.Error(err.Error())
		os.Exit(1)
		return
	}

	db.Where("time < ?", timeToDelete).Delete(internal.ClaymoreStat{})
	db.Where("time < ?", timeToDelete).Delete(internal.HardwareStat{})
}
