package grapho

import (
	"encoding/hex"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/wallerprogramm/btc-csv/neo4jcsv"
	"github.com/wallerprogramm/btc-csv/timelog"

	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/serial"
)

var W Writer = neo4jcsv.NewWriter()
var T = timelog.NewWriter()

// Create a graph in neo4j containing all the transactions and their relations
func Build(x interface{}) (func(b *models.Block) error, error) {
	var startHeight uint32
	var err error

	db, err := leveldb.OpenFile(os.Getenv("DBDIR"), nil)
	if err != nil {
		log.Fatalln("Cannot open utxo store", err)
	}

	return func(b *models.Block) error {
		if b == nil {
			log.Info("All files read")
			return nil
		}
		if b.NHeight < startHeight {
			return fmt.Errorf("Jump to height %d", startHeight)
		}
		//		c.Height = b.NHeight
		decodeBlockHeader(&b.BlockHeader, W)
		blockHash := hex.EncodeToString(serial.ReverseHex(b.BlockHeader.Hash))
		for _, tx := range b.Txs {
			if err = decodeTx(&tx, blockHash, W, db); err != nil {
				log.Warn(fmt.Sprintf("Error: %d %x", b.BlockHeader.NHeight, serial.ReverseHex(tx.Hash)))
				return err
			}
		}
		// Time measure
		T.Log_time(b.NHeight)
		return nil
	}, nil

}
