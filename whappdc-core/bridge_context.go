package core

import (
	"log"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
)

type BridgeContext struct {
	Config         *Config
	WhappConn      *whatsapp.Conn
	DCContext      *deltachat.Context
	DB             *Database
	MessageTracker *MessageTracker
	DCUserID       uint32
	DCUserChatID   uint32
}

// Find or create a deltachat verified group chat for a whatsapp JID and return it's ID.
func (b *BridgeContext) GetOrCreateDCIDForJID(JID string) (uint32, error) {
	if DCID, _ := b.DB.GetDCIDForWhappJID(JID); DCID != nil {
		return *DCID, nil
	}

	chatName := JID
	chat, ok := b.WhappConn.Store.Chats[JID]

	if ok {
		chatName = chat.Name
	} else if sender, ok := b.WhappConn.Store.Contacts[JID]; ok {
		chatName = sender.Name
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

func (b *BridgeContext) MessageWasSent(ID string) bool {
	sent, err := b.MessageTracker.WasSent(ID)

	if err != nil {
		log.Println(err)
		b.SendLog(err.Error())
	}

	return sent
}

func (b *BridgeContext) ShouldMessageBeSent(info whatsapp.MessageInfo) bool {
	// Skip if the message has already been sent
	if b.MessageWasSent(info.Id) {
		return false
	}

	// send if not from user
	if !info.FromMe {
		return true
	}

	// If from user, only send when it is enabled in the config
	return b.Config.App.ShowFromMe
}
