/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package internal

import (
	"encoding/json"
	"mstats-new/logger"
)

type handlerError struct {
	Status int    `json:"status"`
	Reason string `json:"reason"`
}

func (c handlerError) ToJson() string {
	bytes, err := json.Marshal(&c)
	if err != nil {
		logger.Error("Something went wrong with parsing stat to json..")
		return ""
	}
	return string(bytes)
}

func CreateHandlerError(status int, reason string) handlerError {
	return handlerError{Status: status, Reason: reason}
}
