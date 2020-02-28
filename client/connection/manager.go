package connection

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mstats-new/client/config"
	"mstats-new/client/hardware"
	"mstats-new/client/license"
	"mstats-new/client/programs"
	"mstats-new/client/updater"
	"mstats-new/internal"
	"mstats-new/logger"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var client net.Conn
var connected = false
var logTask = false
var aliveTask = false
var registerAtempt = 0

func handleMessagesFromServer(conn net.Conn) {
	for {
		message := make([]byte, 4096)
		length, err := conn.Read(message)

		if err != nil {
			logger.Error("Something went wrong when reading packet..")
			logger.Error(err.Error())
			connected = false
			break
		}

		if length <= 0 {
			continue
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
				logger.Error(unmarshal.Error())
				continue
			}

			if packet.Password != config.GetPassword() {
				logger.Warning("client with IP " + conn.LocalAddr().String() + " trying to send packet with bad password!")
				continue
			}

			switch packet.PacketType {
			case internal.RegisterExportsPacket:
				logger.Info("I am registering exports program..")

				var cfg internal.ClayMoreBasicConfig
				unmarshal := json.Unmarshal([]byte(packet.Message), &cfg)

				if unmarshal != nil {
					logger.Error("Something went wrong when receiving packet from " + conn.RemoteAddr().String())
					logger.Error(unmarshal.Error())
					connected = false
					conn.Close()
					break
				}

				split := strings.Split(cfg.Params, "|")
				for _, val := range split {

					str := "bash -c " + val
					args := strings.Fields(str)

					run := exec.Command(args[0], args[1:]...).Run()

					if run != nil {
						logger.Error("Can't run export command..")
						logger.Error(run.Error())
						continue
					}
				}

				logger.Info("Successfully registering exports...")
			case internal.RunClaymorePacket:

				if !license.IsValidated() {
					logger.Info("License is not valid, can't start mining")
					continue
				}

				if programs.CurrentWorking["claymore"] == true {
					logger.Info("Claymore already running.. Can't run second time")
					continue
				}

				logger.Info("I am starting claymore program..")
				programs.ClayMoreDownloadIfNotExist()
				programs.GetClaymoreRebootScriptIfNotExist()

				var cfg internal.ClayMoreBasicConfig
				unmarshal := json.Unmarshal([]byte(packet.Message), &cfg)

				if unmarshal != nil {
					logger.Error("Something went wrong when parsing json packet from " + conn.RemoteAddr().String())
					logger.Error(unmarshal.Error())
					continue
				}

				h := sha512.New()
				h.Write([]byte(time.Now().String()))
				pass := hex.EncodeToString(h.Sum(nil))
				passw := pass[0:20]
				programs.ClaymoreTempPassword = passw

				var basic = "screen -S claymore -d -m claymore/ethdcrminer64 " + cfg.Params + " -mpsw " + passw

				basic = strings.Replace(basic, "$WALLET_IDENTIFY$", config.GetClientName(), -1)

				args := strings.Fields(basic)

				run := exec.Command(args[0], args[1:]...).Run()

				if run != nil {
					logger.Error("Can't run screen with claymore..")
					logger.Error(run.Error())
					continue
				}

				programs.CurrentWorking["claymore"] = true
				go requestLogsTask()

				logger.Info("Successfully start claymore..")

			case internal.StopClaymorePacket:
				programs.StopClaymore()
			case internal.SendShellCommandPacket:
				command := packet.Message
				args := strings.Fields(command)

				run := exec.Command(args[0], args[1:]...).Run()

				if run != nil {
					logger.Error("Some problem with run command..")
					logger.Error(run.Error())
					continue
				}

				logger.Info("Command (" + command + ") executed successfully")

			case internal.SpeedChangePacket:
				split := strings.Split(packet.Message, "|")
				//speed 0 toggles automatic control
				hardware.SetFanSpeed(split[0], split[1])
				SendHardwareStatPacket()
			case internal.ClientAlreadyRegistered:
				connected = false
				registerAtempt++
				if registerAtempt > 5 {
					logger.Error("Client with this name exist, choose another")
					conn.Close()
					os.Exit(1)
				}

				logger.Error("Client with this name is connected, trying to reconnect.. Attempt " + strconv.Itoa(registerAtempt) + "/5")
				time.Sleep(1 * time.Minute)
				packet := internal.CreatePacket(config.GetClientName(), "", internal.RegisterPacket, config.GetPassword())
				sendPacketToServer(packet)
			case internal.RegisterPacket:
				connected = true
				registerAtempt = 0

				if !logTask {
					logTask = true
					go startHardwareLogTask()
				}

				if !aliveTask {
					aliveTask = true
					go startIAmAliveLogging()
				}
			case internal.UpdateClientPacket:
				updater.Init()
			case internal.CardStatePacket:
				programs.StopClaymore()

				split := strings.Split(packet.Message, "|")
				var state string
				if split[1] == "0" {
					state = "0"
				} else {
					state = "1"
				}
				hardware.SetCardState(split[0], state)
				RequestJob()
			}
		}
	}
}

func InitClient() {
	go forceConnect()
}

func forceConnect() {
	for {

		if connected {
			time.Sleep(30 * time.Second)
			continue
		}

		connection, err := net.Dial("tcp", config.GetStringAddress())

		if err != nil {
			logger.Error("Something went wrong when trying to connect..")
			logger.Error(err.Error())
			time.Sleep(30 * time.Second)
			continue
		}

		client = connection

		go handleMessagesFromServer(connection)

		packet := internal.CreatePacket(config.GetClientName(), "", internal.RegisterPacket, config.GetPassword())

		sendPacketToServer(packet)

		time.Sleep(10 * time.Second)
	}

}

func sendPacketToServer(packet *internal.Packet) {
	fmt.Fprintf(client, packet.ToJson())
	logger.Warning("I am sending packet.. " + packet.ToJson())
}

func RequestJob() {
	packet := internal.CreatePacket(config.GetClientName(), "", internal.RequestWorkPacket, config.GetPassword())
	sendPacketToServer(packet)
}
