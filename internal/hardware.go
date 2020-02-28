package internal

import (
	"encoding/json"
	"mstats-new/logger"
)

type HardwareStat struct {
	MinerName string `json:"minerName"`
	GpuList   string `json:"gpuList" sql:"size:1024;"`
	Time      string `json:"time"`
}

func (c HardwareStat) ToJson() string {
	bytes, err := json.Marshal(&c)
	if err != nil {
		logger.Error("Something went wrong with parsing stat to json..")
		return ""
	}
	return string(bytes)
}
