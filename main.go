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

	value := big.NewInt(100)            // in wei (1 eth)
	gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	gasLimit := uint64(21000)           // in units

	var total int
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	ticker1 := time.NewTicker(10 * time.Second)
	defer ticker1.Stop()

	fmt.Println("Start to send transactions...")
	for {
		select {
		case <-ticker.C:
			//get one account
			account := sender.GetSender()
			//fmt.Printf("select account %#v", account.Account.Address.Hex())
			from := account.Account.Address
			to := recipient.GetRecipient()

			for i := 0; i < *txsPerRound/10; i++ {

				tx := types.NewTransaction(account.Nonce, from, to, value, gasLimit, gasPrice, nil, 0)

				signedTx, err := account.Ks.SignTx(account.Account, tx, nil)
				if err != nil {
					fmt.Println("signtx error", err)
				}

				err = conn.SendTransaction(context.Background(), signedTx)
				if err != nil {
					fmt.Println("   from ", from.Hex(), "to ", to.Hex())
					glog.Fatal("SendTransaction error ", err)

				}
				if !*silent {
					fmt.Printf("tx sent: %s %v\n", signedTx.Hash().Hex(), signedTx.Nonce())

				}
				account.Nonce = account.Nonce + 1
			}
			total += *txsPerRound / 10
			if !*silent {
				fmt.Print("total tx sent ", total)
				fmt.Println("   from ", from.Hex(), "to ", to.Hex())
			}
		case <-ticker1.C:

			fmt.Println("total tx sent ", total)

		}
	}
}
