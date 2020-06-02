/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package connection

import (
	json2 "encoding/json"
	"mstats-new/client/config"
	"mstats-new/client/hardware"
	"mstats-new/client/programs"
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/util"
	"time"
)

func requestLogsTask() {
	//allow claymore to turn on, fixes empty records in database
	time.Sleep(60 * time.Second)
	for {
		if !(programs.CurrentWorking["claymore"]) {
			time.Sleep(5 * time.Second)
			logger.Warning("Claymore is offline, skipping..")
			continue
		}
		stat := programs.RequestClayMoreStat()

		//@todo code below should close claymore when it freezes, but it didn't work :(
		//if stat.TotalEthHashRate == 0 {
		//	programs.StopClaymore()
		//}

		json := stat.ToJson()
		sendPacketToServer(internal.CreatePacket(config.GetClientName(), json, internal.SendClaymoreStatPacket, config.GetPassword()))
		time.Sleep(30 * time.Second)
	}
}

func startHardwareLogTask() {
	for {
		if !connected {
			time.Sleep(30 * time.Second)
			continue
		}

		time.Sleep(30 * time.Second)

		SendHardwareStatPacket()

	}
}

func CheckIfMinerWorks() {
	time.Sleep(10 * time.Second)
	for {
		if !programs.IsCurrentlyWorking() {
			logger.Warning("Not working, requesting job from server..")
			RequestJob()
		}
		time.Sleep(30 * time.Second)
	}
}

func startIAmAliveLogging() {
	for {
		if !connected {
			time.Sleep(10 * time.Second)
			continue
		}
		time.Sleep(30 * time.Second)
		packet := internal.CreatePacket(config.GetClientName(), "", internal.AlivePacket, config.GetPassword())
		sendPacketToServer(packet)
	}
}

func SendHardwareStatPacket() {
	bytes, err := json2.Marshal(hardware.Cards)

	if err != nil {
		logger.Warning(err.Error())
		return
	}
	json := internal.HardwareStat{MinerName: config.GetClientName(), GpuList: string(bytes), Time: util.CurrentTime()}.ToJson()
	sendPacketToServer(internal.CreatePacket(config.GetClientName(), json, internal.SendHardwareStatPacket, config.GetPassword()))
}
