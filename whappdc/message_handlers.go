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

func MakeTextMessageAction(w *WhappContext, m whatsapp.TextMessage) MessageAction {
	return func() error {
		if !w.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := w.GetOrCreateDCIDForJID(JID)

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(w.BridgeCtx, m.Info)

		w.DCCtx().SendTextMessage(
			DCID,
			fmt.Sprintf("%s:\n%s", senderName, m.Text),
		)

		return w.MessageTracker.MarkSent(&m.Info.Id)
	}
}

func MakeImageMessageAction(w *WhappContext, m whatsapp.ImageMessage) MessageAction {
	return func() error {
		if !w.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := w.GetOrCreateDCIDForJID(JID)

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(w.BridgeCtx, m.Info)

		imageData, err := m.Download()

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			w.markSentIf404Download(&m.Info.Id, err)

			return err
		}

		filename, err := WriteTempFile(w.BridgeCtx, imageData, "img")

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		message := w.DCCtx().NewMessage(deltachat.DC_MSG_IMAGE)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:\n%s", senderName, m.Caption))
		message.SetFile(filename, m.Type)

		w.DCCtx().SendMessage(DCID, message)

		return w.MessageTracker.MarkSent(&m.Info.Id)
	}
}

func MakeDocumentMessageAction(w *WhappContext, m whatsapp.DocumentMessage) MessageAction {
	return func() error {
		if !w.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := w.GetOrCreateDCIDForJID(JID)

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(w.BridgeCtx, m.Info)

		documentData, err := m.Download()

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			w.markSentIf404Download(&m.Info.Id, err)

			return err
		}

		filename, err := WriteTempFile(w.BridgeCtx, documentData, "doc")

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		message := w.DCCtx().NewMessage(deltachat.DC_MSG_FILE)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:\n%s", senderName, m.Title))
		message.SetFile(filename, m.Type)

		w.DCCtx().SendMessage(DCID, message)

		return w.MessageTracker.MarkSent(&m.Info.Id)
	}
}

func MakeAudioMessageAction(w *WhappContext, m whatsapp.AudioMessage) MessageAction {
	return func() error {
		if !w.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := w.GetOrCreateDCIDForJID(JID)

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(w.BridgeCtx, m.Info)

		audioData, err := m.Download()

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			w.markSentIf404Download(&m.Info.Id, err)

			return err
		}

		filename, err := WriteTempFile(w.BridgeCtx, audioData, "audio")

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		message := w.DCCtx().NewMessage(deltachat.DC_MSG_AUDIO)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:", senderName))
		message.SetFile(filename, m.Type)

		w.DCCtx().SendMessage(DCID, message)

		return w.MessageTracker.MarkSent(&m.Info.Id)
	}
}

func MakeVideoMessageAction(w *WhappContext, m whatsapp.VideoMessage) MessageAction {
	return func() error {
		if !w.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := w.GetOrCreateDCIDForJID(JID)

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(w.BridgeCtx, m.Info)

		videoData, err := m.Download()

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			w.markSentIf404Download(&m.Info.Id, err)

			return err
		}

		filename, err := WriteTempFile(w.BridgeCtx, videoData, "vid")

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		message := w.DCCtx().NewMessage(deltachat.DC_MSG_VIDEO)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:", senderName))
		message.SetFile(filename, m.Type)

		w.DCCtx().SendMessage(DCID, message)

		return w.MessageTracker.MarkSent(&m.Info.Id)
	}
}

// 2020-01-14 10:34 TODO: Find out why this doesn't work.
func MakeContactMessageAction(w *WhappContext, m whatsapp.ContactMessage) MessageAction {
	return func() error {
		if !w.ShouldMessageBeSent(m.Info) {
			return nil
		}

		JID := m.Info.RemoteJid
		DCID, err := w.GetOrCreateDCIDForJID(JID)

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		senderName := DetermineSenderName(w.BridgeCtx, m.Info)

		filename, err := WriteTempFile(w.BridgeCtx, []byte(m.Vcard), "vcf")

		if err != nil {
			w.BridgeCtx.SendLog(err.Error())
			return err
		}

		message := w.DCCtx().NewMessage(deltachat.DC_MSG_FILE)
		defer message.Unref()
		message.SetText(fmt.Sprintf("%s:", senderName))
		message.SetFile(filename, mime.TypeByExtension(".vcf"))

		w.DCCtx().SendMessage(DCID, message)

		return w.MessageTracker.MarkSent(&m.Info.Id)
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
