package web

import (
	"encoding/json"
	"fmt"
	"log"
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/server/client"
	"mstats-new/server/config"
	"mstats-new/server/connection"
	"mstats-new/server/license"
	"mstats-new/server/mail"
	"mstats-new/server/programs"
	"mstats-new/util"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

var api string

func initRouters() *mux.Router {
	log.Println("Initing routers..")
	r := mux.NewRouter()

	r.HandleFunc("/", basicHandler).Methods(http.MethodGet)
	r.HandleFunc("/claymore/manager/{name}/{cfg}/{option}", claymoreHandler).Methods(http.MethodGet)

	r.HandleFunc("/claymore/log/{name}/limit/{p1}/{p2}", showClaymoreStatLogsWithLimitHandler).Methods(http.MethodGet)
	r.HandleFunc("/claymore/log/{name}/date/{p1}/{p2}", showClaymoreStatLogsWithDateLimitHandler).Methods(http.MethodGet)
	r.HandleFunc("/claymore/log/{name}", showClaymoreStatLogsHandler).Methods(http.MethodGet)

	r.HandleFunc("/hardware/log/{name}/limit/{p1}/{p2}", showHardwareStatLogsWithLimitHandler).Methods(http.MethodGet)
	r.HandleFunc("/hardware/log/{name}/date/{p1}/{p2}", showHardwareStatLogsWithDateLimitHandler).Methods(http.MethodGet)
	r.HandleFunc("/hardware/log/{name}", showHardwareStatLogsHandler).Methods(http.MethodGet)

	r.HandleFunc("/claymore/config/basic", createBasicConfigHandler).Methods(http.MethodPost)
	r.HandleFunc("/claymore/config/basic", getClaymoreBasicConfigsHandler).Methods(http.MethodGet)
	r.HandleFunc("/claymore/config/basic/delete", deleteBasicConfigHandler).Methods(http.MethodPost)
	//r.HandleFunc("/claymore/config/basic/{name}", createBasicConfigHandler).Methods(http.MethodGet)
	//r.HandleFunc("/claymore/config/basic/{name}", getClaymoreBasicConfigHandler).Methods(http.MethodGet)
	//r.HandleFunc("/claymore/config/basic/{name}", deleteBasicConfigHandler).Methods(http.MethodDelete)

	r.HandleFunc("/client", getCustomClientsHandler).Methods(http.MethodGet)
	r.HandleFunc("/client", createCustomClientHandler).Methods(http.MethodPost)
	r.HandleFunc("/client/updateapp", updateClientAppHandler).Methods(http.MethodPost)
	r.HandleFunc("/client/delete", deleteCustomClientHandler).Methods(http.MethodPost)

	r.HandleFunc("/clients", getClientsHandler).Methods(http.MethodGet)
	r.HandleFunc("/clients/disconected", getDisconectedClientsHandler).Methods(http.MethodGet)

	// zmiana względem poprzedniego api - post z jsonem {minerName: , reason:}
	r.HandleFunc("/reboot", rebootClientHandler).Methods(http.MethodPost)
	r.HandleFunc("/reboot/log/{name}", getRebootInfoHandler).Methods(http.MethodGet)
	r.HandleFunc("/reboot/log/{name}/limit/{p1}/{p2}", getRebootInfoWithLimitHandler).Methods(http.MethodGet)
	r.HandleFunc("/reboot/log/{name}/date/{p1}/{p2}", getRebootInfoWithDateLimitHandler).Methods(http.MethodGet)

	// zmiana względem poprzedniego api - post z jsonem {minerName: , reason:, cfg: }
	r.HandleFunc("/restart", restartAppHandler).Methods(http.MethodPost)
	r.HandleFunc("/restart/log/{name}", getRestartInfoHandler).Methods(http.MethodGet)
	r.HandleFunc("/restart/log/{name}/limit/{p1}/{p2}", getRestartInfoWithLimitHandler).Methods(http.MethodGet)
	r.HandleFunc("/restart/log/{name}/date/{p1}/{p2}", getRestartInfoWithDateLimitHandler).Methods(http.MethodGet)

	// json {minerName: , reason:}
	r.HandleFunc("/shutdown", shutdownClientHandler).Methods(http.MethodPost)
	r.HandleFunc("/shutdown/log/{name}", getShutdownInfoHandler).Methods(http.MethodGet)
	r.HandleFunc("/shutdown/log/{name}/limit/{p1}/{p2}", getShutdownInfoWithLimitHandler).Methods(http.MethodGet)
	r.HandleFunc("/shutdown/log/{name}/date/{p1}/{p2}", getShutdownInfoWithDateLimitHandler).Methods(http.MethodGet)

	r.HandleFunc("/latest/stat", getLatestStatHandler).Methods(http.MethodGet)
	r.HandleFunc("/latest/reboot", getRebootLastInfoHandler).Methods(http.MethodGet)
	r.HandleFunc("/latest/restart", getRestartLastInfoHandler).Methods(http.MethodGet)
	r.HandleFunc("/latest/temperatures", getTemperatureLastInfoHandler).Methods(http.MethodGet)

	r.HandleFunc("/mail", createRecipientHandler).Methods(http.MethodPost)
	r.HandleFunc("/mail", deleteRecipientHandler).Methods(http.MethodDelete)

	r.HandleFunc("/arduino/temps", arduinoTemperatureHandler).Methods(http.MethodGet)
	r.HandleFunc("/arduino/resets", arduinoResetHandler).Methods(http.MethodGet)
	r.HandleFunc("/arduino/pins", arduinoAddPinsHandler).Methods(http.MethodPost)
	r.HandleFunc("/arduino/pins/clear", arduinoClearPinsHandler).Methods(http.MethodPost)
	r.HandleFunc("/arduino/pins", arduinoPinsAllHandler).Methods(http.MethodGet)
	r.HandleFunc("/arduino/add_reset", arduinoAddResetHandler).Methods(http.MethodPost)

	r.HandleFunc("/temperatures", temperatureStatisticHandler).Methods(http.MethodGet)

	r.HandleFunc("/fans", fansSpeedHandler).Methods(http.MethodPost)
	r.HandleFunc("/card_state", cardStateHandler).Methods(http.MethodPost)

	r.HandleFunc("/settings", settingsGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/settings", settingsAddHandler).Methods(http.MethodPost)

	r.HandleFunc("/wallet", walletGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/wallet", walletAddHandler).Methods(http.MethodPost)
	r.HandleFunc("/wallet/delete", walletDeleteHandler).Methods(http.MethodPost)
	r.HandleFunc("/wallet/set_default", walletSetDefaultHandler).Methods(http.MethodPost)

	// options handlers
	r.HandleFunc("/reboot", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/restart", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/shutdown", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/client", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/client/updateapp", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/client/delete", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/arduino/add_reset", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/arduino/pins", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/arduino/pins/clear", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/settings", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/fans", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/wallet", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/wallet/delete", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/wallet/set_default", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/card_state", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/claymore/config/basic", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/claymore/config/basic/delete", optionsHandler).Methods(http.MethodOptions)

	log.Println("Routers successfully inited..")

	// to show api
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {

		t, err := route.GetPathTemplate()
		if err != nil {
			logger.Error("Problem with geting api path..")
			return err
		}

		api += t + "\n"
		return nil
	})
	return r
}

func InitWebServer() {
	logger.Info("I am starting web server..")
	logger.Info("Listening on port " + config.GetWebPort())
	err := http.ListenAndServe(":"+config.GetWebPort(), initRouters())

	if err != nil {
		logger.Error("Something went wrong when init web server..")
		logger.Error(err.Error())
		os.Exit(1)
		return
	}
}

func getCustomClientsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var clients []client.CustomClient
	data.Find(&clients)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clients)
}

func getClaymoreBasicConfigHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	params := mux.Vars(r)

	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var basicConfig internal.ClayMoreBasicConfig
	data.Where("name = ?", params["name"]).First(&basicConfig)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(basicConfig)
}

func getClaymoreBasicConfigsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var basicConfigs []internal.ClayMoreBasicConfig
	data.Find(&basicConfigs)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(basicConfigs)
}

func deleteBasicConfigHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var basicConfig internal.ClayMoreBasicConfig
	_ = json.NewDecoder(r.Body).Decode(&basicConfig)

	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data.Delete(&basicConfig)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully deleted basic config.."))
}

func deleteCustomClientHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var clientInfo client.CustomClient
	_ = json.NewDecoder(r.Body).Decode(&clientInfo)

	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data.Where("name = ?", clientInfo.Name).Delete(client.CustomClient{})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully deleted client.."))
}

func editBasicConfigHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	params := mux.Vars(r)

	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var basicConfig internal.ClayMoreBasicConfig

	data.Where("name = ?", params["name"]).First(&basicConfig)
	data.Delete(&basicConfig)

	_ = json.NewDecoder(r.Body).Decode(&basicConfig)

	data.Create(&basicConfig)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully updated basic config.."))
}

func createBasicConfigHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var basicConfig internal.ClayMoreBasicConfig
	_ = json.NewDecoder(r.Body).Decode(&basicConfig)

	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	data.Create(&basicConfig)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully created basic config.."))
}

func createCustomClientHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var customClient client.CustomClient
	_ = json.NewDecoder(r.Body).Decode(&customClient)

	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	data.Where("name = ?", customClient.Name).Assign(customClient).FirstOrCreate(&customClient)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully created custom client.."))
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	/*fmt.Fprintf(w,

	"/\n\n"+

	"/claymore/manager/{name}/{cfg}/{option}\n"+
	"/claymore/log/{name}/limit/{p1}/{p2}\n"+
	"/claymore/log/{name}/date/{p1}/{p2}\n"+
	"/claymore/log/{name}\n\n"+

	"/hardware/log/{name}/limit/{p1}/{p2}\n"+
	"/hardware/log/{name}/date/{p1}/{p2}\n"+
	"/hardware/log/{name}\n\n"+

	"/claymore/config/basic\n"+
	"/claymore/config/export\n"+
	"/claymore/config/basic/{name}\n"+
	"/claymore/config/export/{name}\n\n"+

	"/reboot\n\n"+

	"/client\n"+
	"/client/{name}\n\n"+

	"/arduino\n\n"+

	"/clients\n\n"+

	"/latest/stat\n\n")
	*/
	fmt.Fprint(w, api)
}

func claymoreHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	option := vars["option"]

	if option == "off" {
		programs.StopClaymore(name)
		return
	}

	if option == "on" {
		programs.StartClaymore(name, vars["cfg"])
	}

}

func showClaymoreStatLogsWithLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	var stats []internal.ClaymoreStat

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("miner_name = ?", name).Order("time desc").Offset(p1).Limit(p2).Find(&stats)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func showClaymoreStatLogsWithDateLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	var stats []internal.ClaymoreStat

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("miner_name = ? AND time between ? AND ?", name, p1, p2).Order("time desc").Find(&stats)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func showClaymoreStatLogsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	var stats []internal.ClaymoreStat

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("miner_name = ?", name).Find(&stats)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func showHardwareStatLogsWithLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	var stats []internal.HardwareStat

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("miner_name = ?", name).Order("time desc").Offset(p1).Limit(p2).Find(&stats)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func showHardwareStatLogsWithDateLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	var stats []internal.HardwareStat

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("miner_name = ? AND time between ? AND ?", name, p1, p2).Order("time desc").Find(&stats)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func showHardwareStatLogsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	var stats []internal.HardwareStat

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("miner_name = ?", name).Find(&stats)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func arduinoTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tempInfo internal.TempInfo

	data := make(map[string]string)
	for key, val := range r.URL.Query() {
		data[key] = val[0]
	}

	jsonData, _ := json.Marshal(data)
	tempInfo.Data = string(jsonData)
	tempInfo.Time = util.CurrentTime()

	db.Create(&tempInfo)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func arduinoResetHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v := "|"

	var machinesToReset []internal.HardReset
	db.Find(&machinesToReset)

	for _, value := range machinesToReset {
		v += value.MachineNumber + "R|"
	}

	var machinesToShutdown []internal.HardShutdown
	db.Find(&machinesToShutdown)
	for _, value := range machinesToShutdown {
		if value.Function == "0" {
			v += value.MachineNumber + "D|"
		} else {
			v += value.MachineNumber + "U|"
		}
	}

	fmt.Fprint(w, v)
	db.Delete(&machinesToReset)
	db.Delete(&machinesToShutdown)
	w.WriteHeader(http.StatusOK)

}

func getClientsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	clients := make(map[string]internal.Miner)
	clients = connection.GetConnectedClients()

	names := make([]string, 0, len(clients))

	for name := range clients {
		names = append(names, name)
	}

	var customClients []client.CustomClient

	if len(names) == 0 {
		db.Find(&customClients)
	} else {
		db.Where("name NOT IN (?)", names).Find(&customClients)
	}
	for _, customClient := range customClients {
		clients[customClient.Name] = internal.Miner{MinerName: customClient.Name}
	}

	json.NewEncoder(w).Encode(clients)
}

func getLatestStatHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var claymoreLogs []internal.ClaymoreStat
	var claymoreStats []string

	if license.IsValidated() {
		db.Table("claymore_stats").Group("miner_name").Pluck("MAX(time)", &claymoreStats)
		db.Where("time IN (?)", claymoreStats).Find(&claymoreLogs)
	}

	miners := make(map[string]internal.ClaymoreStat)
	for _, val := range claymoreLogs {
		miners[val.MinerName] = val
	}

	minersJSON, _ := json.Marshal(miners)

	var resets []internal.HardReset

	if license.IsValidated() {
		db.Find(&resets)
	}

	var v []int
	for _, value := range resets {
		number, _ := strconv.Atoi(value.MachineNumber)
		v = append(v, number)
	}
	resetsJSON, _ := json.Marshal(v)

	var hardwareStats []string
	var hardwareLogs []internal.HardwareStat

	if license.IsValidated() {
		db.Table("hardware_stats").Group("miner_name").Pluck("MAX(time)", &hardwareStats)
		db.Where("time IN (?)", hardwareStats).Find(&hardwareLogs)
	}

	hardware := make(map[string]internal.HardwareStat)
	for _, val := range hardwareLogs {
		hardware[val.MinerName] = val
	}

	hardwareJSON, _ := json.Marshal(hardware)

	type CustomClient struct {
		Id         int    `json:"id"`
		Name       string `json:"Name"`
		Config     string `json:"Config"`
		WalletID   int    `json:"WalletID" sql:"type:int"`
		Currency   string `json:"Currency"`
		WalletName string `json:"WalletName"`
	}

	var clientsData []CustomClient
	db.Table("custom_clients").Select("*").Joins("join wallets on custom_clients.wallet_id = wallets.id").Find(&clientsData)

	clients := make(map[string]CustomClient)
	for _, val := range clientsData {
		clients[val.Name] = val
	}

	clientsJSON, _ := json.Marshal(clients)

	licenseStatus := strconv.FormatBool(license.IsValidated())

	output := "{\"miners\":" + string(minersJSON) + ", \"resets\":" + string(resetsJSON) + ", \"hardware\":" + string(hardwareJSON) + ", \"config\":" + string(clientsJSON) + ", \"license\": " + licenseStatus + "}"
	fmt.Fprint(w, output)
}

func rebootClientHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	var rebootInfo internal.Reboot
	_ = json.NewDecoder(r.Body).Decode(&rebootInfo)
	rebootInfo.Time = util.CurrentTime()

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Create(&rebootInfo)
	connection.RebootClient(rebootInfo.MinerName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func getRebootInfoHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rebootInfo []internal.Reboot
	db.Where("miner_name = ?", name).Find(&rebootInfo)
	json.NewEncoder(w).Encode(rebootInfo)
}

func getRebootLastInfoHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rebootInfo []internal.Reboot
	var miners []int

	db.Table("reboots").Group("miner_name").Find(&miners).Pluck("MAX(id)", &miners)
	db.Where("id IN (?)", miners).Find(&rebootInfo)
	json.NewEncoder(w).Encode(rebootInfo)
}

func getRebootInfoWithLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rebootInfo []internal.Reboot
	db.Where("miner_name = ?", name).Order("time desc").Offset(p1).Limit(p2).Find(&rebootInfo)
	json.NewEncoder(w).Encode(rebootInfo)
}

func getRebootInfoWithDateLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rebootInfo []internal.Reboot
	db.Where("miner_name = ? AND time between ? AND ?", name, p1, p2).Order("time desc").Find(&rebootInfo)
	json.NewEncoder(w).Encode(rebootInfo)
}

func createRecipientHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var email mail.Recipient
	_ = json.NewDecoder(r.Body).Decode(&email)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	db.Create(&email)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully added email to recipeints.."))
}

func deleteRecipientHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var email mail.Recipient
	_ = json.NewDecoder(r.Body).Decode(&email)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("email = ?", email.Email).First(&email)
	db.Delete(&email)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully deleted email from recipients.."))
}

func restartAppHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var restartInfo internal.Restart
	_ = json.NewDecoder(r.Body).Decode(&restartInfo)
	restartInfo.Time = util.CurrentTime()

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Create(&restartInfo)

	programs.StopClaymore(restartInfo.MinerName)
	programs.StartClaymore(restartInfo.MinerName, restartInfo.Config)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func getRestartInfoHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var restartInfo []internal.Restart
	db.Where("miner_name = ?", name).Find(&restartInfo)
	json.NewEncoder(w).Encode(restartInfo)
}

func getRestartLastInfoHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var restartInfo []internal.Restart
	var miners []int

	db.Table("restarts").Group("miner_name").Find(&miners).Pluck("MAX(id)", &miners)
	db.Where("id IN (?)", miners).Find(&restartInfo)
	json.NewEncoder(w).Encode(restartInfo)
}

func getRestartInfoWithLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var restartInfo []internal.Restart
	db.Where("miner_name = ?", name).Order("time desc").Offset(p1).Limit(p2).Find(&restartInfo)
	json.NewEncoder(w).Encode(restartInfo)
}

func getRestartInfoWithDateLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var restartInfo []internal.Restart
	db.Where("miner_name = ? AND time between ? AND ?", name, p1, p2).Order("time desc").Find(&restartInfo)
	json.NewEncoder(w).Encode(restartInfo)
}

func shutdownClientHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	var shutdownClient internal.Shutdown
	_ = json.NewDecoder(r.Body).Decode(&shutdownClient)
	shutdownClient.Time = util.CurrentTime()

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Create(&shutdownClient)
	connection.ShutdownClient(shutdownClient.MinerName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func getShutdownInfoHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shutdownInfo []internal.Shutdown
	db.Where("miner_name = ?", name).Find(&shutdownInfo)
	json.NewEncoder(w).Encode(shutdownInfo)
}

func getShutdownLastInfoHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shutdownInfo []internal.Shutdown
	var miners []int

	db.Table("shutdowns").Group("miner_name").Find(&miners).Pluck("MAX(id)", &miners)
	db.Where("id IN (?)", miners).Find(&shutdownInfo)
	json.NewEncoder(w).Encode(shutdownInfo)
}

func getShutdownInfoWithLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shutdownInfo []internal.Shutdown
	db.Where("miner_name = ?", name).Order("time desc").Offset(p1).Limit(p2).Find(&shutdownInfo)
	json.NewEncoder(w).Encode(shutdownInfo)
}

func getShutdownInfoWithDateLimitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shutdownInfo []internal.Shutdown
	db.Where("miner_name = ? AND time between ? AND ?", name, p1, p2).Order("time desc").Find(&shutdownInfo)
	json.NewEncoder(w).Encode(shutdownInfo)
}

func temperatureStatisticHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tempInfo []internal.TempInfo
	db.Find(&tempInfo)

	for _, data := range tempInfo {
		fmt.Fprintf(w, data.Data)
	}
}

func getTemperatureLastInfoHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tempInfo internal.TempInfo
	db.Last(&tempInfo)

	fmt.Fprintf(w, tempInfo.Data)
}

func fansSpeedHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	data := struct {
		Id      string `json:"id"`
		Machine string `json:"machine"`
		Speed   string `json:"speed"`
	}{"", "", ""}

	_ = json.NewDecoder(r.Body).Decode(&data)

	speed, _ := strconv.Atoi(data.Speed)
	if speed > 100 || speed < 0 {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Wrong speed >100 or <0"))
		w.WriteHeader(http.StatusBadRequest)
	}

	programs.SetFanSpeed(data.Machine, data.Id, data.Speed)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func arduinoAddPinsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	var pinFunction internal.PinFunction

	_ = json.NewDecoder(r.Body).Decode(&pinFunction)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db.Save(&pinFunction)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")

}

func arduinoPinsAllHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var pinFunction []internal.PinFunction
	db.Find(&pinFunction)
	json.NewEncoder(w).Encode(pinFunction)
}

func arduinoClearPinsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db.Delete(internal.PinFunction{})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")

}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	w.WriteHeader(http.StatusOK)
}

func getDisconectedClientsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	clients := make(map[string]internal.Miner)
	clients = connection.GetConnectedClients()

	names := make([]string, 0, len(clients))

	for name := range clients {
		names = append(names, name)
	}

	var customClients []client.CustomClient

	if len(names) == 0 {
		db.Find(&customClients)
	} else {
		db.Where("name NOT IN (?)", names).Find(&customClients)
	}

	for _, customClient := range customClients {
		names = append(names, customClient.Name)
	}

	var claymoreStats []internal.ClaymoreStat

	if len(names) == 0 {
		db.Group("miner_name").Find(&claymoreStats)
	} else {
		db.Where("miner_name NOT IN (?)", names).Group("miner_name").Find(&claymoreStats)
	}

	disconected := make([]string, 0, len(claymoreStats))

	for _, stat := range claymoreStats {
		disconected = append(disconected, stat.MinerName)
	}

	json.NewEncoder(w).Encode(disconected)
}

func updateClientAppHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var data map[string]string
	_ = json.NewDecoder(r.Body).Decode(&data)
	connection.UpdateClient(data["Name"])
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully sent update request.."))
}

func arduinoAddResetHandler(w http.ResponseWriter, r *http.Request) {
	// 0 - reset
	// 1 - shutdown
	enableCors(w)
	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := struct {
		MachineName string `json:"MachineName"`
		Function    string `json:"Function"`
	}{"", ""}

	_ = json.NewDecoder(r.Body).Decode(&data)

	if data.Function == "reset" {
		var pin []string
		db.Table("pin_functions").Where("miner_name = ?", data.MachineName).Where("function = ?", "0").Pluck("id", &pin)
		var machineToReset internal.HardReset
		machineToReset.MachineNumber = pin[0]
		db.Save(&machineToReset)
	}

	if data.Function == "shutdown" {
		var pin []string
		db.Table("pin_functions").Where("miner_name = ?", data.MachineName).Where("function = ?", "1").Pluck("id", &pin)
		var machineToShutdown internal.HardShutdown
		machineToShutdown.Function = "0"
		machineToShutdown.MachineNumber = pin[0]
		db.Save(&machineToShutdown)
	}

	if data.Function == "poweron" {
		var pin []string
		db.Table("pin_functions").Where("miner_name = ?", data.MachineName).Where("function = ?", "1").Pluck("id", &pin)
		var machineToShutdown internal.HardShutdown
		machineToShutdown.Function = "1"
		machineToShutdown.MachineNumber = pin[0]
		db.Save(&machineToShutdown)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func settingsGetHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var settings []internal.Settings
	db.Find(&settings)
	settingsMap := make(map[string]string)
	for _, val := range settings {
		settingsMap[val.Name] = val.Value
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(settingsMap)
}

func settingsAddHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	var settings internal.Settings
	_ = json.NewDecoder(r.Body).Decode(&settings)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.Where("name = ?", settings.Name).Assign(settings).FirstOrCreate(&settings)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Settings Saved.."))
}

func walletGetHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var wallets []internal.Wallet
	db.Find(&wallets)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallets)
}

func walletAddHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	var wallet internal.Wallet
	_ = json.NewDecoder(r.Body).Decode(&wallet)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	db.Where("id = ?", wallet.ID).Assign(wallet).FirstOrCreate(&wallet)
	json.NewEncoder(w).Encode(wallet.ID)
}

func walletDeleteHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	var wallet internal.Wallet
	_ = json.NewDecoder(r.Body).Decode(&wallet)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	db.Where("id = ?", wallet.ID).Delete(&wallet)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully deleted wallet.."))
}

func walletSetDefaultHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	var wallet internal.Wallet
	_ = json.NewDecoder(r.Body).Decode(&wallet)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

	db.Model(&internal.Wallet{}).Where("currency = ?", wallet.Currency).Update("is_default", "0")
	db.Model(&wallet).Where("currency = ?", wallet.Currency).Where("id = ?", wallet.ID).Update("is_default", "1")
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully set wallet as default.."))
}

func cardStateHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	var cardInfo internal.CardDisable
	_ = json.NewDecoder(r.Body).Decode(&cardInfo)

	db, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer db.Close()

	if err != nil {
		json.NewEncoder(w).Encode(internal.CreateHandlerError(1, "Can't connect to database.."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	miner := cardInfo.Miner
	id := cardInfo.CardID
	state := cardInfo.State
	programs.CardState(miner, id, state)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(internal.CreateHandlerError(0, "Successfully set GPU state.."))
}
