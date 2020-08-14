package config

import (

	"fmt"
	"io/ioutil"
	"encoding/json"
)

type DatabaseConfig struct {

	Host string		`json:"host"`
	Port string		`json:"port"`
	DBname string 	`json:"dbname"`
	Username string	`json:"username"`
	Password string	`json:"password"`
	SSLmode string	`json:"sslmode"`
}

var (

	UsingConfig = DatabaseConfig{}
)

func init()  {

	LoadConfig()
}

func LoadConfig() {

	data, _ := ioutil.ReadFile("./.env")
	err := json.Unmarshal(data, &UsingConfig)
	if err != nil {
		fmt.Println(err)
	}
}