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

	var cfg *config.Config
	cfg = config.GetConfig("./config.json")

	txsPerRound := cfg.Rate
	chainAmount := cfg.ChainAmount
	silent := cfg.Silent
	endpoints := cfg.Endpoints
	workerSize := cfg.Worker
	rpcEndPoint := endpoints[0]

	conn, err := ethclient.Dial(rpcEndPoint)
	ctx := context.Background()
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	sender.InitSender(chainAmount)
	sender.UpdateNonce(ctx, conn)

	var total int

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	// update tx sent periodically
	ticker1 := time.NewTicker(10 * time.Second)
	defer ticker1.Stop()

	// channel to buffer txs
	ch := make(chan *types.Transaction, cfg.TxBuffer)

	go sendTx(ctx, conn, ch, workerSize)
	fmt.Println("Start to send transactions...")
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

			fmt.Println("total tx sent ", total)

		}
	}
}

func sendTx(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction, workerSize int) {
	for i := 0; i < workerSize; i++ {
		go txworker(ctx, conn, txsCh)
	}
}

func txworker(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction) {
	for signedTx := range txsCh {
		err := conn.SendTransaction(ctx, signedTx)
		if err != nil {
			log.Fatal("SendTransaction error ", err, signedTx)
		}
	}

}
