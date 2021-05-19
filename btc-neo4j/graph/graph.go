package graph

import (
	"encoding/hex"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/serial"
)

// Create a graph in neo4j containing all the transactions and their relations
func Build(x interface{}) (func(b *models.Block) error, error) {
	var startHeight uint32
	var err error

	fmt.Println("something")
	// TODO: Check current block height in neo4j db

	// TODO: Return function that writes to neo4j
	return func(b *models.Block) error {
		//		if b == nil {
		//			return c.DB.Close()
		//		}
		if b.NHeight < startHeight {
			return fmt.Errorf("Jump to height %d", startHeight)
		}
		//		c.Height = b.NHeight
		decodeBlockHeader(&b.BlockHeader)
		blockHash := hex.EncodeToString(serial.ReverseHex(b.BlockHeader.HashBlock))
		for _, tx := range b.Txs {
			if err = decodeTx(&tx, blockHash); err != nil {
				log.Warn(fmt.Sprintf("Error: %d %x", b.BlockHeader.NHeight, serial.ReverseHex(tx.Hash)))
				return err
			}
		}
		return nil
	}, nil

}
