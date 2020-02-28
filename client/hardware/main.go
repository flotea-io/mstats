package hardware

import (
	"io/ioutil"
	"math"
	"mstats-new/internal"
	"mstats-new/logger"
	"mstats-new/util"
	"os"
	"regexp"
	"strconv"
)

var Cards = make(map[string]internal.Card)

const mainDir = "/sys/class/drm/"

//const mainDir = "./"

func SetUpMachine() {
	files, err := ioutil.ReadDir(mainDir)

	if err != nil {
		logger.Error(err.Error())
	}

	cardRegexp := regexp.MustCompile("^card[0-9]$").MatchString
	id := 1
	for _, file := range files {
		if !cardRegexp(file.Name()) {
			continue
		}
		//check if card have hwmon
		if _, err := os.Stat(mainDir + file.Name() + "/device/hwmon/"); os.IsNotExist(err) {
			continue
		}

		//get cards data
		card := readCardInfo(file.Name())

		card.CardID = strconv.Itoa(id)

		//save card
		Cards[card.CardID] = card
		id++
	}
}

func readCardInfo(cardName string) internal.Card {
	card := internal.Card{"", "", "", false, 0, 0, 0, 0, 0, 80, "0"}

	monitorRegexp := regexp.MustCompile("^hwmon[0-9]$").MatchString
	//get card monitor
	monitors, err := ioutil.ReadDir(mainDir + cardName + "/device/hwmon/")
	if err != nil {
		logger.Error(err.Error())
		return card
	}

	enabled := util.ReadFileToInt(mainDir + cardName + "/device/enable")
	card.Enabled = strconv.Itoa(enabled)

	monitor := monitors[0]
	if !monitorRegexp(monitor.Name()) {
		return card

	}

	maxValue := util.ReadFileToInt(mainDir + cardName + "/device/hwmon/" + monitor.Name() + "/pwm1_max")
	currentValue := util.ReadFileToInt(mainDir + cardName + "/device/hwmon/" + monitor.Name() + "/pwm1")
	currentRPM := util.ReadFileToInt(mainDir + cardName + "/device/hwmon/" + monitor.Name() + "/fan1_input")
	temp := util.ReadFileToInt(mainDir+cardName+"/device/hwmon/"+monitor.Name()+"/temp1_input") / 1000
	currentPercent := 100 * currentValue / maxValue

	/*
		if maxValue < 0 || currentValue < 0 || currentRPM < 0 || currentPercent < 0 || temp < 0 || temp > 150 {
			logger.Error("Something went wrong when reading values")
			return	Card{"", "", false, 0, 0, 0, 0, 0}
		}
	*/
	card.CardName = cardName
	card.MonitorName = monitor.Name()
	card.Temperature = temp
	card.MaxValue = maxValue
	card.CurrentValue = currentValue
	card.CurrentRPM = currentRPM
	card.CurrentPercent = currentPercent
	return card
}

func updateCardInfo(cardID string) internal.Card {
	card, exist := Cards[cardID]

	if !exist {
		card = internal.Card{"", "", "", false, 0, 0, 0, 0, 0, 80, "0"}
		return card
	}

	maxValue := util.ReadFileToInt(mainDir + card.CardName + "/device/hwmon/" + card.MonitorName + "/pwm1_max")
	currentValue := util.ReadFileToInt(mainDir + card.CardName + "/device/hwmon/" + card.MonitorName + "/pwm1")
	currentRPM := util.ReadFileToInt(mainDir + card.CardName + "/device/hwmon/" + card.MonitorName + "/fan1_input")
	temp := util.ReadFileToInt(mainDir+card.CardName+"/device/hwmon/"+card.MonitorName+"/temp1_input") / 1000
	currentPercent := 100 * currentValue / maxValue

	/*
		if maxValue < 0 || currentValue < 0 || currentRPM < 0 || currentPercent < 0 || temp < 0 || temp > 150 {
			logger.Error("Something went wrong when reading values")
			return	Card{"", "", false, 0, 0, 0, 0, 0}
		}
	*/
	card.Temperature = temp
	card.MaxValue = maxValue
	card.CurrentValue = currentValue
	card.CurrentRPM = currentRPM
	card.CurrentPercent = currentPercent
	return card
}

func SetFanSpeed(cardID string, speed string) {

	percent, err := strconv.Atoi(speed)

	if err != nil {
		logger.Error("Something went wrong when converting value in " + cardID)
		return
	}

	if percent < 0 || percent > 100 {
		logger.Error("Wrong fan speed at card " + cardID)
		return
	}

	for id, card := range Cards {
		if card.CardID == cardID || cardID == "all" {
			card.DeclaredPercent = percent
			if percent != 0 {
				card.DeclaredPercent = percent
				card.ManualFanControl = true
			} else {
				card.DeclaredPercent = int(util.Round(float64(card.Temperature), 10))
				card.ManualFanControl = false
			}
			Cards[id] = card
			setCardSpeed(card)

			//reload card info
			Cards[id] = updateCardInfo(id)
		}
	}
}

//write speed for a specific card to file
func setCardSpeed(card internal.Card) {
	if math.Abs(float64(card.DeclaredPercent-card.CurrentPercent)) < 5 {
		return
	}

	value := int(math.Round(float64(card.MaxValue*card.DeclaredPercent)) / 100)
	filePath := mainDir + card.CardName + "/device/hwmon/" + card.MonitorName + "/pwm1"

	file, err := os.Create(filePath)
	if err != nil {
		logger.Error("Cannot open file " + filePath + "to writing")
		return
	}
	defer file.Close()
	logger.Info("Setting speed at " + card.CardName + " to " + strconv.Itoa(card.DeclaredPercent) + "%%")
	if _, err := file.Write([]byte(strconv.Itoa(value))); err != nil {
		logger.Error("Cannot write speed to file " + filePath)
	}
}

//write card state (disabled/enabled) to files
func SetCardState(cardID string, state string) {
	for id, card := range Cards {
		if id == cardID || cardID == "all" {

			card.Enabled = state
			Cards[id] = card
			setState(card)
			//reload card info
			Cards[id] = updateCardInfo(id)
		}
	}
}

//write state to files
func setState(card internal.Card) {
	filePath := mainDir + card.CardName + "/device/enable"

	file, err := os.Create(filePath)
	if err != nil {
		logger.Error("Cannot open file " + filePath + " to writing")
		return
	}
	defer file.Close()

	logger.Info("Setting state at " + card.CardName + " to " + card.Enabled)
	if _, err := file.Write([]byte(card.Enabled)); err != nil {
		logger.Error("Cannot write state to file " + filePath)
	}
}
