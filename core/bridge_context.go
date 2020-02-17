package core

import (
	"fmt"
	"log"
	"os"

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
	logger       deltachat.Logger
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

	logFile, err := os.OpenFile(
		b.Config.App.DataFolder+"/whapp-deltachat.log",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0644,
	)

	b.logger = log.New(logFile, "", log.LstdFlags)

	fmt.Printf("Logs will be written to %s\n", logFile.Name())

	dcClient, err := BootstrapDcClientFromConfig(*b.Config, b)

	b.SendLog("Whapp-Deltachat started.")

	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		b.SendLog(fmt.Sprintf("Attempting whapp login (attempt %d)", i+1))
		err = CreateAndLoginWhappConnection(b.Config.App.DataFolder, b)

		if err == nil {
			b.SendLog("Whapp login was successful")
			break
		}
	}

	if err != nil {
		return err
	}

	b.WhappConn.AddHandler(whappHandler)

	bot := deltabot.NewBot(b.logger)

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

func (b *BridgeContext) Logger() deltachat.Logger {
	return b.logger
}

func (b *BridgeContext) SendLog(logString string) {
	b.logger.Println(logString)
	b.DCContext.SendTextMessage(b.DCUserChatID, logString)
}

// Returns true when a DC message is eligible to be bridged.
func (b *BridgeContext) Accepts(c *deltachat.Chat, m *deltachat.Message) bool {
	chatID := c.GetID()

	chatJID, err := b.DB.GetWhappJIDForDCID(chatID)

	if err != nil {
		// The database is failing, very much an edge case.
		b.SendLog(err.Error())

		return false
	}

	// Only forward messages for known groups,
	// Don't forward info messages like "group name changed" etc.
	return chatJID != nil && !m.IsInfo()
}
