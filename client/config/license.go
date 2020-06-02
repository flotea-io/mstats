/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package config

var licenseServerIP = ""
var licenseServerPort = ""

func GetLicenseAddress() string {
	return licenseServerIP + ":" + licenseServerPort
}
