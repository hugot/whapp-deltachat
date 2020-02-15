package botcommands

import (
	"fmt"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
	"github.com/hugot/whapp-deltachat/core"
)

func NewWhappBridge(bridgeContext *core.BridgeContext) *WhappBridge {
	return &WhappBridge{
		bridgeCtx: bridgeContext,
	}
}

type WhappBridge struct {
	bridgeCtx *core.BridgeContext
}

func (b *WhappBridge) Accepts(c *deltachat.Chat, m *deltachat.Message) bool {
	chatID := c.GetID()

	chatJID, err := b.bridgeCtx.DB.GetWhappJIDForDCID(chatID)

	if err != nil {
		// The database is failing, very much an edge case.
		b.bridgeCtx.SendLog(err.Error())

		return false
	}

	// Only forward messages for known groups,
	// Don't forward info messages like "group name changed" etc.
	return chatJID != nil && !m.IsInfo()
}

func (b *WhappBridge) Execute(
	c *deltachat.Context,
	chat *deltachat.Chat,
	m *deltachat.Message,
) {
	JID, err := b.bridgeCtx.DB.GetWhappJIDForDCID(chat.GetID())

	if err != nil {
		b.bridgeCtx.SendLog(
			fmt.Sprintf(
				"Database error in Whapp bridge: %s",
				err.Error(),
			),
		)

		return
	}

	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: *JID,
		},
		Text: m.GetText(),
	}

	_, err = b.bridgeCtx.WhappConn.Send(text)

	if err != nil {
		b.bridgeCtx.SendLog(
			fmt.Sprintf(
				"Error sending message to %s. \nMessage contents: %s\nError: %s",
				*JID,
				m.GetText(),
				err.Error(),
			),
		)
	}
}
