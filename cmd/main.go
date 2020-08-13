package main

import (

	"github.com/RainrainWu/gopherbot/db"
)

func main()  {

	db.CreateResource("status", "status_url")
	db.CreateResource("status2", "status2_url")
	db.CreateTeam("sponsorship")
	db.RegisterResource("status", "sponsorship")
	db.RegisterResource("status2", "sponsorship")
	db.QueryResources("sponsorship")
	db.DisconnectDatabase()
	return
}