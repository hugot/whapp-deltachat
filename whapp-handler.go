package main

import (
	"log"

	"github.com/Rhymen/go-whatsapp"
)

type WhappHandler struct {
	BridgeContext *BridgeContext
	MessageWorker *MessageWorker
}

func (h *WhappHandler) HandleError(err error) {
	log.Println("Whatsapp Error: " + err.Error())
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
