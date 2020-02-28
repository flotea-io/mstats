package files

import (
	"mstats-new/client/config"
	"mstats-new/logger"
	"mstats-new/util"
	"os"
	"os/exec"
)

func InitFilesCheck() {
	logger.Info("I am checking required files..")
	if _, err := os.Stat("scripts"); os.IsNotExist(err) {
		logger.Warning("Folder with required scripts not found. I am starting downloading..")

		fileUrl := config.ScriptsFileUrl
		err := util.DownloadFile(util.ProgramPath()+"/scripts.tar.gz", fileUrl)

		if err != nil {
			logger.Error("Failed to download required scripts..")
			logger.Error(err.Error())
			os.Exit(1)
			return
		}

		exec.Command("tar", "-zxvf", "scripts.tar.gz").Run()
		exec.Command("chmod", "-R", "777", "scripts").Run()
		exec.Command("rm", "scripts.tar.gz").Run()

		logger.Info("Successfully download required scripts..")
		return
	}

}
