/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package updater

import (
	"encoding/json"
	"io/ioutil"
	"mstats-new/client/config"
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/util"
	"net/http"
	"os"
	"os/exec"
)

func Init() {
	logger.Info("I am checking about updates..")
	url := config.UpdaterDetectUrl

	resp, err := http.Get(url)

	if err != nil {
		logger.Error("Can't connect to update server..")
		logger.Error(err.Error())
		return
	}

	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logger.Error("Can't read information from update server..")
		logger.Error(err.Error())
		return
	}

	var update internal.Update

	err2 := json.Unmarshal(html, &update)

	if err2 != nil {
		logger.Error("Can't parse data from update server to object..")
		logger.Error(err2.Error())
		return
	}

	if config.IsCurrentVersion(update.Client) {
		logger.Info("You got newest version..")
		return
	}

	logger.Warning("=========================")
	logger.Warning("I found new version!")
	logger.Warning("Type: " + update.Client)
	logger.Warning("I am starting updating..")
	logger.Warning("=========================")

	fileUrl := config.ClientFileUrl

	exec.Command("rm", "client").Run()
	err3 := util.DownloadFile(util.ProgramPath()+"/client", fileUrl)
	exec.Command("chmod", "777", "client").Run()

	if err3 != nil {
		panic(err3.Error())
	}

	logger.Info("Files updated successfully!")
	exec.Command("reboot", "-f").Run()
	os.Exit(1)
}
