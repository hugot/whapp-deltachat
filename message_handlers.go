package main

import (
	"fmt"
	"io/ioutil"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
)

type MessageHandler struct {
	Action MessageAction
	Jid    string
}

type MessageAction func() error

func MakeTextMessageAction(b *BridgeContext, m whatsapp.TextMessage) MessageAction {
	return func() error {
		if b.MessageWasSent(m.Info.Id) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := b.GetOrCreateDCIDForJID(JID)

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(b, m.Info)

		b.DCContext.SendTextMessage(
			DCID,
			fmt.Sprintf("%s:\n%s", senderName, m.Text),
		)

		return b.MessageTracker.MarkSent(&m.Info.Id)
	}
}

func MakeImageMessageAction(b *BridgeContext, m whatsapp.ImageMessage) MessageAction {
	return func() error {
		if b.MessageWasSent(m.Info.Id) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := b.GetOrCreateDCIDForJID(JID)

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(b, m.Info)

		imageData, err := m.Download()

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		tmpFile, err := ioutil.TempFile(
			b.Config.App.DataFolder+"/tmp",
			"XXXXXXX-img",
		)

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		err = ioutil.WriteFile(tmpFile.Name(), imageData, 0600)

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		message := b.DCContext.NewMessage(deltachat.DC_MSG_IMAGE)
		message.SetText(fmt.Sprintf("%s:\n%s", senderName, m.Caption))
		message.SetFile(tmpFile.Name(), m.Type)

		b.DCContext.SendMessage(DCID, message)

		return b.MessageTracker.MarkSent(&m.Info.Id)
	}
}

////
// Helpers

func DetermineSenderName(b *BridgeContext, info whatsapp.MessageInfo) string {
	senderName := info.Source.GetParticipant()

	// No participant probably means that this isn't a group chat.
	if senderName == "" {
		senderName = info.RemoteJid
	}

	if info.FromMe == true {
		return b.DCContext.GetContact(b.DCUserID).GetDisplayName()
	}

	contact, ok := b.WhappConn.Store.Contacts[senderName]
	if ok {
		senderName = contact.Name
	}

	return senderName
}
