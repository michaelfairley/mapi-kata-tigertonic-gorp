package api

import ()

type Auth struct {
	Username string
	Password string
}

type Token struct {
	Value string `json:"token"`
}
