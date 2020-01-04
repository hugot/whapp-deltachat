package main

import (
	"encoding/binary"

	"go.etcd.io/bbolt"
)

type Database struct {
	db     *bbolt.DB
	dbPath string
}

const (
	JID_TO_DCID_INT uint8 = iota
	DCID_TO_JID_INT
	KEY_VALUE_INT
)

var (
	JID_TO_DCID = []byte{JID_TO_DCID_INT}
	DCID_TO_JID = []byte{DCID_TO_JID_INT}
	KEY_VALUE   = []byte{KEY_VALUE_INT}
)

func (d *Database) Init() error {
	db, err := bbolt.Open(d.dbPath, 0600, nil)

	if err != nil {
		return err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(JID_TO_DCID)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(DCID_TO_JID)

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(KEY_VALUE)

		return err
	})

	d.db = db

	return nil
}

func (d *Database) GetDCIDForWhappJID(JID string) (*uint32, error) {
	var DCID *uint32

	err := d.db.View(func(tx *bbolt.Tx) error {
		rawDCID := tx.Bucket(JID_TO_DCID).Get([]byte(JID))

		if rawDCID == nil {
			DCID = nil
		} else {
			i := binary.LittleEndian.Uint32(rawDCID)
			DCID = &i
		}

		return nil
	})

	return DCID, err
}

func (d *Database) GetWhappJIDForDCID(DCID uint32) (*string, error) {
	var JID *string

	rawDCID := make([]byte, 4)
	binary.LittleEndian.PutUint32(rawDCID, DCID)

	err := d.db.View(func(tx *bbolt.Tx) error {

		rawJID := tx.Bucket(DCID_TO_JID).Get(rawDCID)

		if rawJID == nil {
			JID = nil
		} else {
			str := string(rawJID)
			JID = &str
		}

		return nil
	})

	return JID, err
}

func (d *Database) StoreDCIDForJID(JID string, DCID uint32) error {
	err := d.db.Update(func(tx *bbolt.Tx) error {

		DCIDbs := make([]byte, 4)

		binary.LittleEndian.PutUint32(DCIDbs, DCID)

		err := tx.Bucket(JID_TO_DCID).Put([]byte(JID), DCIDbs)

		if err != nil {
			return err
		}

		err = tx.Bucket(DCID_TO_JID).Put(DCIDbs, []byte(JID))

		return err
	})

	return err
}

func (d *Database) Put(key []byte, value []byte) error {
	err := d.db.Update(func(tx *bbolt.Tx) error {
		err := tx.Bucket(KEY_VALUE).Put(key, value)

		return err
	})

	return err
}

func (d *Database) Get(key []byte) []byte {
	var value []byte

	d.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket(KEY_VALUE).Get(key)

		return nil
	})

	return value
}
