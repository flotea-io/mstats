/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package dependencies

import (
	"mstats-new/logger"
	"os"
	"os/exec"
	"strings"
)

var dependencies = []string{"lm-sensors", "screen"}

func InitDependenciesCheck() {
	logger.Info("I am checking for required dependencies..")

	for _, val := range dependencies {
		logger.Info("I am checking " + val + " dependency")
		out, err := exec.Command("apt-cache", "policy", val).Output()

		if err != nil {
			logger.Error("Something went wrong with checking dependency " + val)
			logger.Error(err.Error())
			os.Exit(1)
			return
		}

		stringOut := string(out)

		if strings.Contains(stringOut, "(none)") {
			logger.Info("I am installing required dependency " + val)
			exec.Command("sudo", "apt-get", "install", val).Run()
		} else {
			logger.Info("Dependency found, skipping..")
		}

	}
}
