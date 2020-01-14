package whappdc

import (
	"fmt"
	"io/ioutil"
	"mime"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
	core "github.com/hugot/whapp-deltachat/whappdc-core"
)

type MessageHandler struct {
	Action MessageAction
	Jid    string
}

type MessageAction func() error

func MakeTextMessageAction(b *core.BridgeContext, m whatsapp.TextMessage) MessageAction {
	return func() error {
		if !b.ShouldMessageBeSent(m.Info) {
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

func MakeImageMessageAction(b *core.BridgeContext, m whatsapp.ImageMessage) MessageAction {
	return func() error {
		if !b.ShouldMessageBeSent(m.Info) {
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

		filename, err := WriteTempFile(b, imageData, "img")

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		message := b.DCContext.NewMessage(deltachat.DC_MSG_IMAGE)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:\n%s", senderName, m.Caption))
		message.SetFile(filename, m.Type)

		b.DCContext.SendMessage(DCID, message)

		return b.MessageTracker.MarkSent(&m.Info.Id)
	}
}

func MakeDocumentMessageAction(b *core.BridgeContext, m whatsapp.DocumentMessage) MessageAction {
	return func() error {
		if !b.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := b.GetOrCreateDCIDForJID(JID)

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(b, m.Info)

		documentData, err := m.Download()

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		filename, err := WriteTempFile(b, documentData, "doc")

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		message := b.DCContext.NewMessage(deltachat.DC_MSG_FILE)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:\n%s", senderName, m.Title))
		message.SetFile(filename, m.Type)

		b.DCContext.SendMessage(DCID, message)

		return b.MessageTracker.MarkSent(&m.Info.Id)
	}
}

func MakeAudioMessageAction(b *core.BridgeContext, m whatsapp.AudioMessage) MessageAction {
	return func() error {
		if !b.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := b.GetOrCreateDCIDForJID(JID)

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(b, m.Info)

		audioData, err := m.Download()

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		filename, err := WriteTempFile(b, audioData, "audio")

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		message := b.DCContext.NewMessage(deltachat.DC_MSG_AUDIO)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:", senderName))
		message.SetFile(filename, m.Type)

		b.DCContext.SendMessage(DCID, message)

		return b.MessageTracker.MarkSent(&m.Info.Id)
	}
}

func MakeVideoMessageAction(b *core.BridgeContext, m whatsapp.VideoMessage) MessageAction {
	return func() error {
		if !b.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := b.GetOrCreateDCIDForJID(JID)

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(b, m.Info)

		videoData, err := m.Download()

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		filename, err := WriteTempFile(b, videoData, "vid")

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		message := b.DCContext.NewMessage(deltachat.DC_MSG_VIDEO)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:", senderName))
		message.SetFile(filename, m.Type)

		b.DCContext.SendMessage(DCID, message)

		return b.MessageTracker.MarkSent(&m.Info.Id)
	}
}

// 2020-01-14 10:34 TODO: Find out why this doesn't work.
func MakeContactMessageAction(b *core.BridgeContext, m whatsapp.ContactMessage) MessageAction {
	return func() error {
		if !b.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := b.GetOrCreateDCIDForJID(JID)

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(b, m.Info)

		filename, err := WriteTempFile(b, []byte(m.Vcard), "vcf")

		if err != nil {
			b.SendLog(err.Error())
			return err
		}

		message := b.DCContext.NewMessage(deltachat.DC_MSG_FILE)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:", senderName))
		message.SetFile(filename, mime.TypeByExtension(".vcf"))

		b.DCContext.SendMessage(DCID, message)

		return b.MessageTracker.MarkSent(&m.Info.Id)

	}
}

////
// Helpers

func DetermineSenderName(b *core.BridgeContext, info whatsapp.MessageInfo) string {
	senderName := info.Source.GetParticipant()

	// No participant probably means that this isn't a group chat.
	if senderName == "" {
		senderName = info.RemoteJid
	}

	if info.FromMe == true {
		dcContact := b.DCContext.GetContact(b.DCUserID)
		defer dcContact.Unref()

		return dcContact.GetDisplayName()
	}

	contact, ok := b.WhappConn.Store.Contacts[senderName]
	if ok {
		senderName = contact.Name
	}

	return senderName
}

func WriteTempFile(b *core.BridgeContext, data []byte, template string) (string, error) {
	tmpFile, err := ioutil.TempFile(
		b.Config.App.DataFolder+"/tmp",
		template,
	)

	if err != nil {
		b.SendLog(err.Error())
		return "", err
	}

	return tmpFile.Name(), ioutil.WriteFile(tmpFile.Name(), data, 0600)
}
