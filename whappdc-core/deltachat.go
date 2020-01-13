package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hugot/go-deltachat/deltachat"
	"github.com/mdp/qrterminal"
)

func DcClientFromConfig(databasePath string, config map[string]string) *deltachat.Client {
	client := &deltachat.Client{}

	// Handler for info logs from libdeltachat
	client.On(deltachat.DC_EVENT_INFO, func(c *deltachat.Context, e *deltachat.Event) {
		info, _ := e.Data2.String()

		log.Println(*info)
	})

	client.Open(databasePath)

	for key, value := range config {
		client.SetConfig(key, value)
	}

	// TODO: Make this configurable for users
	client.SetConfig(
		"server_flags",
		fmt.Sprintf(
			"%d",
			deltachat.DC_LP_AUTH_NORMAL|
				deltachat.DC_LP_IMAP_SOCKET_SSL|
				deltachat.DC_LP_SMTP_SOCKET_STARTTLS,
		),
	)

	client.Configure()

	return client
}

// Note: this manipulates the BridgeContext.
func BootstrapDcClientFromConfig(config Config, ctx *BridgeContext) (*deltachat.Client, error) {
	dcClient := DcClientFromConfig(config.App.DataFolder+"/deltachat.db", config.Deltachat)

	DCCtx := dcClient.Context()

	log.Println("Waiting for deltachat client to be configured")
	for !DCCtx.IsConfigured() {
	}

	userName := "user"
	dcUserID := DCCtx.CreateContact(&userName, &config.App.UserAddress)

	// Send a message in a 1:1 chat first, this will let the user's client know that the
	// crypto setup has changed if it has
	DCCtx.SendTextMessage(
		DCCtx.CreateChatByContactID(dcUserID),
		"Hi, Whapp-Deltachat is initializing",
	)

	userChatIDRaw := ctx.DB.Get([]byte("user-chat-id"))
	var (
		userChatID uint32
		err        error
	)

	if userChatIDRaw == nil {
		userChatID, err = AddUserAsVerifiedContact(dcUserID, dcClient)
		if err != nil {
			return nil, err
		}
	} else {
		userChatID = binary.LittleEndian.Uint32(userChatIDRaw)
	}

	userChatIDbs := make([]byte, 4)
	binary.LittleEndian.PutUint32(userChatIDbs, userChatID)
	err = ctx.DB.Put([]byte("user-chat-id"), userChatIDbs)

	ctx.DCUserID = dcUserID
	ctx.DCUserChatID = userChatID
	ctx.DCContext = DCCtx

	return dcClient, err
}

// Add a user as verified contact to the deltachat context. This is necessary to be able
// to create verified groups.
func AddUserAsVerifiedContact(dcUserID uint32, client *deltachat.Client) (uint32, error) {
	confirmChan := make(chan bool)

	client.On(
		deltachat.DC_EVENT_SECUREJOIN_INVITER_PROGRESS,
		func(c *deltachat.Context, e *deltachat.Event) {
			contactIDInt, err := e.Data1.Int()

			if err != nil {
				log.Println(err)

				// Something weird is going on here, deltachat is not passing expected
				// values. Make the join process fail.
				confirmChan <- false
				return
			}

			contactID := uint32(*contactIDInt)

			if contactID != dcUserID {
				log.Println(
					fmt.Sprintf(
						"Unexpected contact ID encountered for securejoin progress: %v",
						contactID,
					),
				)

				return
			}

			progress, err := e.Data2.Int()

			if err != nil {
				log.Println(err)

				confirmChan <- false
				return

			}

			if *progress == 1000 {
				confirmChan <- true
			}
		},
	)

	ctx := client.Context()

	chatID := ctx.CreateGroupChat(true, "Whapp-DC ***status***")

	log.Println("Scan this qr code with your DC client")
	qrterminal.Generate(
		ctx.GetSecurejoinQR(chatID),
		qrterminal.L,
		os.Stdout,
	)

	success := <-confirmChan

	if !success {
		return chatID, errors.New("Contact Verification process failed")
	}

	log.Println("Securejoin verification completed")

	return chatID, nil
}
