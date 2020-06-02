/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package client

type CustomClient struct {
	Id       int    `json:"id"`
	Name     string `json:"Name"`
	Config   string `json:"Config" gorm:"default:\"autostart\""`
	WalletID int    `json:"WalletID" sql:"type:int"`
}
