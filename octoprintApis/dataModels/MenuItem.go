package dataModels


type MenuItem struct {
	Name  string     `json:"name"`
	Icon  string     `json:"icon"`
	Panel string     `json:"panel"`
	Items []MenuItem `json:"items"`
}
