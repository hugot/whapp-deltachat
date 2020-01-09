package main

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

	sessionFile := storageDir + "/whapp-session.json"
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		session, err = WhappQrLogin(storageDir, ctx)

		if err != nil {
			return err
		}
	} else {
		session = whatsapp.Session{}

		sessionJson, err := ioutil.ReadFile(sessionFile)

		err = json.Unmarshal(sessionJson, &session)

		if err != nil {
			return err
		}

		session, err = wac.RestoreWithSession(session)

		if err != nil {
			return err
		}
	}

	err = StoreWhappSession(session, storageDir)

	return err
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
			log.Fatal(err)
		}

		err = qrcode.WriteFile(qrCode, qrcode.Medium, 256, tmpFile.Name())

		if err != nil {
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
