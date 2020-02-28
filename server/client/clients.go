package client

type CustomClient struct {
	Id       int    `json:"id"`
	Name     string `json:"Name"`
	Config   string `json:"Config" gorm:"default:\"autostart\""`
	WalletID int    `json:"WalletID" sql:"type:int"`
}
