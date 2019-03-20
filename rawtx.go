package main

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/wuxiangzhou2010/txsender/sender"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/wuxiangzhou2010/txsender/config"
	"github.com/wuxiangzhou2010/txsender/recipient"
)

type rawTxGenerator struct {
	task
	rawTxCh chan *types.Transaction
}

func (rg *rawTxGenerator) initGenerator(cfg *config.Config) {
	rg.config = cfg
	accs := sender.GetSender()
	rg.acc = accs[0]
	rg.total = uint32(cfg.Last * cfg.Rate)

	rg.rawTxCh = make(chan *types.Transaction, cfg.RawTxBuffer)
}

func (rg *rawTxGenerator) InitGenerator(cfg *config.Config) {
	rg.initGenerator(cfg)

}

func (rg *rawTxGenerator) generateTxs() chan *types.Transaction {

	value := big.NewInt(1)    // in wei (100 wei)
	gasPrice := big.NewInt(1) // in wei (30 gwei)
	gasLimit := uint64(21000) // in units
	go func() {
		var temp uint32
		from := rg.acc.Account.Address
		round := rg.config.TxPerRecipient

		log.Info("[rawTx] start to generate raw transactions")
		defer func() {
			log.Println("[rawTx] stop generate raw transacitons, total generated", temp)
		}()
		defer close(rg.rawTxCh)
		for rg.total > temp {
			//get one recipient
			to := recipient.GetRecipient()
			for i := 0; i < round; i++ {
				tx := types.NewTransaction(rg.acc.Nonce, to, value, gasLimit, gasPrice, nil)

				atomic.AddUint64(&rg.acc.Nonce, 1)
				log.Debug("generated a transaction")
				rg.rawTxCh <- tx
			}
			temp += uint32(round)
			if !silent {
				fmt.Println("[generateRawTx] generate tx  from ", from.Hex(), "to ", to.Hex(), "amount", round)
			}
		}

	}()
	return rg.rawTxCh
}

func (rg *rawTxGenerator) GenerateTxs() chan *types.Transaction {
	return rg.generateTxs()
}
