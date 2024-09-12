package models

type User struct {
	UID      int64
	IP       string
	Email    string
	PassHash []byte
	Secret   string
	Token
}
