package handlers

import (

	"fmt"
	"log"
	"strings"
	
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/RainrainWu/gopherbot/db"
	"github.com/RainrainWu/gopherbot/config"
)

type Update struct {
	
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	
	Text     string   `json:"text"`
	Chat     Chat     `json:"chat"`
}

type Chat struct {
	
	Id int `json:"id"`
}

var (

	bot *tgbotapi.BotAPI
	err interface{}
)

func init() {

	authorizeBot(config.Config.TGConfig.Token)
}

func authorizeBot(token string) {

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func PollingBot() {
	
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		command := strings.Split(update.Message.Text, " ")
		if command[0] != "/gopher" {
			continue
		}

		var text string
		switch command[1] {
		case "res":
			text = resourceHandler(command[2:])
		case "tag":
			text = tagHandler(command[2:])
		default:
			text = "help message"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		// msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func resourceHandler(command []string) string {

	switch command[0] {
	case "ls", "list":
		var list []db.Resource
		if len(command) == 1 {
			list = db.ListResources()
		} else if len(command) == 2 {
			list = db.QueryResources(command[1])
		} else {
			return "Too many parameters"
		}

		reply := ""
		for _, i := range list {
			tpl := "%s\t%s"
			reply += fmt.Sprintf(tpl, i.Name, i.Url) + "\n"
		}
		if reply == "" {
			reply = "No resources found."
		}
		return reply

	case "new":
		db.CreateResource(command[1], command[2])
		tpl := "Create resource %s with url %s."
		return fmt.Sprintf(tpl, command[1], command[2])

	case "del":
		db.DeleteResource(command[1])
		tpl := "Delete resource %s."
		return fmt.Sprintf(tpl, command[1])

	case "tag":
		db.RegisterResource(command[1], command[2])
		tpl := "Tag resource %s with tag %s."
		return fmt.Sprintf(tpl, command[1], command[2])

	case "detag":
		db.DeregisterResource(command[1], command[2])
		tpl := "Detag resource %s with tag %s."
		return fmt.Sprintf(tpl, command[1], command[2])
	
	default:
		return db.GetResource(command[0]).Url
	}
}

func tagHandler(command []string) string {

	switch command[0] {
	default:
		var list []db.Team
		if len(command) == 1 {
			list = db.ListTeams()
		} else {
			return "Too many parameters"
		}

		reply := ""
		for _, i := range list {
			reply += i.Name + "\n"
		}
		if reply == "" {
			reply = "No teams found."
		}
		return reply
	}
}