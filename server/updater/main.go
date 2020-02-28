package updater

import (
	"encoding/json"
	"io/ioutil"
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/server/config"
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

	if config.IsCurrentVersion(update.Server) {
		logger.Info("You got newest version..")
		return
	}

	logger.Warning("=========================")
	logger.Warning("I found new version!")
	logger.Warning("Type: " + update.Server)
	logger.Warning("I am starting updating..")
	logger.Warning("=========================")

	fileUrl := config.ServerFileUrl

	exec.Command("rm", "server").Run()
	err3 := util.DownloadFile(util.ProgramPath()+"/server", fileUrl)
	exec.Command("chmod", "777", "server").Run()

	if err3 != nil {
		panic(err3.Error())
	}

	logger.Info("Files updated successfully!")
	exec.Command("reboot", "-f").Run()
	os.Exit(1)
}
