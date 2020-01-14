package whappdc

import (
	"fmt"
	"log"

	"github.com/Rhymen/go-whatsapp"
	core "github.com/hugot/whapp-deltachat/whappdc-core"
)

type WhappHandler struct {
	BridgeContext *core.BridgeContext
	MessageWorker *MessageWorker
}

func (h *WhappHandler) HandleError(err error) {
	// Err might be nil.
	if err == nil {
		return
	}

	// there is a weird edge case in which calling err.Error() causes a nil pointer
	// dereference in this function. This is probably a bug in rhymen/go-whatsapp. Let's
	// keep this here to recover from panics.
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in *WhappHandler.HandleError")
		}
	}()

	// If connection to the whapp servers failed for some reason, just retry.
	if _, connectionFailed := err.(*whatsapp.ErrConnectionFailed); connectionFailed {
		err = core.RestoreWhappSessionFromStorage(
			h.BridgeContext.Config.App.DataFolder,
			h.BridgeContext.WhappConn,
		)

		if err != nil {
			logString := "Failed to restore whatsapp connection: " + err.Error()
			log.Println(logString)
			h.BridgeContext.SendLog(logString)
			return
		}
	}

	typeLogString := fmt.Sprintf("Whatsapp Error of type: %T", err)
	log.Println(typeLogString)

	// Calling err.Error() here may cause a nil pointer dereference panic. See defer
	// statement above.
	logString := "Whatsapp Error: " + err.Error()
	log.Println(logString)

	// Invalid ws data seems to be pretty common, let's not bore the user with that.xg
	if err.Error() != "error processing data: "+whatsapp.ErrInvalidWsData.Error() {
		h.BridgeContext.SendLog(logString)
	}
}

func (h *WhappHandler) HandleTextMessage(m whatsapp.TextMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeTextMessageAction(h.BridgeContext, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleImageMessage(m whatsapp.ImageMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeImageMessageAction(h.BridgeContext, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleDocumentMessage(m whatsapp.DocumentMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeDocumentMessageAction(h.BridgeContext, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleAudioMessage(m whatsapp.AudioMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeAudioMessageAction(h.BridgeContext, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleVideoMessage(m whatsapp.VideoMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeVideoMessageAction(h.BridgeContext, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleContactMessage(m whatsapp.VideoMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeVideoMessageAction(h.BridgeContext, m),
	}

	h.MessageWorker.HandleMessage(handler)
}
