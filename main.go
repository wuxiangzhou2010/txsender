package main

import (
	"./recipient"
	"./sender"

	"context"
	"flag"
	"fmt"
	"math/big"
	"time"

	"github.com/comatrix/go-comatrix/core/types"
	"github.com/comatrix/go-comatrix/ethclient"
	"github.com/golang/glog"
)

func main() {

	// parse flag
	txsPerRound := flag.Int("rate", 300, "txs per second")
	chainAmount := flag.Int("amount", 1, "chain amount")
	silent := flag.Bool("silent", true, "keep silent")

	var rpcEndPoint string
	flag.StringVar(&rpcEndPoint, "ip", "http://3.0.218.180:8546", "rpc endpoint")
	// flag.StringVar(&rpcEndPoint, "ip", "http://13.228.196.190:8546", "rpc endpoint")

	flag.Parse()
	fmt.Println("flags: rate ", *txsPerRound, "silent ", *silent, "rpc endpoint ", rpcEndPoint)

	conn, err := ethclient.Dial(rpcEndPoint)
	ctx := context.Background()
	if err != nil {
		glog.Fatal("Whoops something went wrong!", err)
	}

	sender.InitSender(*chainAmount)
	sender.UpdateNonce(ctx, conn)

	var total int

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	// update tx sent periodically
	ticker1 := time.NewTicker(10 * time.Second)
	defer ticker1.Stop()

	// channel to buffer txs
	ch := make(chan *types.Transaction, 500)

	go sendTx(ctx, conn, ch)
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

			for i := 0; i < *txsPerRound/10; i++ {

				tx := types.NewTransaction(account.Nonce, from, to, value, gasLimit, gasPrice, nil, 0)

				signedTx, err := account.Ks.SignTx(account.Account, tx, nil)
				if err != nil {
					fmt.Println("signtx error", err)
				}
				ch <- signedTx

				account.Nonce = account.Nonce + 1
			}
			total += *txsPerRound / 10
			if !*silent {
				fmt.Println(" generate tx  from ", from.Hex(), "to ", to.Hex(), "amount", *txsPerRound/10)
			}
		case <-ticker1.C:

			fmt.Println("total tx sent ", total)

		}
	}
}

func sendTx(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction) {
	for i := 0; i < 4; i++ {
		go txworker(ctx, conn, txsCh)
	}
}

func txworker(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction) {
	for signedTx := range txsCh {
		err := conn.SendTransaction(ctx, signedTx)
		if err != nil {
			glog.Fatal("SendTransaction error ", err)
		}
	}

}
