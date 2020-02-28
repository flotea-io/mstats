package connection

import (
	"mstats-new/logger"
	"mstats-new/server/config"
	"net"
	"os"
	"time"
)

func InitSocketServer() {
	logger.Info("I am starting socket server..")

	listener, error := net.Listen("tcp", config.GetStringAddress())

	if error != nil {
		logger.Error("Socket server start failed!")
		logger.Error(error.Error())
		os.Exit(0)
		return
	}

	logger.Info("Server listening on " + config.GetStringAddress())

	go removeDeadClientsDaemon()

	for {
		connection, error := listener.Accept()

		if error != nil {
			logger.Error("Can't accept connection..")
			logger.Error(error.Error())
			continue
		}

		go receiveFromConnection(connection)
	}
}

func removeDeadClientsDaemon() {
	twoMinutes := 2 * time.Minute
	for {
		aliveMutex.RLock()
		myClientsAliveTime := clientsAliveTime
		aliveMutex.RUnlock()
		clientsMutex.RLock()
		myClients := clients
		clientsMutex.RUnlock()
		for client, aliveTime := range myClientsAliveTime {
			if time.Since(aliveTime) > twoMinutes {
				logger.Warning("Client with name " + client + " is not responding since: " + aliveTime.Format("2006-01-02 15:04:05") + ". Unregistering...")
				unregisterClient(myClients[client].RemoteAddr().String())
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
