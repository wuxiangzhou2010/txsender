package sender

import (
	"testing"
)

func TestReadKeystore(t *testing.T) {
	path := "../keystore"
	readKeystore(path)
}
