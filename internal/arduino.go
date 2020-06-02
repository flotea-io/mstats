/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package internal

type HardReset struct {
	MachineNumber string `json:"machine"`
}

type HardShutdown struct {
	MachineNumber string `json:"machine"`
	Function      string `json:"func" sql:"type:int(1) unsigned"`
}

type TempInfo struct {
	ID   int    `json:"id" gorm:"AUTO_INCREMENT"`
	Data string `json:"data"`
	Time string `json:"time"`
}

type PinFunction struct {
	ID        *string `json:"ID" sql:"type:int PRIMARY KEY"`
	MinerName *string `json:"MinerName"`
	Function  *string `json:"Function" sql:"type:int(1) unsigned"`
	AutoReset *string `json:"AutoReset" sql:"type:int(1) unsigned" gorm:"default:0"`
}
