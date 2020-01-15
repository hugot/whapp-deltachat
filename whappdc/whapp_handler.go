package whappdc

import (
	"fmt"
	"log"
	"time"

	"github.com/Rhymen/go-whatsapp"
	core "github.com/hugot/whapp-deltachat/whappdc-core"
)

// WhappHandler implements go-whaptsapp.Handler
type WhappHandler struct {
	WhappCtx      *WhappContext
	MessageWorker *MessageWorker
}

func NewWhappHandler(
	bridgeCtx *core.BridgeContext,
	msgTrackerFlushInterval time.Duration,
) *WhappHandler {
	return &WhappHandler{
		WhappCtx:      NewWhappContext(bridgeCtx, msgTrackerFlushInterval),
		MessageWorker: NewMessageWorker(),
	}
}

func (h *WhappHandler) Start() {
	h.MessageWorker.Start()
}

func (h *WhappHandler) Stop() error {
	h.MessageWorker.Stop()
	return h.WhappCtx.Close()
}

func (h *WhappHandler) HandleError(err error) {
	// If connection to the whapp servers failed for some reason, just retry.
	if _, connectionFailed := err.(*whatsapp.ErrConnectionFailed); connectionFailed {
		err = core.RestoreWhappSessionFromStorage(
			h.WhappCtx.BridgeCtx.Config.App.DataFolder,
			h.WhappCtx.BridgeCtx.WhappConn,
		)

		if err != nil {
			logString := "Failed to restore whatsapp connection: " + err.Error()
			log.Println(logString)
			h.WhappCtx.BridgeCtx.SendLog(logString)
		}

		return
	}

	typeLogString := fmt.Sprintf("Whatsapp Error of type: %T", err)
	log.Println(typeLogString)

	logString := "Whatsapp Error: " + err.Error()
	log.Println(logString)

	// Invalid ws data seems to be pretty common, let's not bore the user with that.xg
	if err.Error() != "error processing data: "+whatsapp.ErrInvalidWsData.Error() {
		h.WhappCtx.BridgeCtx.SendLog(logString)
	}
}

func (h *WhappHandler) HandleTextMessage(m whatsapp.TextMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeTextMessageAction(h.WhappCtx, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleImageMessage(m whatsapp.ImageMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeImageMessageAction(h.WhappCtx, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleDocumentMessage(m whatsapp.DocumentMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeDocumentMessageAction(h.WhappCtx, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleAudioMessage(m whatsapp.AudioMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeAudioMessageAction(h.WhappCtx, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleVideoMessage(m whatsapp.VideoMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeVideoMessageAction(h.WhappCtx, m),
	}

	h.MessageWorker.HandleMessage(handler)
}

func (h *WhappHandler) HandleContactMessage(m whatsapp.VideoMessage) {
	handler := MessageHandler{
		Jid:    m.Info.RemoteJid,
		Action: MakeVideoMessageAction(h.WhappCtx, m),
	}

	h.MessageWorker.HandleMessage(handler)
}
