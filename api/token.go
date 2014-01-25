package api

import (
	"crypto/rand"
	"fmt"
)

type Auth struct {
	Username string
	Password string
}

type Token struct {
	Value  string `json:"token"`
	UserId uint64 `json:"-" db:"user_id"`
}

func GenerateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	b[8] = (b[8] | 0x80) & 0xBF
	b[6] = (b[6] | 0x40) & 0x4F

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
