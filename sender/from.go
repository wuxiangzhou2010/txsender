package sender

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/comatrix/go-comatrix/accounts"
	"github.com/comatrix/go-comatrix/accounts/keystore"
	"github.com/comatrix/go-comatrix/ethclient"
	"io/ioutil"
	"log"
)

var keypaths = []string{
	"./keystore/UTC--2018-11-05T07-13-33.829662100Z--5c2f960a954be76c71b890287463ec81be020e43",
	"./keystore/UTC--2018-11-05T07-14-20.837583500Z--9246eebcc9e71e5f69ca48c9fd1f39a5fd9ad3e8",
	"./keystore/UTC--2019-01-11T09-44-43.418478500Z--80371043454fd85c609860a8545f9456e6caef9d",
	"./keystore/UTC--2019-01-11T09-45-22.450285600Z--000b45d515b6a0098787571eb407caf8ff7a670a",
}

type Acc struct {
	Ks      *keystore.KeyStore
	Account accounts.Account
	Nonce   uint64
}

var senderAccounts []*Acc
var chainAmount int

func GetSender() *Acc {
	if senderAccounts == nil {
		panic("GetSender nil senderAccounts ")
	}
	//if chainAmount == 2 {
	//	chainAmount = 1
	//}
	return senderAccounts[rand.Intn(chainAmount)]
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

		//fmt.Println("path ", filePath)
		fmt.Println("imported ", "0x"+path[48:])
		result = append(result, &Acc{Ks: ks, Account: account})

	}
	//fmt.Printf("\n\n%+v\n", result)
	return result

}

func UpdateNonce(conn *ethclient.Client, ctx context.Context) {
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
	}
	//fmt.Printf("accounts with Nonce %+v", senderAccounts)
	fmt.Println("Get Nonce ok")
}

func InitSender(amount int) {
	chainAmount = amount
	senderAccounts = getAccountFromPath(keypaths)
}
