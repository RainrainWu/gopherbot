package main

import (

	"log"

	"github.com/RainrainWu/gopherbot/db"
)

func main()  {

	db.CreateResource("status", "status_url")
	db.CreateResource("status2", "status2_url")
	db.CreateResource("status2", "status2_url")
	log.Print(db.ListResources())
	db.CreateTeam("sponsorship")
	db.CreateTeam("program")
	db.CreateTeam("program")
	log.Print(db.ListTeams())
	return
}