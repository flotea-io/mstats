/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package internal

type PacketType int
type Status int

const (
	RegisterPacket          PacketType = 0
	SpeedChangePacket       PacketType = 1
	RegisterExportsPacket   PacketType = 2
	RunClaymorePacket       PacketType = 3
	StopClaymorePacket      PacketType = 4
	SendClaymoreStatPacket  PacketType = 5
	SendShellCommandPacket  PacketType = 6
	SendHardwareStatPacket  PacketType = 7
	ClientAlreadyRegistered PacketType = 8
	RequestWorkPacket       PacketType = 9
	UpdateClientPacket      PacketType = 10
	AlivePacket             PacketType = 11
	CardStatePacket         PacketType = 12
)
