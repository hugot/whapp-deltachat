package botcommands

import (
	"fmt"
	"log"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
)

type Database interface {
	GetWhappJIDForDCID(DCID uint32) (*string, error)
}

func NewWhappBridge(
	wac *whatsapp.Conn,
	db Database,
	UserChatID uint32,
) *WhappBridge {
	return &WhappBridge{
		wac:        wac,
		db:         db,
		UserChatID: UserChatID,
	}
}

type WhappBridge struct {
	wac        *whatsapp.Conn
	db         Database
	UserChatID uint32
}

func (b *WhappBridge) Accepts(c *deltachat.Chat, m *deltachat.Message) bool {
	chatID := c.GetID()

	chatJID, err := b.db.GetWhappJIDForDCID(chatID)

	if err != nil {
		// The database is failing, time to die :(
		log.Fatal(err)
	}

	return chatJID != nil
}

func (b *WhappBridge) Execute(
	c *deltachat.Context,
	chat *deltachat.Chat,
	m *deltachat.Message,
) {
	JID, err := b.db.GetWhappJIDForDCID(chat.GetID())

	if err != nil {
		c.SendTextMessage(
			b.UserChatID,
			fmt.Sprintf(
				"Whapp bridge dying: %s",
				err.Error(),
			),
		)

		log.Fatal(err)
	}

	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: *JID,
		},
		Text: m.GetText(),
	}

	_, err = b.wac.Send(text)

	if err != nil {
		c.SendTextMessage(
			b.UserChatID,
			fmt.Sprintf(
				"Error sending message to %s. Message contents: %s",
				JID,
				m.GetText(),
			),
		)
	}
}
