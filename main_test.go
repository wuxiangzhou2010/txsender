package main

import (
	"context"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wuxiangzhou2010/txsender/config"
	"github.com/wuxiangzhou2010/txsender/sender"
)

func TestInit(t *testing.T) {
	var cfg *config.Config
	cfg = config.GetConfig()

	endpoint := cfg.Endpoint

	rpcEndPoint := endpoint

	conn, err := ethclient.Dial(rpcEndPoint)
	ctx := context.Background()
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}
	sender.InitSender()
	if err := sender.UpdateNonce(ctx, conn); err != nil {
		log.Fatal("updateNonce failed")
	}
}
