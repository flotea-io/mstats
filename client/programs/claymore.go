/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package programs

import (
	"bufio"
	"fmt"
	"mstats-new/client/config"
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/util"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

var ClaymoreTempPassword = "abcabc"
var clientIP = "000.000.000.000"

var statsRequest = "{\"id\":0,\"jsonrpc\":\"2.0\",\"method\":\"miner_getstat1\",\"psw\":\"pass\"}"

func ClayMoreDownloadIfNotExist() {
	logger.Info("I am checking file exist..")

	out, err := exec.Command("ls").Output()

	if err != nil {
		logger.Error("Something went wrong with checking files")
		logger.Error(err.Error())
		return
	}

	stringOut := string(out)

	if strings.Contains(stringOut, "claymore") {
		logger.Info("Claymore is installed on that machine")
		return
	}

	fileUrl := "http://" + clientIP + ":9920/files/claymore.tar.gz"
	err1 := util.DownloadFile(util.ProgramPath()+"/claymore.tar.gz", fileUrl)

	if err1 != nil {
		logger.Error("Failed to download required files..")
		logger.Error(err.Error())
		os.Exit(1)
		return
	}

	exec.Command("tar", "-zxvf", "claymore.tar.gz").Run()
	exec.Command("chmod", "-R", "777", "claymore").Run()
	exec.Command("rm", "claymore.tar.gz").Run()

	logger.Info("Successfully installed claymore..")
}

func RequestClayMoreStat() *internal.ClaymoreStat {
	statistics := requestStatistics()

	if !(CurrentWorking["claymore"]) {
		logger.Warning("Claymore is offline")
		return &internal.ClaymoreStat{}
	}

	if statistics == "null" {
		logger.Error("Claymore frozen ")
		//reboot here
	}

	if len(statistics) < 10 {
		logger.Warning("Claymore not loaded yet..")
		return &internal.ClaymoreStat{}
	}

	value := gjson.Get(statistics, "result")

	resultArray := value.Array()
	var version = resultArray[0].String()
	var time, _ = strconv.Atoi(resultArray[1].String())

	/**
	ETH SECTION
	*/
	ethSharesArray := strings.Split(resultArray[2].String(), ";")
	var totalEthShares, _ = strconv.Atoi(ethSharesArray[0])
	var ethShares, _ = strconv.Atoi(ethSharesArray[1])
	var ethRejectedShares, _ = strconv.Atoi(ethSharesArray[2])
	var detailedEthGPU = resultArray[3].String()

	/**
	DCR SECTION
	*/

	dcrSharesArray := strings.Split(resultArray[4].String(), ";")
	var totalDcrShares, _ = strconv.Atoi(dcrSharesArray[0])
	var dcrShares, _ = strconv.Atoi(dcrSharesArray[1])
	var dcrRejectedShares, _ = strconv.Atoi(dcrSharesArray[2])
	var detailedDcrGPU = resultArray[5].String()

	var temp = resultArray[6].String()
	var pool = resultArray[7].String()

	/**
	Other SECTION
	*/

	otherShares := strings.Split(resultArray[8].String(), ";")

	//192.168.88.159
	var ethInvalidShares, _ = strconv.Atoi(otherShares[0])
	var ethPoolSwitches, _ = strconv.Atoi(otherShares[1])
	var dcrInvalidShares, _ = strconv.Atoi(otherShares[2])
	var dcrPoolSwitches, _ = strconv.Atoi(otherShares[3])

	var stat = internal.ClaymoreStat{MinerName: config.GetClientName(),
		MinerVersion:              version,
		RunningTime:               time,
		EthShares:                 ethShares,
		EthRejectedShares:         ethRejectedShares,
		DetailedEthHashRatePerGPU: detailedEthGPU,
		DcrShares:                 dcrShares,
		DcrRejectedShares:         dcrRejectedShares,
		DetailedDcrHashRatePerGPU: detailedDcrGPU,
		Temperatures:              temp,
		MiningPool:                pool,
		EthInvalidShares:          ethInvalidShares,
		EthPoolSwitches:           ethPoolSwitches,
		DcrInvalidShares:          dcrInvalidShares,
		DcrPoolSwitches:           dcrPoolSwitches,
		TotalEthHashRate:          totalEthShares,
		TotalDcrHashRate:          totalDcrShares,
		Time:                      util.CurrentTime()}

	return &stat
}

func requestStatistics() string {

	//@todo check
	var ip = "0.0.0.0:3333"
	conn, err := net.Dial("tcp", ip)

	if err != nil {
		logger.Error("Can't connect to claymore..")
		logger.Error(err.Error())
		return "null"
	}

	fmt.Fprintf(conn, strings.Replace(statsRequest, "pass", ClaymoreTempPassword, -1)+"\n")

	message, err := bufio.NewReader(conn).ReadString('\n')

	return message
}

func StopClaymore() {
	run := exec.Command("screen", "-S", "claymore", "-X", "quit").Run()

	if run != nil {
		logger.Error("Problem with close screen with claymore..")
		logger.Error(run.Error())
	}
	CurrentWorking["claymore"] = false
	logger.Info("Successfully stop claymore..")
}

func GetClaymoreRebootScriptIfNotExist() {
	if _, err := os.Stat("claymore/reboot.sh"); os.IsNotExist(err) {
		logger.Warning("Script to reboot claymore not found. Downloading..")

		fileUrl := config.ClaymoreRebootFileUrl
		err := util.DownloadFile(util.ProgramPath()+"/claymore/reboot.sh", fileUrl)
		exec.Command("chmod", "+x", util.ProgramPath()+"/claymore/reboot.sh").Run()

		if err != nil {
			logger.Error("Failed to download script..")
			logger.Error(err.Error())
			os.Exit(1)
			return
		}
		logger.Info("Successfully download reboot script..")
		return
	}
}
