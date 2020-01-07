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

		wasSent, err := b.MessageTracker.WasSent(m.Info.Id)

		if err != nil {
			log.Println(err)
			b.SendLog(err.Error())
		}

		// Messgae has already been sent
		if wasSent == true {
			return nil
		}

		DCID, err := b.GetOrCreateDCIDForJID(JID)

		if err != nil {
			log.Println(err)
			b.SendLog(err.Error())
		}

		senderName := m.Info.Source.GetParticipant()

		// No participant probably means that this isn't a group chat.
		if senderName == "" {
			senderName = m.Info.RemoteJid
		}

		if m.Info.FromMe == true {
			senderName = b.DCContext.GetContact(b.DCUserID).GetDisplayName()
		}

		contact, ok := b.WhappConn.Store.Contacts[senderName]
		if ok {
			senderName = contact.Name
		}

		b.DCContext.SendTextMessage(
			DCID,
			fmt.Sprintf("%s:\n%s", senderName, m.Text),
		)

		return b.MessageTracker.MarkSent(&JID)
	}
}
