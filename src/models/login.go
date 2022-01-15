package models

type Login struct {
	Login  string `json:"login"`
	Secret string `json:"secret"`
}
