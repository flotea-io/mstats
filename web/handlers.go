package web

import (
	"fmt"
	"log"
	"mstats/logger"
	"mstats/mail"
	"mstats/miner"
	"mstats/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
}

func initRouters() *mux.Router {
	log.Println("Initing routers..")
	r := mux.NewRouter()

	r.HandleFunc("/", basicHandler)
	r.HandleFunc("/miners", minerListHandler)
	r.HandleFunc("/miners/remove/{name}", removeHandler)
	r.HandleFunc("/miners/edit/{name}/{option}/{value}", editMinerHandler)
	r.HandleFunc("/miners/add/{name}/{ip}/{port}", addNewMinerHandler)

	r.HandleFunc("/latest/stat", latestStatLogHandler)
	r.HandleFunc("/latest/reboot", latestRebootLogHandler)
	r.HandleFunc("/latest/restart", latestRestartLogHandler)

	r.HandleFunc("/reboot/{name}/{reason}", rebootHandler)
	r.HandleFunc("/restart/{name}/{reason}", restartHandler)

	r.HandleFunc("/log/stat/{name}/limit/{p1}/{p2}", showStatLogsHandlerWithLimit)
	r.HandleFunc("/log/restart/{name}/limit/{p1}/{p2}", showRestartLogsHandlerWithLimit)
	r.HandleFunc("/log/reboot/{name}/limit/{p1}/{p2}", showRebootLogsHandlerWithLimit)

	r.HandleFunc("/log/stat/{name}/date/{p1}/{p2}", showStatLogsHandlerWithDateLimit)
	r.HandleFunc("/log/restart/{name}/date/{p1}/{p2}", showRestartLogsHandlerWithDateLimit)
	r.HandleFunc("/log/reboot/{name}/date/{p1}/{p2}", showRebootLogsHandlerWithDateLimit)

	r.HandleFunc("/log/stat/{name}", showStatLogsHandler)
	r.HandleFunc("/log/reboot/{name}", showRebootLogsHandler)
	r.HandleFunc("/log/restart/{name}", showRestartLogsHandler)

	r.HandleFunc("/mail/add/{name}", showRebootLogsHandler)
	r.HandleFunc("/mail/delete/{name}", showRestartLogsHandler)

	r.HandleFunc("/time", timeHandler)

	log.Println("Routers successfully inited..")
	return r
}

func InitWebServer() {
	log.Println("Miner app starting...")
	log.Println("Listening on port 9922")
	err := http.ListenAndServe(":9922", initRouters())

	if err != nil {
		log.Fatalln("Someting went wrong... ", err)
		return
	}

}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	fmt.Fprintf(w, "/miners\n"+
		"/miners/add/<name>/<ip>/<port>\n"+
		"/miners/remove/<name>\n"+
		"/miners/edit/<name>/<name/ip/port>/<value>\n"+
		"\n"+
		"/latest/stat\n"+
		"/latest/restart\n"+
		"/latest/reboot\n"+
		"\n"+
		"/reboot/<name>/reason_with_floor\n"+
		"/restart/<name>/reason_with_floor\n"+
		"\n"+
		"/log/stat/<name>\n"+
		"/log/reboot/<name>\n"+
		"/log/restart/<name>\n"+
		"\n"+
		"/log/stat/<name>/limit/<limit_1>/<limit_2>\n"+
		"/log/reboot/<name>/limit/<limit_1>/<limit_2>\n"+
		"/log/restart/<name>/limit/<limit_1>/<limit_2>\n"+
		"\n"+
		"/log/stat/<name>/date/<2010-04-05>/<2020-04-20>\n"+
		"/log/reboot/<name>/date/<2010-04-05>/<2020-04-20>\n"+
		"/log/restart/<name>/date/<2010-04-05>/<2020-04-20>\n"+
		"\n"+
		"/mail/add/<mail>\n"+
		"/mail/delete/<mail>\n"+
		"/time")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	fmt.Fprintf(w, util.CurrentTime())
}

func minerListHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	fmt.Fprintf(w, miner.GetMinersJson())
}

func latestStatLogHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	fmt.Fprintf(w, logger.GetLastStatJson())
}

func latestRebootLogHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	fmt.Fprintf(w, logger.GetLastRebootJson())
}

func latestRestartLogHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	fmt.Fprintf(w, logger.GetLastRestartJson())
}

func showRebootLogsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, logger.GetRebootLogsFromMiner(miner.GetMinerByName(name).MinerIP))
}

func mailDeleteHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Mail not found\"}")
		return
	}

	fmt.Fprintf(w, mail.DeleteRecipient(name))
}

func mailAddHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	if mail.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Mail has already been added\"}")
		return
	}

	fmt.Fprintf(w, mail.AddRecipient(name))
}

func showRebootLogsHandlerWithLimit(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, logger.GetRebootLogsFromMinerWithLimit(miner.GetMinerByName(name).MinerIP, p1, p2))
}

func editMinerHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	option := strings.ToLower(vars["option"])
	value := vars["value"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, miner.EditMiner(miner.GetMinerByName(name), option, value))

}

func addNewMinerHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	ip := vars["ip"]
	port := vars["port"]
	i, e := strconv.Atoi(port)

	if e != nil {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Error where parsing port to integer\"}")
		return
	}

	fmt.Fprintf(w, miner.CreateMiner(name, ip, i))

}

func showRestartLogsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "This miner is not exist..")
		return
	}

	fmt.Fprintf(w, logger.GetRestartLogsFromMiner(miner.GetMinerByName(name).MinerIP))
}

func showRestartLogsHandlerWithLimit(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, logger.GetRestartLogsFromMinerWithLimit(miner.GetMinerByName(name).MinerIP, p1, p2))
}

func showStatLogsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, logger.GetStatsLogsFromMiner(miner.GetMinerByName(name).MinerIP))
}

func showStatLogsHandlerWithLimit(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, logger.GetStatsLogsFromMinerWithLimit(miner.GetMinerByName(name).MinerIP, p1, p2))
}

func showStatLogsHandlerWithDateLimit(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, logger.GetStatsLogsFromMinerWithDateLimit(miner.GetMinerByName(name).MinerIP, p1, p2))
}

func showRestartLogsHandlerWithDateLimit(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, logger.GetRestartLogsFromMinerWithDateLimit(miner.GetMinerByName(name).MinerIP, p1, p2))
}

func showRebootLogsHandlerWithDateLimit(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	p1 := vars["p1"]
	p2 := vars["p2"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, logger.GetRebootLogsFromMinerWithDateLimit(miner.GetMinerByName(name).MinerIP, p1, p2))
}

func rebootHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	reason := vars["reason"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	go miner.ReBoot(miner.GetMinerByName(name), reason)
	fmt.Fprintf(w, miner.Restart(miner.GetMinerByName(name), reason))
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]
	reason := vars["reason"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, miner.Restart(miner.GetMinerByName(name), reason))
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	vars := mux.Vars(r)
	name := vars["name"]

	if !miner.IsExist(name) {
		fmt.Fprintf(w, "{\"status\":\"1\",\"log\":\"Miner is not exist\"}")
		return
	}

	fmt.Fprintf(w, miner.RemoveMiner(miner.GetMinerByName(name)))
}
