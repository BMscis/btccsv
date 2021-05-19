package grapho

import "github.com/chainswatch/bclib/models"

type Writer interface {
	Close()
	Add_block(blockHeader *models.BlockHeader)
	Add_transaction(txid string, blockHash string)
	Add_sends_rel(txid string, addr string, value uint64)
	Add_receives_rel(txid string, address string, output_nr uint32, value uint64)
}
