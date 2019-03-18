package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wuxiangzhou2010/txsender/config"
	"github.com/wuxiangzhou2010/txsender/recipient"
	"github.com/wuxiangzhou2010/txsender/sender"
)

var silent bool
var totalSent int32

var cfg *config.Config

func main() {

	cfg = config.GetConfig()

	silent = cfg.Silent
	endpoints := cfg.Endpoints

	conns := getConnections(endpoints)

	ctx := context.Background()

	senderOkCh := make(chan struct{})

	go sender.InitSender(senderOkCh)
	<-senderOkCh

	sender.UpdateNonce(ctx, conns[0])

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// update total tx generated and sent
	tickerPrint := time.NewTicker(1 * time.Second)
	defer tickerPrint.Stop()
	want := cfg.Rate * cfg.Last

	senderCh := make(chan *types.Transaction, cfg.SignedTxBuffer)
	go generateTx(senderCh, want)

	for {
		select {
		case <-ticker.C:
			go sendTx(ctx, conns, senderCh, cfg.Rate)

		case <-tickerPrint.C:
			sent := atomic.LoadInt32(&totalSent)
			log.Println("total tx sent ", sent)
			if want == sent {
				log.Println("all txs sent...")
				tickerPrint.Stop()
			}

		}
	}
}

func generateTx(senderCh chan *types.Transaction, total int32) {
	log.Println("[generateTx] Start to Generate", total, " transactions...")

	//get one sender
	accounts := sender.GetSender()
	accountsLen := len(accounts)
	if accountsLen == 0 {
		panic("[generateTx] wrong accounts")
	}

	totalPerAccount := total / int32(accountsLen)

	//var signedTotal uint32
	var wgAll sync.WaitGroup
	for _, account := range accounts {
		wgAll.Add(1)
		go func(acc *sender.Acc, w *sync.WaitGroup) {
			//rawTxCh := make(chan *types.Transaction, cfg.RawTxBuffer)
			go generateRawTx(senderCh, acc, totalPerAccount)

			//var wg sync.WaitGroup
			//wg.Add(10)
			//for i := 0; i < 10; i++ {
			//	go txSigner(rawTxCh, signedTxCh, acc, &signedTotal, &wg)
			//}
			//wg.Wait()
			w.Done()
		}(account, &wgAll)
	}
	wgAll.Wait()
	log.Println("[generateTx] all signed txs are generated")

}
func generateRawTx(rawTxCh chan *types.Transaction, account *sender.Acc, rawCount int32) {
	defer close(rawTxCh)

	value := big.NewInt(1)    // in wei (100 wei)
	gasPrice := big.NewInt(1) // in wei (30 gwei)
	gasLimit := uint64(21000) // in units
	from := account.Account.Address

	var totalRaw int
	round := cfg.TxPerRecipient
	for rawCount > int32(totalRaw) {
		//get one recipient
		to := recipient.GetRecipient()
		for i := 0; i < round; i++ {
			tx := types.NewTransaction(account.Nonce, to, value, gasLimit, gasPrice, nil)

			atomic.AddUint64(&account.Nonce, 1)
			rawTxCh <- tx
		}
		totalRaw += round
		if !silent {
			fmt.Println("[generateRawTx] generate tx  from ", from.Hex(), "to ", to.Hex(), "amount", 20)
		}

	}
	log.Println("[generateRawTx] all raw txs are generated, totalRaw", totalRaw)

}

//
//func txSigner(rawTxCh chan *types.Transaction, signedTxCh chan *types.Transaction, account *sender.Acc, signedTotal *uint32, wg *sync.WaitGroup) {
//	for tx := range rawTxCh {
//		signedTx, err := account.Ks.SignTx(account.Account, tx, nil)
//		if err != nil {
//			log.Println("[txSigner] SignTx error", err, account.Account.Address.Hex())
//		}
//		signedTxCh <- signedTx
//		atomic.AddUint32(signedTotal, 1)
//		if *signedTotal%20000 == 0 {
//			log.Println("[txSigner] signedTotal ", *signedTotal)
//		}
//	}
//	wg.Done()
//}

func sendTx(ctx context.Context, conns []*ethclient.Client, txsCh chan *types.Transaction, count int32) {
	for _, conn := range conns {
		go txWorker(ctx, conn, txsCh, count)
	}
}

func txWorker(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction, count int32) {
	var c int32
	for signedTx := range txsCh {
		err := conn.SendTransaction(ctx, signedTx)
		if err != nil {
			log.Fatal("SendTransaction error ", err, signedTx)
		}
		atomic.AddInt32(&totalSent, 1)
		c++
		if c == count {
			return
		}
	}
}

func getConnections(rpcEndPoints []string) []*ethclient.Client {
	var conns []*ethclient.Client
	for _, endPoint := range rpcEndPoints {
		conn, err := ethclient.Dial(endPoint)
		if err != nil {
			panic("")
		}
		conns = append(conns, conn)
	}
	return conns
}
