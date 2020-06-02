/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package mail

type Recipient struct {
	ID    int    `json:"id gorm:"AUTO_INCREMENT"`
	Email string `json:"email" gorm:"UNIQUE"`
}

type MailGunConfig struct {
	Domain     string `json:"domain"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Sender     string `json:"sender"`
}

type EmailHistory struct {
	Recipients string `json:"recipients"`
	Title      string `json:"title"`
	Message    string `json:"message"`
}

func join(rec []Recipient) string {
	var s = ""

	for _, val := range rec {
		s = s + val.Email + ","
	}

	return s
}
