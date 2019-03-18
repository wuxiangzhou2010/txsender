package sender

import (
	"testing"
)

func TestInitSender(t *testing.T) {
	senderOkCh := make(chan struct{})
	InitSender(senderOkCh)
}
