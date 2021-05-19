package db

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/syndtr/goleveldb/leveldb"
)

func PutMsg(db interface{}, key []byte, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	switch db.(type) {
	case *leveldb.DB:
		return db.(*leveldb.DB).Put(key, data, nil)
	case *leveldb.Batch:
		db.(*leveldb.Batch).Put(key, data)
		return nil
	}
	return fmt.Errorf("db.PutMsg(): Unrecognized input type")
}

func GetMsg(db *leveldb.DB, key []byte, msg proto.Message) error {
	data, err := db.Get(key, nil)
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, msg)
}
