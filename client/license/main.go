package license

import (
	"mstats-new/client/config"
	"mstats-new/client/programs"
	"mstats-new/internal"
	"mstats-new/logger"
	"os/exec"
	"time"
)

var isLicenseValid bool

func InitLicense(checkPeriod time.Duration) {
	isLicenseValid = false
	for {
		isLicenseValid = internal.CheckLicense(" "+config.GetClientName(), config.GetLicenseAddress())

		if !IsValidated() {
			logger.Warning("License is not valid, mining stopped..")
			run := exec.Command("screen", "-S", "claymore", "-X", "quit").Run()

			if run != nil {
				logger.Error("Problem with close screen with claymore..")
				logger.Error(run.Error())
				continue
			}
			programs.CurrentWorking["claymore"] = false

		}
		time.Sleep(checkPeriod)
	}
}

func IsValidated() bool {
	return isLicenseValid
}
