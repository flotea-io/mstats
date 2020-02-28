package internal

type Settings struct {
	ID    int    `json:"id" gorm:"AUTO_INCREMENT"`
	Name  string `json:"Name" gorm:"unique_index"`
	Value string `json:"Value"`
}
