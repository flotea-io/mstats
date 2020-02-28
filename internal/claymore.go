package internal

import (
	"encoding/json"
	"mstats-new/logger"
)

type ClayMoreBasicConfig struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Params   string `json:"params"`
	Currency string `json:"currency"`
}

func (c *ClayMoreBasicConfig) ToJson() string {
	bytes, err := json.Marshal(&c)
	if err != nil {
		logger.Error("Something went wrong with parsing config to json..")
		return ""
	}
	return string(bytes)
}

type ClaymoreStat struct {
	MinerName                 string `json:"MinerName"`
	MinerVersion              string `json:"MinerVersion"`
	RunningTime               int    `json:"RunningTime"`
	TotalEthHashRate          int    `json:"TotalEthHashRate"`
	EthShares                 int    `json:"EthShares"`
	EthRejectedShares         int    `json:"EthRejectedShares"`
	DetailedEthHashRatePerGPU string `json:"DetailedEthHashRatePerGPU"`
	TotalDcrHashRate          int    `json:"TotalDcrHashRate"`
	DcrShares                 int    `json:"DcrShares"`
	DcrRejectedShares         int    `json:"DcrRejectedShares"`
	DetailedDcrHashRatePerGPU string `json:"DetailedDcrHashRatePerGPU"`
	Temperatures              string `json:"Temperatures"`
	MiningPool                string `json:"MiningPool"`
	EthInvalidShares          int    `json:"EthInvalidShares"`
	EthPoolSwitches           int    `json:"EthPoolSwitches"`
	DcrInvalidShares          int    `json:"DcrInvalidShares"`
	DcrPoolSwitches           int    `json:"DcrPoolSwitches"`
	Time                      string `json:"Time"`
}

func (c *ClaymoreStat) ToJson() string {
	bytes, err := json.Marshal(&c)

	if err != nil {
		logger.Error("Something went wrong with parsing stat to json..")
		return ""
	}
	return string(bytes)
}
