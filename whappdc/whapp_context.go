package whappdc

import (
	"log"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
	core "github.com/hugot/whapp-deltachat/whappdc-core"
)

func NewWhappContext(
	bridgeCtx *core.BridgeContext,
	msgTrackerFlushInterval time.Duration,
) *WhappContext {
	messageTracker := NewMessageTracker(bridgeCtx.DB, msgTrackerFlushInterval)

	return &WhappContext{
		BridgeCtx:      bridgeCtx,
		MessageTracker: messageTracker,
	}
}

type WhappContext struct {
	BridgeCtx      *core.BridgeContext
	MessageTracker *MessageTracker
}

// Flushes the message tracker.  While this method is called "Close", it currently doesn't
// make the context unusable. It might do so in the future though so no reason to change
// the name atm.
func (w *WhappContext) Close() error {
	return w.MessageTracker.Flush()
}

// Find or create a deltachat verified group chat for a whatsapp JID and return it's ID.
func (w *WhappContext) GetOrCreateDCIDForJID(JID string) (uint32, error) {
	if DCID, _ := w.BridgeCtx.DB.GetDCIDForWhappJID(JID); DCID != nil {
		return *DCID, nil
	}

	chatName := JID
	chat, ok := w.BridgeCtx.WhappConn.Store.Chats[JID]

	if ok {
		chatName = chat.Name
	} else if sender, ok := w.BridgeCtx.WhappConn.Store.Contacts[JID]; ok {
		chatName = sender.Name
	}

	DCID := w.BridgeCtx.DCContext.CreateGroupChat(true, chatName)

	err := w.BridgeCtx.DB.StoreDCIDForJID(JID, DCID)

	if err != nil {
		return DCID, err
	}

	w.BridgeCtx.DCContext.AddContactToChat(DCID, w.BridgeCtx.DCUserID)

	return DCID, err
}

func (w *WhappContext) MessageWasSent(ID string) bool {
	sent, err := w.MessageTracker.WasSent(ID)

	if err != nil {
		log.Println(err)
		w.BridgeCtx.SendLog(err.Error())
	}

	return sent
}

func (w *WhappContext) ShouldMessageBeSent(info whatsapp.MessageInfo) bool {
	// Skip if the message has already been sent
	if w.MessageWasSent(info.Id) {
		return false
	}

	// send if not from user
	if !info.FromMe {
		return true
	}

	// If from user, only send when it is enabled in the config
	return w.BridgeCtx.Config.App.ShowFromMe
}

// Alias for easy access of DC Context
func (w *WhappContext) DCCtx() *deltachat.Context {
	return w.BridgeCtx.DCContext
}

// Helper to easily mark failed media downloads as sent anyways.
func (w *WhappContext) markSentIf404Download(ID *string, err error) {
	if err == whatsapp.ErrMediaDownloadFailedWith404 {
		err := w.MessageTracker.MarkSent(ID)

		if err != nil {
			log.Println(err)
			w.BridgeCtx.SendLog(err.Error())
		}
	}
}
