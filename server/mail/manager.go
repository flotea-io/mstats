package mail

import (
	"mstats-new/logger"
	"mstats-new/server/config"
	"os"

	"github.com/jinzhu/gorm"
)

func sendMessage(mg mailgun.Mailgun, subject, body string) {
	data, err := gorm.Open("mysql", config.GetDatabaseDSN())

	defer data.Close()

	if err != nil {
		logger.Error("Problem with connect to database..")
		logger.Error(err.Error())
		os.Exit(1)
		return
	}

	var cfg MailGunConfig

	data.First(&cfg)

	var recipients []Recipient

	data.Find(&recipients)

	for _, key := range recipients {
		message := mg.NewMessage(cfg.Sender, subject, body, key.Email)
		_, _, err := mg.Send(message)

		if err != nil {
			logger.Error("Something went wrong with sending email..")
			logger.Error(err.Error())
			continue
		}
	}

	data.Create(&EmailHistory{Recipients: join(recipients), Title: subject, Message: body})

}
