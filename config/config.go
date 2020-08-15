package config

import (

	"fmt"
	"io/ioutil"
	"encoding/json"
)

type UsingConfig struct {

	DBConfig DatabaseConfig	`json:"dbconfig"`
	TGConfig TelegramConfig	`json:"tgconfig"`
}

type DatabaseConfig struct {

	Host string		`json:"host"`
	Port string		`json:"port"`
	DBname string 	`json:"dbname"`
	Username string	`json:"username"`
	Password string	`json:"password"`
	SSLmode string	`json:"sslmode"`
}

type TelegramConfig struct {

	Token string	`json:token`
}

var (

	Config = UsingConfig{

		DBConfig: DatabaseConfig{},
		TGConfig: TelegramConfig{},
	}
)

func init()  {

	LoadConfig()
}

func LoadConfig() {

	data, _ := ioutil.ReadFile("./.env")
	err := json.Unmarshal(data, &Config)
	if err != nil {
		fmt.Println(err)
	}
}