package license

import (
	"mstats-new/internal"
	"mstats-new/server/config"
	"time"
)

var isLicenseValid bool

func InitLicense(checkPeriod time.Duration, name string) {
	isLicenseValid = false
	for {
		isLicenseValid = internal.CheckLicense(name, config.GetLicenseAddress())
		time.Sleep(checkPeriod)
	}
}

func IsValidated() bool {
	return isLicenseValid
}
