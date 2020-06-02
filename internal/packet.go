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

type Packet struct {
	Sender     string     `json:"sender"`
	Message    string     `json:"message"`
	PacketType PacketType `json:"packetType"`
	Password   string     `json:"password"`
}

func CreatePacket(sender string, message string, packetType PacketType, password string) *Packet {
	packet := &Packet{}
	packet.Sender = sender
	packet.Message = message
	packet.PacketType = packetType
	packet.Password = password
	return packet
}

func (packet *Packet) ToJson() string {
	bytes, err := json.Marshal(&packet)
	if err != nil {
		logger.Error("Something went wrong with parsing packet to json..")
		return ""
	}
	return string(bytes)
}

func FromJson(data string) Packet {
	var receivedPacket Packet

	json.Unmarshal([]byte(data), &receivedPacket)

	return receivedPacket
}
