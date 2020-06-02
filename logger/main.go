/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package logger

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

var infoLogPath = programPath() + "/logs/info.log"
var warningLogPath = programPath() + "/logs/warning.log"
var errorLogPath = programPath() + "/logs/error.log"
var skipped = false

func InitLoggerSystem() {
	exec.Command("mkdir", "logs").Run()

	var infoFile, err1 = os.Create(infoLogPath)
	var warningFile, err2 = os.Create(warningLogPath)
	var errorFile, err3 = os.Create(errorLogPath)

	if err1 != nil || err2 != nil || err3 != nil {
		skipped = true
		Error("Can't set up logs files, skipping saving to file..")
		return
	}

	infoLogger = log.New(infoFile, "", 0)
	warningLogger = log.New(warningFile, "", 0)
	errorLogger = log.New(errorFile, "", 0)

}

func Error(message string) {
	var msg = getTime() + " [ERROR] " + message
	color.Red(msg)
	if !skipped {
		errorLogger.Printf(msg)
	}

}

func Info(message string) {
	var msg = getTime() + " [INFO] " + message
	color.White(msg)
	if !skipped {
		infoLogger.Printf(msg)

	}
}

func Warning(message string) {
	var msg = getTime() + " [WARNING] " + message
	color.Yellow(msg)
	if !skipped {
		warningLogger.Printf(msg)
	}
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func programPath() string {
	ex, err := os.Executable()

	if err != nil {
		Error("Can't read program folder path")
		Error(err.Error())
		return ""
	}

	exPath := filepath.Dir(ex)
	return exPath
}
