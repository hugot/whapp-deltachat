package botcommands

import (
	"github.com/hugot/go-deltachat/deltachat"
)

type Echo struct{}

const echoPrefix = "!echo"

func (e *Echo) Accepts(c *deltachat.Chat, m *deltachat.Message) bool {
	messageText := m.GetText()

	return len(messageText) > len(echoPrefix) && messageText[0:len(echoPrefix)] == echoPrefix
}

func (e *Echo) Execute(c *deltachat.Context, chat *deltachat.Chat, message *deltachat.Message) {
	chatID := chat.GetID()

	messageText := message.GetText()

	c.SendTextMessage(chatID, messageText[len(echoPrefix)+1:])
}
