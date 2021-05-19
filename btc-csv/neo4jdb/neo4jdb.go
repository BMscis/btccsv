package neo4jdb

import (
	"encoding/hex"
	"fmt"

	"github.com/chainswatch/bclib/models"
	"github.com/chainswatch/bclib/serial"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var (
	Driver  neo4j.Driver
	Session neo4j.Session
	err     error
)

//func init() {
//	Driver, err = neo4j.NewDriver("bolt://localhost:7687",
//		neo4j.BasicAuth("neo4j", "blob", ""))
//
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	Session, err = Driver.Session(neo4j.AccessModeWrite)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//}

func Close() {
	Driver.Close()
	Session.Close()
}

func run_cypher(cypher string) {
	result, err := Session.Run(cypher,
		map[string]interface{}{})
	if err != nil {
		fmt.Println("Error creating constraints")
	}

	_, err = result.Consume()
	if err != nil {
		fmt.Println("Error consuming after constraints")
		fmt.Println(err.Error())
	}
}

func Add_constraints() {
	result, err := Session.Run("CALL db.constraints", map[string]interface{}{})
	if err != nil {
		fmt.Println("Error checking for constraints")
	}

	_, err = result.Consume()
	if err != nil {
		fmt.Println("Error consuming after checking for constraints")
	}

	fmt.Println(result.Next())

	if result.Next() == false {

		run_cypher("CREATE CONSTRAINT ON (b:Block) ASSERT b.hash IS UNIQUE")
		run_cypher("CREATE CONSTRAINT ON (t:Transaction) ASSERT t.txid IS UNIQUE")
		run_cypher("CREATE CONSTRAINT ON (a:Address) ASSERT a.address IS UNIQUE")
	} else {
		fmt.Println("Constraints already exist")
	}
}

func Add_block(blockHeader *models.BlockHeader) {
	if blockHeader.NHeight != 0 {
		result, err := Session.Run("MATCH (prev:Block {hash:$pHash})\n"+
			"MERGE (b:Block {hash:$hash, height:$height})\n"+
			"MERGE (prev)-[:IS_BEFORE]->(b) RETURN b.height",
			map[string]interface{}{"pHash": hex.EncodeToString(serial.ReverseHex(blockHeader.HashPrevBlock)),
				"hash":   hex.EncodeToString(serial.ReverseHex(blockHeader.HashBlock)),
				"height": blockHeader.NHeight})
		if err != nil {
			fmt.Println("Error creating Block and relation After")
		}

		_, err = result.Consume()
		if err != nil {
			fmt.Println("Error consuming after Block")
		}
	} else {
		result, err := Session.Run(
			"MERGE (b:Block {hash:$hash, height:$height})",
			map[string]interface{}{"hash": hex.EncodeToString(serial.ReverseHex(blockHeader.HashBlock)),
				"height": blockHeader.NHeight})
		if err != nil {
			fmt.Println("Error creating Block and relation After")
		}

		_, err = result.Consume()
		if err != nil {
			fmt.Println("Error consuming after Block")
		}
	}
}

func Add_address(address string) {
	result, err := Session.Run("MERGE (a:Address {address:$address}) RETURN a.address",
		map[string]interface{}{"address": address})
	if err != nil {
		fmt.Println("Error creating Address")
	}

	_, err = result.Consume()
	if err != nil {
		fmt.Println("Error consuming after Address")
	}
}

func Add_transaction(txid string, blockHash string) {
	result, err := Session.Run("MATCH (b:Block {hash:$hash})\n"+
		"MERGE (t:Transaction {txid:$txid})\n"+
		"MERGE (t)-[:BELONGS_TO]->(b) RETURN t.txid",
		map[string]interface{}{"txid": txid, "hash": blockHash})
	if err != nil {
		fmt.Println("Error creating Transaction")
	}

	_, err = result.Consume()
	if err != nil {
		fmt.Println("Error consuming after Transaction")
	}
}

func Add_sends_rel(txid string, txInput *models.TxInput) {
	result, err := Session.Run("MATCH (pTx:Transaction {txid:$pTxid})\n"+
		"MATCH (a:Address)<-[r:RECEIVES {output_nr:$nr}]-(pTx)\n"+
		"MATCH (t:Transaction {txid:$txid})\n"+
		"CREATE (a)-[:SENDS {value:r.value}]->(t) RETURN t.txid",
		map[string]interface{}{
			"txid":  txid,
			"pTxid": hex.EncodeToString(serial.ReverseHex(txInput.Hash)),
			"nr":    txInput.Index})
	if err != nil {
		fmt.Println("Error creating SENDS")
	}

	_, err = result.Consume()
	if err != nil {
		fmt.Println("Error consuming after SENDS")
	}
}

func Add_receives_rel(txid string, address string, output_nr uint32, value uint64) {
	result, err := Session.Run("MATCH (a:Address {address:$address})\n"+
		"MATCH (t:Transaction {txid:$txid})\n"+
		"CREATE (t)-[:RECEIVES {value:$value, output_nr:$output_nr}]->(a) RETURN t.txid",
		map[string]interface{}{"txid": txid, "address": address, "output_nr": output_nr, "value": value})
	if err != nil {
		fmt.Println("Error creating RECEIVES")
	}

	_, err = result.Consume()
	if err != nil {
		fmt.Println("Error consuming after RECEIVES")
	}
}
