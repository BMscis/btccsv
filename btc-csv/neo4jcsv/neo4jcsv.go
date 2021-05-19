package neo4jcsv

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/serial"
	log "github.com/sirupsen/logrus"
)

type Graphwriter struct {
	blocks     *os.File
	txs        *os.File
	addrs      *os.File
	before     *os.File
	belongs_to *os.File
	receives   *os.File
	sends      *os.File

	wblocks *csv.Writer
	wtxs    *csv.Writer
	waddr   *csv.Writer
	wbef    *csv.Writer
	wbel    *csv.Writer
	wre     *csv.Writer
	wse     *csv.Writer

	//	utxo map[string][]string
}

func NewWriter() *Graphwriter {
	var err error
	w := new(Graphwriter)
	w.blocks, err = os.Create("blocks.csv")
	checkError("Error creating blocks.csv", err)

	w.txs, err = os.Create("transactions.csv")
	checkError("Error creating transactions.csv", err)

	w.addrs, err = os.Create("addresses.csv")
	checkError("Error creating addresses.csv", err)

	w.before, err = os.Create("before_rel.csv")
	checkError("Error creating before_rel.csv", err)

	w.belongs_to, err = os.Create("belongs_to_rel.csv")
	checkError("Error creating belongs_to_rel.csv", err)

	w.receives, err = os.Create("receives_rel.csv")
	checkError("Error creating receives_rel.csv", err)

	w.sends, err = os.Create("sends_rel.csv")
	checkError("Error creating sends_rel.csv", err)

	w.wblocks = csv.NewWriter(w.blocks)
	w.wtxs = csv.NewWriter(w.txs)
	w.waddr = csv.NewWriter(w.addrs)
	w.wbef = csv.NewWriter(w.before)
	w.wbel = csv.NewWriter(w.belongs_to)
	w.wre = csv.NewWriter(w.receives)
	w.wse = csv.NewWriter(w.sends)

	//	w.utxo = make(map[string][]string, 7E7)

	return w
}

func (w *Graphwriter) Close() {
	w.wblocks.Flush()
	w.wtxs.Flush()
	w.waddr.Flush()
	w.wbef.Flush()
	w.wbel.Flush()
	w.wre.Flush()
	w.wse.Flush()

	w.blocks.Close()
	w.txs.Close()
	w.addrs.Close()
	w.before.Close()
	w.belongs_to.Close()
	w.receives.Close()
	w.sends.Close()
}

func (w *Graphwriter) Add_block(blockHeader *models.BlockHeader) {
	hash := hex.EncodeToString(serial.ReverseHex(blockHeader.Hash))
	// Write Block: (hash:ID,height:int)
	err := w.wblocks.Write([]string{hash, fmt.Sprint(blockHeader.NHeight)})
	checkError("Error writing block to csv:", err)

	if blockHeader.NHeight != 0 {
		// Write IS_BEFORE: (:START_ID,:END_ID)
		err = w.wbef.Write([]string{hex.EncodeToString(serial.ReverseHex(blockHeader.HashPrev)), hash})
		checkError("Error writing IS_BEFORE to csv", err)
	}
}

// TODO: Check what to do here!!!
//func Add_address(address string) {
//	result, err := Session.Run("MERGE (a:Address {address:$address}) RETURN a.address",
//		map[string]interface{}{"address": address})
//	if err != nil {
//		fmt.Println("Error creating Address")
//	}
//
//	_, err = result.Consume()
//	if err != nil {
//		fmt.Println("Error consuming after Address")
//	}
//}

func (w *Graphwriter) Add_transaction(txid string, blockHash string) {
	// Write Transaction: (txid:ID)
	err := w.wtxs.Write([]string{txid})
	checkError("Error writing transaction to csv:", err)

	// Write BELONGS_TO: (:START_ID,:END_ID)
	err = w.wbel.Write([]string{txid, blockHash})
	checkError("Error writing BELONGS_TO to csv:", err)
}

func (w *Graphwriter) Add_sends_rel(txid string, addr string, value uint64) {
	//	a := w.utxo[hex.EncodeToString(serial.ReverseHex(txInput.Hash))+fmt.Sprint(txInput.Index)]
	//	delete(w.utxo, hex.EncodeToString(serial.ReverseHex(txInput.Hash))+fmt.Sprint(txInput.Index))

	// Write SENDS: (:START_ID,value:int,:END_ID)
	err := w.wse.Write([]string{addr, fmt.Sprint(value), txid})
	checkError("Error writing SENDS to csv:", err)
}

func (w *Graphwriter) Add_receives_rel(txid string, address string, output_nr uint32, value uint64) {
	//	w.utxo[txid+fmt.Sprint(output_nr)] = []string{address, fmt.Sprint(value)}

	// Write Address: (address:ID)
	err := w.waddr.Write([]string{address})
	checkError("Error writing address to csv:", err)

	// Write RECEIVES: (:START_ID,value:int,output_nr:int,:END_ID
	err = w.wre.Write([]string{txid, fmt.Sprint(value), fmt.Sprint(output_nr), address})
	checkError("Error writing RECEIVES to csv:", err)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatalln(message, err)
	}
}
