package internal

type Wallet struct {
	ID         int    `json:"ID" gorm:"AUTO_INCREMENT unique_index"`
	Address    string `json:"Address"`
	Currency   string `json:"Currency"`
	WalletName string `json:"WalletName"`
	IsDefault  string `json:"IsDefault" sql:"type:int(1) unsigned" gorm:"default:0"`
}
