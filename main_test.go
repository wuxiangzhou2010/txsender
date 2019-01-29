package main

import (
	"context"
	"log"
	"testing"

	"./config"
	"./sender"
	"github.com/comatrix/go-comatrix/ethclient"
)

func TestMain(t *testing.T) {
	var cfg *config.Config
	cfg = config.GetConfig("./config.json")

	// txsPerRound := cfg.Rate
	chainAmount := cfg.ChainAmount
	// silent := cfg.Silent
	endpoints := cfg.Endpoints
	// workerSize := cfg.Worker
	rpcEndPoint := endpoints[0]

	conn, err := ethclient.Dial(rpcEndPoint)
	ctx := context.Background()
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	sender.InitSender(chainAmount)
	sender.UpdateNonce(ctx, conn)
}
