/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package config

var license = "0dc42904d13d7288e1a3309323f3e4afb145ebea0020da5647439baaece76223b0359124a0a65068a09ba60bbd9f622d11a859fe56243e283234487b4b03ddca"
var licenseServerIP = ""
var licenseServerPort = ""

func GetLicenseAddress() string {
	return licenseServerIP + ":" + licenseServerPort
}
