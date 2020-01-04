package main

import (
	"fmt"
	"log"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
)

type WhappHandler struct {
	dcContext *deltachat.Context
	db        *Database
	dcUserID  uint32
	wac       *whatsapp.Conn
}

// Find or create a deltachat verified group chat for a whatsapp JID and return it's ID.
func (h *WhappHandler) getOrCreateDCIDForJID(JID string, isGroup bool) (uint32, error) {
	if DCID, _ := h.db.GetDCIDForWhappJID(JID); DCID != nil {
		return *DCID, nil
	}

	chatName := JID
	if isGroup {
		chat, ok := h.wac.Store.Chats[JID]

		if ok {
			chatName = chat.Name
		}
	} else {
		contact, ok := h.wac.Store.Contacts[JID]

		if ok {
			chatName = contact.Name
		}
	}

	DCID := h.dcContext.CreateGroupChat(true, chatName)

	err := h.db.StoreDCIDForJID(JID, DCID)

	if err != nil {
		return DCID, err
	}

	h.dcContext.AddContactToChat(DCID, h.dcUserID)

	return DCID, err
}

func (h *WhappHandler) HandleError(err error) {
	log.Println("Whatsapp Error: " + err.Error())
}

func (h *WhappHandler) HandleTextMessage(m whatsapp.TextMessage) {
	JID := m.Info.RemoteJid

	DCID, err := h.getOrCreateDCIDForJID(JID, m.Info.RemoteJid != m.Info.SenderJid)

	if err != nil {
		log.Println(err)

		chatID := h.dcContext.GetChatIDByContactID(h.dcUserID)
		h.dcContext.SendTextMessage(
			chatID,
			err.Error(),
		)
	}

	senderName := m.Info.SenderJid
	contact, ok := h.wac.Store.Contacts[m.Info.SenderJid]
	if ok {
		senderName = contact.Name
	}

	h.dcContext.SendTextMessage(
		DCID,
		fmt.Sprintf("%s:\n%s", senderName, m.Text),
	)
}
