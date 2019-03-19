package main

import (
	"fmt"
	"github.com/wuxiangzhou2010/txsender/config"
	"log"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"
)

type signedTxGenerator struct {
	rawTxGenerator
	signedTxCh chan *types.Transaction
}

func (sg *signedTxGenerator) InitGenerator(cfg *config.Config) {
	sg.initGenerator(cfg)
	sg.signedTxCh = make(chan *types.Transaction, cfg.SignedTxBuffer)
}

func (sg *signedTxGenerator) generateSignedTxs() chan *types.Transaction {
	rawtxs := sg.generateTxs()

	go func() {
		for rawtx := range rawtxs {

			fmt.Println("sign raw tx, ", rawtx.Nonce())
			signedTx, err := sg.acc.Ks.SignTx(sg.acc.Account, rawtx, nil)
			if err != nil {
				log.Println("[txSigner] SignTx error", err, sg.acc.Account.Address.Hex())
			}
			sg.signedTxCh <- signedTx
			atomic.AddUint32(&sg.total, 1)
			if sg.total%20000 == 0 {
				log.Println("[txSigner] signedTotal ", sg.total)
			}

		}
	}()
	return sg.signedTxCh
}

func (sg *signedTxGenerator) GenerateTxs() chan *types.Transaction {
	return sg.generateSignedTxs()
}