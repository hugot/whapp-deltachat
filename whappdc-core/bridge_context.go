package core

import (
	"log"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltabot"
	"github.com/hugot/go-deltachat/deltachat"
)

type BridgeContext struct {
	Config       *Config
	WhappConn    *whatsapp.Conn
	DCContext    *deltachat.Context
	DCClient     *deltachat.Client
	DB           *Database
	DCUserID     uint32
	DCUserChatID uint32
}

func NewBridgeContext(config *Config) *BridgeContext {
	db := NewDatabase(config.App.DataFolder + "/app.db")

	return &BridgeContext{
		Config: config,
		DB:     db,
	}
}

// Do all the initialization stuff like intializing databases & services, configuring
// clients, connecting to remote servers, logging in etc.
func (b *BridgeContext) Init(
	whappHandler whatsapp.Handler,
	botCommands []deltabot.Command,
) error {
	err := b.DB.Init()

	if err != nil {
		return err
	}

	dcClient, err := BootstrapDcClientFromConfig(*b.Config, b)

	b.SendLog("Whapp-Deltachat started.")

	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		log.Println("Attempting whapp login")
		err = CreateAndLoginWhappConnection(b.Config.App.DataFolder, b)

		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}

	b.WhappConn.AddHandler(whappHandler)

	bot := &deltabot.Bot{}

	for _, command := range botCommands {
		bot.AddCommand(command)
	}

	dcClient.On(deltachat.DC_EVENT_INCOMING_MSG, bot.HandleMessage)

	return nil
}

func (b *BridgeContext) Close() error {
	_, err := b.WhappConn.Disconnect()

	if err != nil {
		return err
	}

	b.DCClient.Close()

	err = b.DB.Close()

	return err
}

func (b *BridgeContext) SendLog(logString string) {
	b.DCContext.SendTextMessage(b.DCUserChatID, logString)
}
