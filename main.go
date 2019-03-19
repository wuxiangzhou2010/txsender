package main

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wuxiangzhou2010/txsender/config"
	"github.com/wuxiangzhou2010/txsender/sender"
)

var silent bool
var totalSent int32

var cfg *config.Config
var cons []*ethclient.Client

func main() {

	ctx := context.Background()

	sender.UpdateNonce(ctx, cons[0])

	senderTicker := time.NewTicker(1 * time.Second)
	defer senderTicker.Stop()

	// update total tx generated and sent
	printTicker := time.NewTicker(1 * time.Second)
	defer printTicker.Stop()
	want := cfg.Rate * cfg.Last

	var generator txGenerator
	generator = &signedTxGenerator{}
	generator.InitGenerator(cfg)
	txChannel := generator.GenerateTxs()

	for {
		select {
		case <-senderTicker.C:
			go sendTx(ctx, cons, txChannel, cfg.Rate)
			// fmt.Println("signed tx ", <-txChannel)

		case <-printTicker.C:
			sent := atomic.LoadInt32(&totalSent)
			log.Println("total tx sent ", sent)
			if want == sent {
				log.Println("all txs sent...")
				printTicker.Stop()
				return
			}

		}
	}
}

func sendTx(ctx context.Context, cons []*ethclient.Client, txsCh chan *types.Transaction, count int32) {
	for _, conn := range cons {
		go txWorker(ctx, conn, txsCh, count)
	}
}

func txWorker(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction, count int32) {
	var c int32
	for signedTx := range txsCh {
		err := conn.SendTransaction(ctx, signedTx)
		fmt.Println("sent a trancaction")
		if err != nil {
			fmt.Printf("error signedTx %+v\n", signedTx)
			log.Fatal("SendTransaction error ", err, signedTx)
		}
		atomic.AddInt32(&totalSent, 1)
		c++
		if c == count {
			return
		}
	}
}

func getConnections(rpcEndPoints []string) ([]*ethclient.Client, error) {
	var cons []*ethclient.Client
	for _, endPoint := range rpcEndPoints {
		conn, err := ethclient.Dial(endPoint)
		if err != nil {
			return nil, fmt.Errorf("can't establish connection to [%s], [error]: %v", endPoint, err)
		}
		cons = append(cons, conn)
	}
	return cons, nil
}

func init() {

	//init configuration
	cfg = config.GetConfig()
	silent = cfg.Silent

	//init connection
	var err error
	cons, err = getConnections(cfg.Endpoints)

	if err != nil {
		fmt.Println("getConnections failed, ", err)
		return
	}
}
