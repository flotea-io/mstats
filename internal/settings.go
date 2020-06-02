/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package internal

type Settings struct {
	ID    int    `json:"id" gorm:"AUTO_INCREMENT"`
	Name  string `json:"Name" gorm:"unique_index"`
	Value string `json:"Value"`
}
