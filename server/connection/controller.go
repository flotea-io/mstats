package connection

import (
	"encoding/json"
	"fmt"
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/server/client"
	"mstats-new/server/config"
	"mstats-new/util"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

var clients = make(map[string]net.Conn)
var clientsAliveTime = make(map[string]time.Time)
var aliveMutex = sync.RWMutex{}
var clientsMutex = sync.RWMutex{}

func SendMessageToClients(message string) {
	clientsMutex.RLock()
	for _, val := range clients {
		fmt.Fprintf(val, message)
	}
	clientsMutex.RUnlock()
}

func registerClient(conn net.Conn, packet internal.Packet) {
	logger.Info("I am registering new client " + packet.Sender)
	clientsMutex.Lock()
	clients[packet.Sender] = conn
	clientsMutex.Unlock()
	aliveMutex.Lock()
	clientsAliveTime[packet.Sender] = time.Now()
	aliveMutex.Unlock()
	logger.Info("New client " + packet.Sender + "(" + conn.RemoteAddr().String() + ") successfully registered")

}

func SendMessageToClient(name string, messages string) {
	clientsMutex.RLock()
	if _, found := clients[name]; !found {
		logger.Error("Error sending message (client doesn't exist)")
		return
	}
	fmt.Fprintf(clients[name], messages)
	clientsMutex.RUnlock()
}

func unregisterClient(ip string) {
	for name, val := range clients {
		if val.RemoteAddr().String() == ip {
			clientsMutex.Lock()
			delete(clients, name)
			clientsMutex.Unlock()
			aliveMutex.Lock()
			delete(clientsAliveTime, name)
			aliveMutex.Unlock()
			logger.Info("Unregistered client " + name)
			return
		}
	}
	logger.Error("Can't unregister unknown client with IP " + ip)
}

func isRegistered(name string) bool {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()
	if _, ok := clients[name]; ok {
		return true
	}
	return false
}

func receiveFromConnection(conn net.Conn) {
	for {

		message := make([]byte, 4096*4)
		length, err := conn.Read(message)

		if err != nil {
			unregisterClient(conn.RemoteAddr().String())
			conn.Close()
			break
		}

		if length == 0 {
			break
		}

		//solve problem with receiving two jsons once a time
		packets := strings.SplitAfter(string(message[:length]), "\"password\":\""+config.GetPassword()+"\"}")
		for _, msg := range packets {

			if len(msg) == 0 {
				break
			}

			var packet internal.Packet
			unmarshal := json.Unmarshal([]byte(msg[:len(msg)]), &packet)

			if unmarshal != nil {
				logger.Error("Something went wrong when receiving packet from " + conn.RemoteAddr().String())
				logger.Error(msg)
				logger.Error(unmarshal.Error())
				//unregisterClient(conn.RemoteAddr().String())
				//conn.Close()
				continue
			}

			if packet.Password != config.GetPassword() {
				logger.Warning("Client with IP " + conn.LocalAddr().String() + " trying to connect with bad password!")
				unregisterClient(conn.RemoteAddr().String())
				conn.Close()
				return
			}

			if !isRegistered(packet.Sender) && packet.PacketType != internal.PacketType(internal.RegisterPacket) {
				logger.Warning("Client with IP " + conn.LocalAddr().String() + " trying to send packet but he is not registered!")
				//unregisterClient(conn.RemoteAddr().String())
				conn.Close()
				return
			}

			if isRegistered(packet.Sender) && conn != clients[packet.Sender] {
				aliveMutex.Lock()
				timeFromLastAliveLog := time.Since(clientsAliveTime[packet.Sender])
				aliveMutex.Unlock()
				twoMinutes := 2 * time.Minute
				if timeFromLastAliveLog < twoMinutes {
					logger.Warning("Client with IP " + conn.RemoteAddr().String() + " connecting with used name: " + packet.Sender + " but old client is alive. Skipping..")
					unregister := internal.CreatePacket("Server", "", internal.ClientAlreadyRegistered, config.GetPassword())
					fmt.Fprintf(conn, unregister.ToJson())
					conn.Close()
					return
				}
				logger.Warning("Client with IP " + conn.RemoteAddr().String() + " connecting with used name: " + packet.Sender + " and old client is dead. Unregistering old one..")
				clientsMutex.RLock()
				clients[packet.Sender].Close()
				clientsMutex.RUnlock()
				unregisterClient(clients[packet.Sender].RemoteAddr().String())

			}

			switch packet.PacketType {
			case internal.RegisterPacket:
				registerClient(conn, packet)
				registerPacket := internal.CreatePacket("server", "", internal.RegisterPacket, config.GetPassword()).ToJson()
				SendMessageToClient(packet.Sender, registerPacket)
			case internal.SendClaymoreStatPacket:
				var stat internal.ClaymoreStat
				unmarshal := json.Unmarshal([]byte(packet.Message), &stat)

				if len(stat.MinerName) < 4 {
					continue
				}

				if unmarshal != nil {
					logger.Error("Something went wrong when receiving packet from " + conn.RemoteAddr().String())
					println(unmarshal.Error())
					conn.Close()
					break
				}

				db, err := gorm.Open("mysql", config.GetDatabaseDSN())

				if err != nil {
					db.Close()
					continue
				}

				stat.Time = util.CurrentTime()
				db.Create(&stat)
				db.Close()

				logger.Info("Successfully created new claymore log from " + packet.Sender)

			case internal.SendHardwareStatPacket:
				var stat internal.HardwareStat

				unmarshal := json.Unmarshal([]byte(packet.Message), &stat)

				if unmarshal != nil {
					logger.Error("Something went wrong when receiving packet from " + conn.RemoteAddr().String())
					logger.Error(unmarshal.Error())
					conn.Close()
					break
				}

				db, err := gorm.Open("mysql", config.GetDatabaseDSN())

				if err != nil {
					db.Close()
					continue
				}
				stat.Time = util.CurrentTime()
				db.Create(&stat)
				db.Close()

				logger.Info("Successfully created new hardware log from " + packet.Sender)
			case internal.RequestWorkPacket:
				db, err := gorm.Open("mysql", config.GetDatabaseDSN())

				defer db.Close()

				if err != nil {
					logger.Error("Problem with connect to database..")
					logger.Error(err.Error())
					os.Exit(1)
					return
				}
				var miner client.CustomClient
				db.Where("name = ?", packet.Sender).First(&miner)

				var cfgBasic internal.ClayMoreBasicConfig
				var cfgExport internal.ClayMoreBasicConfig
				var wallet internal.Wallet

				db.Where("name = ?", "export").First(&cfgExport)
				db.Where("name = ?", miner.Config).First(&cfgBasic)
				db.Where("id = ?", miner.WalletID).First(&wallet)

				if cfgBasic.Params == "" {
					db.Where("name = ?", "default").First(&cfgBasic)
				}
				if wallet.Address == "" {
					db.Where("is_default = ?", 1).First(&wallet)
				}

				cfgBasic.Params = strings.Replace(cfgBasic.Params, "$WALLET_ADDRESS$", wallet.Address, -1)

				SendMessageToClient(packet.Sender, internal.CreatePacket("server", cfgExport.ToJson(), internal.RegisterExportsPacket, config.GetPassword()).ToJson())
				//time.Sleep(500000000 * time.Nanosecond)
				SendMessageToClient(packet.Sender, internal.CreatePacket("server", cfgBasic.ToJson(), internal.RunClaymorePacket, config.GetPassword()).ToJson())

			case internal.AlivePacket:
				registerAliveClient(packet.Sender)
			default:
				logger.Warning("Client with IP " + conn.LocalAddr().String() + " send unauthorized packet!")
				unregisterClient(conn.RemoteAddr().String())
				conn.Close()
			}
		}

	}
}

func GetConnectedClients() map[string]internal.Miner {
	data := make(map[string]internal.Miner)
	for name, connection := range clients {
		host, port, err := net.SplitHostPort(connection.RemoteAddr().String())

		if err != nil {
			logger.Error("Something went wrong when splitting remote address to host and port")
			continue
		}

		data[name] = internal.Miner{MinerName: name, MinerIP: host, MinerPort: port}
	}
	return data
}

func RebootClient(name string) {
	packet := internal.CreatePacket("Server", "sudo shutdown -r now", internal.SendShellCommandPacket, config.GetPassword())
	SendMessageToClient(name, packet.ToJson())
}

func ShutdownClient(name string) {
	packet := internal.CreatePacket("Server", "sudo shutdown -h now", internal.SendShellCommandPacket, config.GetPassword())
	SendMessageToClient(name, packet.ToJson())
}

func UpdateClient(name string) {
	packet := internal.CreatePacket("Server", "", internal.UpdateClientPacket, config.GetPassword())
	SendMessageToClient(name, packet.ToJson())
}

func registerAliveClient(name string) {
	aliveMutex.Lock()
	clientsAliveTime[name] = time.Now()
	aliveMutex.Unlock()
}
