package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/hugot/go-deltachat/deltachat"
	"github.com/skip2/go-qrcode"
)

func CreateAndLoginWhappConnection(
	storageDir string,
	ctx *BridgeContext,
) error {
	wac, err := whatsapp.NewConn(30 * time.Second)

	if err != nil {
		return err
	}

	ctx.WhappConn = wac

	var session whatsapp.Session

	sessionFile := WhappSessionFileName(storageDir)
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		session, err = WhappQrLogin(storageDir, ctx)

		if err != nil {
			return err
		}

		return StoreWhappSession(session, storageDir)
	}

	return RestoreWhappSessionFromStorage(storageDir, wac)
}

func RestoreWhappSessionFromStorage(storageDir string, wac *whatsapp.Conn) error {
	storedSession, err := GetStoredWhappSession(storageDir)
	if err != nil {
		return err
	}

	session, err := wac.RestoreWithSession(*storedSession)
	if err != nil {
		return err
	}

	return StoreWhappSession(session, storageDir)
}

func WhappSessionFileName(storageDir string) string {
	return storageDir + "/whapp-session.json"
}

func GetStoredWhappSession(storageDir string) (*whatsapp.Session, error) {
	session := &whatsapp.Session{}

	sessionJson, err := ioutil.ReadFile(WhappSessionFileName(storageDir))

	err = json.Unmarshal(sessionJson, session)

	return session, err
}

func WhappQrLogin(
	storageDir string,
	ctx *BridgeContext,
) (whatsapp.Session, error) {
	qrChan := make(chan string)

	go func() {
		qrCode := <-qrChan

		tmpFile, err := ioutil.TempFile(storageDir+"/tmp", "XXXXXXX-qr")

		if err != nil {
			ctx.SendLog("Failed to create temporarary file: " + err.Error())
			log.Fatal(err)
		}

		err = qrcode.WriteFile(qrCode, qrcode.Medium, 256, tmpFile.Name())

		if err != nil {
			ctx.SendLog("Failed to save qrcode file: " + err.Error())
			log.Fatal(err)
		}

		message := ctx.DCContext.NewMessage(deltachat.DC_MSG_IMAGE)
		defer message.Unref()

		message.SetFile(tmpFile.Name(), mime.TypeByExtension(".png"))

		message.SetText("Scan this QR code from whatsapp")

		ctx.DCContext.SendMessage(
			ctx.DCUserChatID,
			message,
		)
	}()

	session, err := ctx.WhappConn.Login(qrChan)

	return session, err
}

func StoreWhappSession(session whatsapp.Session, storageDir string) error {
	sessionJson, err := json.Marshal(session)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(storageDir+"/whapp-session.json", sessionJson, 0600)

	return err
}
