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

func NewWhappHandler(bridgeCtx *core.BridgeContext, messageWorker *MessageWorker) *WhappHandler {
	return &WhappHandler{
		BridgeContext: bridgeCtx,
		MessageWorker: messageWorker,
	}
}

func (h *WhappHandler) HandleError(err error) {
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
		}

		return
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
