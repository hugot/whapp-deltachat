package main

import (
	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
)

type BridgeContext struct {
	WhappConn    *whatsapp.Conn
	DCContext    *deltachat.Context
	DB           *Database
	DCUserID     uint32
	DCUserChatID uint32
}

// Find or create a deltachat verified group chat for a whatsapp JID and return it's ID.
func (b *BridgeContext) GetOrCreateDCIDForJID(JID string, isGroup bool) (uint32, error) {
	if DCID, _ := b.DB.GetDCIDForWhappJID(JID); DCID != nil {
		return *DCID, nil
	}

	chatName := JID
	if isGroup {
		chat, ok := b.WhappConn.Store.Chats[JID]

		if ok {
			chatName = chat.Name
		}
	} else {
		contact, ok := b.WhappConn.Store.Contacts[JID]

		if ok {
			chatName = contact.Name
		}
	}

	DCID := b.DCContext.CreateGroupChat(true, chatName)

	err := b.DB.StoreDCIDForJID(JID, DCID)

	if err != nil {
		return DCID, err
	}

	b.DCContext.AddContactToChat(DCID, b.DCUserID)

	return DCID, err
}

func (b *BridgeContext) SendLog(logString string) {
	b.DCContext.SendTextMessage(b.DCUserChatID, logString)
}
