/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package programs

import (
	"mstats-new/internal"
	"mstats-new/server/config"
	"mstats-new/server/connection"
)

func SetFanSpeed(machine string, cardId string, speed string) {

	connection.SendMessageToClient(machine, internal.CreatePacket("server", cardId+"|"+speed, internal.SpeedChangePacket, config.GetPassword()).ToJson())

}

func CardState(machine string, id string, state string) {
	packet := internal.CreatePacket("Server", id+"|"+state, internal.CardStatePacket, config.GetPassword())
	connection.SendMessageToClient(machine, packet.ToJson())
}
