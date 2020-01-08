package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	db := &Database{
		dbPath: config.App.DataFolder + "/app.db",
	}

	err = db.Init()

	messageTracker := &MessageTracker{
		DB: db,
	}

	defer messageTracker.Flush()

	if err != nil {
		log.Fatal(err)
	}

	bridgeCtx := &BridgeContext{
		Config:         config,
		DB:             db,
		MessageTracker: messageTracker,
	}

	dcClient, err := BootstrapDcClientFromConfig(*config, bridgeCtx)

	bridgeCtx.SendLog("Whapp-Deltachat started.")

	if err != nil {
		log.Fatal(err)
	}

	// Give dc an opportunity to perform some close-down logic
	// and close it's db etc.
	defer dcClient.Close()

	for i := 0; i < 10; i++ {
		err = CreateAndLoginWhappConnection(
			config.App.DataFolder,
			bridgeCtx,
		)

		if err == nil {
			break
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	messageWorker := NewMessageWorker()
	messageWorker.Start()

	bridgeCtx.WhappConn.AddHandler(&WhappHandler{
		BridgeContext: bridgeCtx,
		MessageWorker: messageWorker,
	})

	bot := &deltabot.Bot{}

	bot.AddCommand(&botcommands.Echo{})
	bot.AddCommand(botcommands.NewWhappBridge(
		bridgeCtx.WhappConn, bridgeCtx.DB, bridgeCtx.DCUserChatID,
	))

	dcClient.On(deltachat.DC_EVENT_INCOMING_MSG, bot.HandleMessage)

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)

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
