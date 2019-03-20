package main

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/wuxiangzhou2010/txsender/config"
)

type txGenerator interface {
	InitGenerator(*config.Config)
	GenerateTxs() chan *types.Transaction
}
