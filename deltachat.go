package main

import (
	"fmt"
	"log"

	"github.com/hugot/go-deltachat/deltabot"
	"github.com/hugot/go-deltachat/deltachat"
	"github.com/hugot/whapp-deltachat/botcommands"
)

func DcClientFromConfig(databasePath string, config map[string]string) *deltachat.Client {
	client := &deltachat.Client{}

	// Handler for info logs from libdeltachat
	client.On(deltachat.DC_EVENT_INFO, func(c *deltachat.Context, e *deltachat.Event) {
		info, _ := e.Data2.String()

		log.Println(*info)
	})

	client.Open(databasePath)

	for key, value := range config {
		client.SetConfig(key, value)
	}

	// TODO: Make this configurable for users
	client.SetConfig(
		"server_flags",
		fmt.Sprintf(
			"%d",
			deltachat.DC_LP_AUTH_NORMAL|
				deltachat.DC_LP_IMAP_SOCKET_SSL|
				deltachat.DC_LP_SMTP_SOCKET_STARTTLS,
		),
	)

	client.Configure()

	return client
}

func BootstrapDcClientFromConfig(config Config) *deltachat.Client {
	dcClient := DcClientFromConfig(config.App.DataFolder+"/deltachat.db", config.Deltachat)

	context := dcClient.Context()
	userName := "user"

	userID := context.CreateContact(&userName, &config.App.UserAddress)

	context.SendTextMessage(context.CreateChatByContactID(userID), "Whapp-Deltachat initialized")

	bot := &deltabot.Bot{}

	bot.AddCommand(&botcommands.Echo{})

	dcClient.On(deltachat.DC_EVENT_INCOMING_MSG, bot.HandleMessage)

	return dcClient
}
