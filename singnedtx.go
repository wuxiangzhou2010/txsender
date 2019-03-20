package main

import (
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/wuxiangzhou2010/txsender/config"
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
		var signedTotal int32

		log.Info("[signTx] start to sign transaciton")
		defer func() { log.Info("[signTx] stop sign transactions, total signed", signedTotal) }()
		defer close(sg.signedTxCh)

		for rawtx := range rawtxs {

			log.Debug("sign raw tx, ", rawtx.Nonce())
			signedTx, err := sg.acc.Ks.SignTx(sg.acc.Account, rawtx, nil)
			if err != nil {
				log.Println("[txSigner] SignTx error", err, sg.acc.Account.Address.Hex())
			}
			sg.signedTxCh <- signedTx

			signedTotal++

			if signedTotal%200 == 0 {
				log.Println("[txSigner] signedTotal ", signedTotal)
			}
		}

	}()
	return sg.signedTxCh
}

func (sg *signedTxGenerator) GenerateTxs() chan *types.Transaction {
	return sg.generateSignedTxs()
}
