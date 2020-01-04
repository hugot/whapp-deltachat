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
	dcContext *deltachat.Context,
	dcUserID uint32,
) (*whatsapp.Conn, error) {
	wac, err := whatsapp.NewConn(20 * time.Second)

	if err != nil {
		return wac, err
	}

	var session whatsapp.Session

	sessionFile := storageDir + "/whapp-session.json"
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		session, err = WhappQrLogin(storageDir, wac, dcContext, dcUserID)

		if err != nil {
			return wac, err
		}
	} else {
		session = whatsapp.Session{}

		sessionJson, err := ioutil.ReadFile(sessionFile)

		err = json.Unmarshal(sessionJson, &session)

		if err != nil {
			return wac, err
		}

		session, err = wac.RestoreWithSession(session)

		if err != nil {
			return wac, err
		}
	}

	err = StoreWhappSession(session, storageDir)

	return wac, err
}

func WhappQrLogin(
	storageDir string,
	wac *whatsapp.Conn,
	dcContext *deltachat.Context,
	dcUserID uint32,
) (whatsapp.Session, error) {
	qrChan := make(chan string)

	go func() {
		qrCode := <-qrChan

		tmpFile, err := ioutil.TempFile(storageDir+"/tmp", "qr")

		if err != nil {
			log.Fatal(err)
		}

		err = qrcode.WriteFile(qrCode, qrcode.Medium, 256, tmpFile.Name())

		if err != nil {
			log.Fatal(err)
		}

		message := dcContext.NewMessage(deltachat.DC_MSG_IMAGE)

		log.Println("MIME: " + mime.TypeByExtension("png"))

		message.SetFile(tmpFile.Name(), "image/png")

		message.SetText("Scan this QR code from whatsapp")

		dcContext.SendMessage(
			dcContext.GetChatIDByContactID(dcUserID),
			message,
		)
	}()

	session, err := wac.Login(qrChan)

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
