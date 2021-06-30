package models

type Building struct {
	Id            int64           `json:"id,omitempty"`
	Address       string          `json:"address"`
	Latitude      float64         `json:"latitude"`
	Longitude     float64         `json:"longitude"`
	Organizations *[]Organization `json:"organizations,omitempty"`
}
