package main

import (
	"fmt"
	"log"

	"github.com/Rhymen/go-whatsapp"
)

type MessageHandler struct {
	Action MessageAction
	Jid    string
}

type MessageAction func() error

func MakeTextMessageAction(b *BridgeContext, m whatsapp.TextMessage) MessageAction {
	return func() error {
		JID := m.Info.RemoteJid

		DCID, err := b.GetOrCreateDCIDForJID(JID, m.Info.RemoteJid != m.Info.SenderJid)

		if err != nil {
			log.Println(err)
			b.SendLog(err.Error())
		}

		senderName := m.Info.Source.GetParticipant()
		contact, ok := b.WhappConn.Store.Contacts[senderName]
		if ok {
			senderName = contact.Name
		}

		b.DCContext.SendTextMessage(
			DCID,
			fmt.Sprintf("%s:\n%s", senderName, m.Text),
		)

		return nil
	}
}
