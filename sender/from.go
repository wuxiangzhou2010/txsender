package sender

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/comatrix/go-comatrix/accounts"
	"github.com/comatrix/go-comatrix/accounts/keystore"
	"github.com/comatrix/go-comatrix/ethclient"
)

// Acc is an accout that with keystore
type Acc struct {
	Ks      *keystore.KeyStore
	Account accounts.Account
	Nonce   uint64
}

var senderAccounts []*Acc

// GetSender get one account from account array
func GetSender() []*Acc {
	if senderAccounts == nil {
		panic("GetSender nil senderAccounts ")
	}

	//return senderAccounts[rand.Intn(len(senderAccounts))]
	return senderAccounts
}

func (aa *Acc) String() string {
	return fmt.Sprintf("Acc{a:%v, Nonce:%v}", aa.Account.Address.Hex(), aa.Nonce)
}

func getAccountFromPath(filePath []string) []*Acc {

	var result []*Acc
	for _, path := range filePath {

		// get keystore accounts
		ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
		jsonBytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		password := "123"
		account, err := ks.Import(jsonBytes, password, password)
		if err != nil {
			log.Fatal(err)
		}

		ks.Unlock(account, "123")

		log.Println("imported ", "0x"+path[48:])
		result = append(result, &Acc{Ks: ks, Account: account})

	}
	if err := os.RemoveAll("./tmp"); err != nil {
		log.Fatal(err)
	}
	return result

}

// UpdateNonce update the nonce of accounts
func UpdateNonce(ctx context.Context, conn *ethclient.Client) {
	if senderAccounts == nil {
		panic("nil sender")
	}
	for _, v := range senderAccounts {

		nonce, err := conn.NonceAt(ctx, v.Account.Address, nil)
		if err != nil {
			panic("err")
		}
		v.Nonce = nonce
		log.Println("Get Nonce ok", v.Account.Address.Hex(), "  nonce:", v.Nonce)
	}

}

// InitSender init the sender
func InitSender(senderOkCh chan struct{}) {

	path := getPath()
	keyPaths := readKeystore(path)

	senderAccounts = getAccountFromPath(keyPaths)
	defer close(senderOkCh)
}

func getPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	keystoreDir := filepath.Join(dir, "keystore")

	//fmt.Println("current dir ", dir, "keystore dir ", keystoreDir)
	return keystoreDir
}

//func getSenderChs(len int, bufferSize int) []chan *types.Transaction {
//	var senderChs []chan *types.Transaction
//	for i := 0; i < len; i++ {
//		senderChs = append(senderChs, make(chan *types.Transaction, bufferSize))
//	}
//	return senderChs
//}
//
//func getOnceCh(chs []chan *types.Transaction) chan *types.Transaction {
//	return chs[rand.Intn(len(chs))]
//}
