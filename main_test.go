package main

import (
	"context"
	"log"
	"testing"

	"./config"
	"./sender"
	"github.com/comatrix/go-comatrix/ethclient"
)

func TestInit(t *testing.T) {
	var cfg *config.Config
	cfg = config.GetConfig()

	endpoints := cfg.Endpoints

	rpcEndPoint := endpoints[0]

	conn, err := ethclient.Dial(rpcEndPoint)
	ctx := context.Background()
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}
	senderOkCh := make(chan struct{})
	sender.InitSender(senderOkCh)
	sender.UpdateNonce(ctx, conn)
}
