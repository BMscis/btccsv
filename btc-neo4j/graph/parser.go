package graph

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/wallerprogramm/btc-neo4j/neo4jdb"

	"github.com/chainswatch/bclib/btc"
	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/serial"
)

func decodeInputs(tx *models.Tx, txid string) error {
	//var err error
	// Go through all Inputs and search for address and BTC value
	// in the neo4j db, then add relationship SENDS

	for _, vin := range tx.Vin {
		// TODO: What to do with blockreward inputs?
		if vin.Index == 0xFFFFFFFF { // block reward
			// Is it needed to do something with that?
		} else {
			// TODO: return an error
			neo4jdb.Add_sends_rel(txid, &vin)
		}
	}

	return nil
}

func decodeOutputs(tx *models.Tx, txid string) error {
	var err error
	var addr string

	// Go through all Outputs
	// - add address to neo4j db if it doesn't already exist
	// - add relationship RECEIVES with BTC value and output_nr
	for _, vout := range tx.Vout {
		if addr, err = btc.DecodeAddr(vout.AddrType, vout.Addr); err != nil {
			log.Fatal(fmt.Sprintf("DecodeOutputs: %s: %d %x", err.Error(), serial.ReverseHex(tx.Hash)))
		}
		neo4jdb.Add_address(addr)
		neo4jdb.Add_receives_rel(txid, addr, vout.Index, vout.Value)
	}

	return nil
}

func decodeTx(tx *models.Tx, blockHash string) (err error) {
	txid := hex.EncodeToString(serial.ReverseHex(tx.Hash))
	// Add Transaction Node to neo4j graph
	neo4jdb.Add_transaction(txid, blockHash)

	if err = decodeInputs(tx, txid); err != nil {
		return err
	}
	if err = decodeOutputs(tx, txid); err != nil {
		return err
	}

	return
}

func decodeBlockHeader(bH *models.BlockHeader) error {
	// Add Block node to neo4j db
	neo4jdb.Add_block(bH)
	return nil
}
