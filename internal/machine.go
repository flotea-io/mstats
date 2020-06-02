/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package internal

type Miner struct {
	MinerName string `json:"MinerName"`
	MinerIP   string `json:"MinerIP"`
	MinerPort string `json:"MinerPort"`
}

type Card struct {
	CardName         string `json:"-"`
	CardID           string `json:"id"`
	MonitorName      string `json:"-"`
	ManualFanControl bool   `json:"manualFanControl"`
	Temperature      int    `json:"temp"`
	MaxValue         int    `json:"-"`
	CurrentValue     int    `json:"-"`
	CurrentRPM       int    `json:"fanSpeed"`
	CurrentPercent   int    `json:"percentSpeed"`
	DeclaredPercent  int    `json:"declaredFanSpeed"`
	Enabled          string `json:"enabled"`
}

type CardDisable struct {
	Miner  string `json:"miner"`
	CardID string `json:"id"`
	State  string `json:"state"`
}
