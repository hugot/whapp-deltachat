package botcommands

import (
	"fmt"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
	"github.com/hugot/whapp-deltachat/core"
)

func NewTextMessageProxy(bridgeContext *core.BridgeContext) *TextMessageProxy {
	return &TextMessageProxy{
		bridgeCtx: bridgeContext,
	}
}

type TextMessageProxy struct {
	bridgeCtx *core.BridgeContext
}

func (b *TextMessageProxy) Accepts(c *deltachat.Chat, m *deltachat.Message) bool {
	return b.Accepts(c, m)
}

func (b *TextMessageProxy) Execute(
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
