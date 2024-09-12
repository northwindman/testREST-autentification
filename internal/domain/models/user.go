package models

type User struct {
	ID       int64
	Email    string
	PassHash []byte
	IP       string
	Token
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
