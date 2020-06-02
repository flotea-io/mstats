/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package internal

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"mstats-new/logger"
	"mstats-new/util"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	filename = "license.json"
)

func generateMessage(text string) string {
	//return "TEST DATA TO COMPUTE"
	return "Nonstromo" + text + " " + util.CurrentTime()
}

func CheckLicense(text string, address string) bool {

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Can't open file with key and license (" + filename + ")")
		logger.Error(err.Error())
	}

	license := struct {
		Key        string `json:"key"`
		License_id int    `json:"license_id"`
	}{}

	err = json.Unmarshal(fileData, &license)
	if err != nil {
		logger.Error("Can't parse file with key and license (" + filename + ")")
		logger.Error(err.Error())
	}
	PEMKey := license.Key
	license_id := strconv.Itoa(license.License_id)

	message := generateMessage(text)
	resp, err := http.PostForm("http://"+address+"/api/check-license", url.Values{"license_id": {license_id}, "message": {message}})

	if err != nil {
		logger.Error("Can't connect to license server..")
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	sig, err := ioutil.ReadAll(resp.Body)
	signature, err := base64.StdEncoding.DecodeString(string(sig))

	if err != nil {
		logger.Error("Can't read body from license server..")
		logger.Error(err.Error())
		os.Exit(1)
	}

	PEMBlock, _ := pem.Decode([]byte(PEMKey))

	if PEMBlock == nil {
		logger.Error("Can't parse Public Key PEM..")
		return false
	}
	if PEMBlock.Type != "PUBLIC KEY" {
		return false
		logger.Error("Found wrong key type")
	}

	pubkey, err := x509.ParsePKIXPublicKey(PEMBlock.Bytes)

	if err != nil {
		logger.Error("Can't read body of Public Key PEM..")
		logger.Error(err.Error())
		return false
	}

	// compute the sha1
	h := sha1.New()
	h.Write([]byte(message))

	// Verify
	err = rsa.VerifyPKCS1v15(pubkey.(*rsa.PublicKey), crypto.SHA1, h.Sum(nil), signature)

	if err != nil {
		logger.Info("Invalid license..")
		return false
	}

	logger.Info("Valid license..")
	return true
}
