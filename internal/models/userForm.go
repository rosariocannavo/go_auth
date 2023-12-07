package models

type UserForm struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	MetamaskAddress string `json:"metamaskAddress"`
}
