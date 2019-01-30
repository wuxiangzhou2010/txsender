package main

import (
	"./config"
	"./recipient"
	"./sender"

	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/comatrix/go-comatrix/core/types"
	"github.com/comatrix/go-comatrix/ethclient"
)

func main() {

	cfg := config.GetConfig()

	txsPerRound := cfg.Rate

	silent := cfg.Silent
	endpoints := cfg.Endpoints
	workerSize := cfg.Worker
	rpcEndPoint := endpoints[0]

	conn, err := getConnection(rpcEndPoint)
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}
	ctx := context.Background()

	senderOkCh := make(chan struct{})

	go sender.InitSender(senderOkCh)
	<-senderOkCh

	sender.UpdateNonce(ctx, conn)

	var total int

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	// update total tx sent
	ticker1 := time.NewTicker(10 * time.Second)
	defer ticker1.Stop()

	// channel to buffer txs
	ch := make(chan *types.Transaction, cfg.TxBuffer)

	go sendTx(ctx, conn, ch, workerSize)

	log.Println("Start to send transactions...")
	for {
		select {
		case <-ticker.C:

			value := big.NewInt(100)            // in wei (1 eth)
			gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
			gasLimit := uint64(21000)           // in units

			//get one account
			account := sender.GetSender()

			from := account.Account.Address
			to := recipient.GetRecipient()
			for i := 0; i < txsPerRound/10; i++ {

				tx := types.NewTransaction(account.Nonce, from, to, value, gasLimit, gasPrice, nil, 0)

				signedTx, err := account.Ks.SignTx(account.Account, tx, nil)
				if err != nil {
					fmt.Println("signtx error", err, account.Account.Address.Hex())
				}
				ch <- signedTx

				account.Nonce = account.Nonce + 1
			}
			total += txsPerRound / 10
			if !silent {
				fmt.Println(" generate tx  from ", from.Hex(), "to ", to.Hex(), "amount", txsPerRound/10)
			}
		case <-ticker1.C:

			log.Println("total tx sent ", total)

		}
	}
}

func sendTx(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction, workerSize int) {
	for i := 0; i < workerSize; i++ {
		go txWorker(ctx, conn, txsCh)
	}
}

func txWorker(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction) {
	for signedTx := range txsCh {
		err := conn.SendTransaction(ctx, signedTx)
		if err != nil {
			log.Fatal("SendTransaction error ", err, signedTx)
		}
	}

}

func getConnection(rpcEndPoint string) (*ethclient.Client, error) {
	return ethclient.Dial(rpcEndPoint)
}
