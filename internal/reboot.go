/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package internal

type Reboot struct {
	ID        int    `json:"id" gorm:"AUTO_INCREMENT"`
	MinerName string `json:"MinerName"`
	Reason    string `json:"reason"`
	Time      string `json:"time"`
}

type Restart struct {
	ID        int    `json:"id" gorm:"AUTO_INCREMENT"`
	MinerName string `json:"MinerName"`
	Reason    string `json:"reason"`
	Config    string `json:"cfg"`
	Time      string `json:"time"`
}

type Shutdown struct {
	ID        int    `json:"id" gorm:"AUTO_INCREMENT"`
	MinerName string `json:"MinerName"`
	Reason    string `json:"reason"`
	Time      string `json:"time"`
}
