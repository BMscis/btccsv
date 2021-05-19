package grapho

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/BMscis/btccsv/protos"
	"github.com/chainswatch/bclib/btc"
	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/serial"
	"github.com/golang/protobuf/proto"
	logs "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

func decodeInputs(tx *models.Tx, txid string, writer Writer, utxo *leveldb.DB) error {

	txOutput := new(protos.TxOutput)

	for _, vin := range tx.Vin {
		// TODO: What to do with blockreward inputs?
		if vin.Index == 0xFFFFFFFF { // block reward

		} else {
			// TODO: return an error
			data, err := utxo.Get([]byte(hex.EncodeToString(serial.ReverseHex(vin.Hash))+fmt.Sprint(vin.Index)), nil)
			if err != nil {
				logs.Warnf("Utxo with txid: " + hex.EncodeToString(serial.ReverseHex(vin.Hash)) + " output_nr: " + fmt.Sprint(vin.Index))
				logs.Warnln("Txid of current transaction: " + txid)
				log.Fatal("Utxo not found: ", err)
			}
			if err = proto.Unmarshal(data, txOutput); err != nil {
				log.Fatal("Unmarshalling error: ", err)
			}
			writer.Add_sends_rel(txid, txOutput.Addr, txOutput.Value)
			// Delete entry from utxo
			utxo.Delete([]byte(hex.EncodeToString(serial.ReverseHex(vin.Hash))+fmt.Sprint(vin.Index)), nil)
		}
	}

	return nil
}

func decodeOutputs(tx *models.Tx, txid string, writer Writer, utxo *leveldb.DB) error {
	var err error
	var addr string
	txOutput := new(protos.TxOutput)

	// Go through all Outputs
	// - add address to neo4j db if it doesn't already exist
	// - add relationship RECEIVES with BTC value and output_nr
	for _, vout := range tx.Vout {
		if addr, err = btc.DecodeAddr(vout.AddrType, vout.Addr); err != nil {
			log.Fatal(fmt.Sprintf("DecodeOutputs: %s: %d %x", err.Error(), serial.ReverseHex(tx.Hash)))
		}
		txOutput.Addr = addr
		txOutput.Value = vout.Value

		data, err := proto.Marshal(txOutput)
		if err != nil {
			log.Fatal("marshalling error: ", err)
		}
		utxo.Put([]byte(txid+fmt.Sprint(vout.Index)), data, nil)

		writer.Add_receives_rel(txid, addr, vout.Index, vout.Value)
	}

	return nil
}

func decodeTx(tx *models.Tx, blockHash string, writer Writer, utxo *leveldb.DB) (err error) {
	txid := hex.EncodeToString(serial.ReverseHex(tx.Hash))
	// Add Transaction Node to neo4j graph
	writer.Add_transaction(txid, blockHash)

	if err = decodeInputs(tx, txid, writer, utxo); err != nil {
		return err
	}
	if err = decodeOutputs(tx, txid, writer, utxo); err != nil {
		return err
	}

	return
}

func decodeBlockHeader(bH *models.BlockHeader, writer Writer) error {
	// Add Block node to neo4j db
	writer.Add_block(bH)
	return nil
}
