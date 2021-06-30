package models

type Organization struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Categories   []int64   `json:"categories"`
	PhoneNumbers []string  `json:"phone_numbers"`
	Building     *Building `json:"building"`
}
