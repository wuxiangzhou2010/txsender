package sender

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
var chainAmount int

// GetSender get one account from account array
func GetSender() *Acc {
	if senderAccounts == nil {
		panic("GetSender nil senderAccounts ")
	}

	return senderAccounts[rand.Intn(len(senderAccounts))]
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

		fmt.Println("imported ", "0x"+path[48:])
		result = append(result, &Acc{Ks: ks, Account: account})

	}

	return result

}

// UpdateNonce update the nonce of accounts
func UpdateNonce(ctx context.Context, conn *ethclient.Client) {
	if senderAccounts == nil {
		panic("nil sender")
	}
	for _, v := range senderAccounts {
		//get Nonce
		nonce, err := conn.NonceAt(ctx, v.Account.Address, nil)
		if err != nil {
			panic("err")
		}
		v.Nonce = nonce
		log.Println("Get Nonce ok", v.Account.Address.Hex(), "  ", v.Nonce)
	}

}

// InitSender init the sender
func InitSender(amount int) {
	chainAmount = amount

	path := getPath()
	keypaths := readKeystore(path)

	senderAccounts = getAccountFromPath(keypaths)
}

func getPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	keystoreDir := filepath.Join(dir, "keystore")

	fmt.Println("current dir ", dir, "keystore dir ", keystoreDir)
	return keystoreDir
}
