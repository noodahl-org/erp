package api

import "fmt"

type ApiConf struct {
	Port        int `validate:"gt=0"`
	BaseUrl     string
	BraveAPIKey string
	DbConf
}

type DbConf struct {
	DbUser string
	DbPass string
	DbHost string
	DbName string
	DbPort int
}

func (c *DbConf) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%v database=%s",
		c.DbHost, c.DbUser, c.DbPass, c.DbPort, c.DbName)
}
