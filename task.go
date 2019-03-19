package main

import (
	"github.com/wuxiangzhou2010/txsender/config"
	"github.com/wuxiangzhou2010/txsender/sender"
)

type task struct {
	config *config.Config
	acc    *sender.Acc
	total  uint32
}
