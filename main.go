package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltabot"
	"github.com/hugot/go-deltachat/deltachat"
	"github.com/hugot/whapp-deltachat/botcommands"
)

func main() {
	argLen := len(os.Args)

	if argLen != 2 {
		fmt.Fprintln(os.Stderr, "Usage: whapp-deltachat CONFIG_FILE")
		return
	}

	configPath := os.Args[1]

	config, err := ConfigFromFile(configPath)

	if err != nil {
		log.Fatal(err)
	}

	ensureDirectoryOrDie(config.App.DataFolder)
	ensureDirectoryOrDie(config.App.DataFolder + "/tmp")

	dcClient := BootstrapDcClientFromConfig(*config)

	// Give dc an opportunity to perform some close-down logic
	// and close it's db etc.
	defer dcClient.Close()

	db := &Database{
		dbPath: config.App.DataFolder + "/app.db",
	}

	err = db.Init()

	if err != nil {
		log.Fatal(err)
	}

	ctx := dcClient.Context()
	userName := "user"
	dcUserID := ctx.CreateContact(&userName, &config.App.UserAddress)

	userChatIDRaw := db.Get([]byte("user-chat-id"))
	var userChatID uint32

	if userChatIDRaw == nil {
		userChatID, err = AddUserAsVerifiedContact(dcUserID, dcClient)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		userChatID = binary.LittleEndian.Uint32(userChatIDRaw)
	}

	userChatIDbs := make([]byte, 4)
	binary.LittleEndian.PutUint32(userChatIDbs, userChatID)
	err = db.Put([]byte("user-chat-id"), userChatIDbs)

	if err != nil {
		log.Fatal(err)
	}

	var wac *whatsapp.Conn

	for i := 0; i < 10; i++ {
		wac, err = CreateAndLoginWhappConnection(
			config.App.DataFolder,
			dcClient.Context(),
			dcUserID,
		)

		if err == nil {
			break
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	wac.AddHandler(&WhappHandler{
		dcContext: dcClient.Context(),
		db:        db,
		dcUserID:  dcUserID,
		wac:       wac,
	})

	bot := &deltabot.Bot{}

	bot.AddCommand(&botcommands.Echo{})
	bot.AddCommand(botcommands.NewWhappBridge(
		wac, db, dcClient.Context().GetChatIDByContactID(dcUserID),
	))

	dcClient.On(deltachat.DC_EVENT_INCOMING_MSG, bot.HandleMessage)

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt)

	for {
		select {
		case sig := <-wait:
			log.Println(sig)
			return
		}
	}

}

func ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}

	}

	err := os.Chmod(dir, 0700)
	if err != nil {
		return err
	}

	return nil
}

func ensureDirectoryOrDie(dir string) {
	err := ensureDirectory(dir)

	if err != nil {
		log.Fatal(err)
	}
}
