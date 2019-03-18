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

	endpoints := cfg.Endpoints

	rpcEndPoint := endpoints[0]

	conn, err := ethclient.Dial(rpcEndPoint)
	ctx := context.Background()
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}
	sender.InitSender()
	sender.UpdateNonce(ctx, conn)
}
