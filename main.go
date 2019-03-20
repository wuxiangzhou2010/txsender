package main

import (
	"context"
	"fmt"
	"os"

	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/wuxiangzhou2010/txsender/config"
	"github.com/wuxiangzhou2010/txsender/sender"
)

var silent bool
var totalSent int32

var cfg *config.Config
var con *ethclient.Client

func main() {

	ctx := context.Background()

	if err := sender.UpdateNonce(ctx, con); err != nil {
		log.Error("Update nonce err")
		return
	}

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

			go sendTx(ctx, con, txChannel, cfg.Rate)

		case <-printTicker.C:
			sent := atomic.LoadInt32(&totalSent)
			log.Println("[sendTx] total tx sent ", sent)
			if want == sent {
				log.Println("[sendTx] all txs sent...")
				printTicker.Stop()
				return
			}

		}
	}
}

func sendTx(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction, count int32) {

	go txWorker(ctx, conn, txsCh, count)

}

func txWorker(ctx context.Context, conn *ethclient.Client, txsCh chan *types.Transaction, count int32) {
	var c int32
	for signedTx := range txsCh {
		err := conn.SendTransaction(ctx, signedTx)
		log.Debug("sent a transaction")
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

func getConnections(rpcEndPoint string) (*ethclient.Client, error) {

	conn, err := ethclient.Dial(rpcEndPoint)
	if err != nil {
		return nil, fmt.Errorf("can't establish connection to [%s], [error]: %v", rpcEndPoint, err)
	}

	return conn, nil
}

func init() {

	//init configuration
	cfg = config.GetConfig()
	silent = cfg.Silent

	//init connection
	var err error
	con, err = getConnections(cfg.Endpoint)

	if err != nil {
		fmt.Println("getConnections failed, ", err)
		return
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

}
