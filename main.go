//package main
//
//import (
//"context"
//"fmt"
//"log"
//
//"github.com/ethereum/go-ethereum/common"
//"github.com/ethereum/go-ethereum/ethclient"
//)
//
//func main() {
//	conn, err := ethclient.Dial("http://3.0.218.180:8545")
//	if err != nil {
//		log.Fatal("Whoops something went wrong!", err)
//	}
//
//	ctx := context.Background()
//	tx, pending, _ := conn.TransactionByHash(ctx, common.HexToHash("0x30999361906753dbf60f39b32d3c8fadeb07d2c0f1188a32ba1849daac0385a8"))
//	if !pending {
//		fmt.Println(tx)
//	}
//
//	account := common.HexToAddress("0xf1854f20482b211b8d62747e5fd62144efba8def")
//	balance, err := conn.BalanceAt(context.Background(), account, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(balance) // 25893180161173005034
//
//}

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/comatrix/go-comatrix/accounts/keystore"
	"github.com/comatrix/go-comatrix/common"
	"github.com/comatrix/go-comatrix/core/types"
	"github.com/comatrix/go-comatrix/ethclient"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"time"
)

const testAddress = "0x5c2f960a954be76c71b890287463ec81be020e43"

//const contractAddress = "0xd50e3913b461d01baceb687b4330fe5cb56dc6b5"

func getAddres() common.Address {
	addresses := []common.Address{
		common.HexToAddress("0x9246eebcc9e71e5f69ca48c9fd1f39a5fd9ad3e8"),
		common.HexToAddress("0x5c2f960a954be76c71b890287463ec81be020e43"),
		common.HexToAddress("0x80371043454fd85c609860a8545f9456e6caef9d"),
		common.HexToAddress("0x000b45d515b6a0098787571eb407caf8ff7a670a"),
		common.HexToAddress("0x592490348b165b85d878735ee66c8439084d267a"),
	}

	return addresses[rand.Intn(len(addresses))]

}

func main() {

	txsPerRound := flag.Int("rate", 10, "txs per round")
	silent := flag.Bool("silent", true, "keep silent")

	var ipaddress_string string
	flag.StringVar(&ipaddress_string, "ip", "http://3.0.218.180:8546", "rpc endpoint")

	flag.Parse()
	fmt.Println("flags: rate ", *txsPerRound, "silent ", *silent, "ip ", ipaddress_string)
	conn, err := ethclient.Dial(ipaddress_string)
	//conn, err := ethclient.Dial("http://172.16.3.191:8545")
	if err != nil {
		log.Fatal("Whoops something went wrong!", err)
	}

	ctx := context.Background()

	//get transaction
	//tx, pending, _ := conn.TransactionByHash(ctx, common.HexToHash("0x30999361906753dbf60f39b32d3c8fadeb07d2c0f1188a32ba1849daac0385a8"))
	//if !pending {
	// fmt.Println(tx)
	//}

	//get nonce

	nonce, err := conn.NonceAt(ctx, common.HexToAddress(testAddress), nil)
	if err != nil {
		panic("err")

	}
	//fmt.Println("test nonce  ", testAddress, " is ", nonce)

	//get block
	//block, err := conn.BlockByNumber(ctx, big.NewInt(0), big.NewInt(1))
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("test block of subchainid 0, number 122 is ", block)

	//count, err := conn.TransactionCount(context.Background(), big.NewInt(0), block.Hash())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("test transaction count ", count) // 144

	//getBalance
	//balance, err := conn.BalanceAt(context.Background(), common.HexToAddress(testAddress), nil)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println("getBalance ", balance) // 25893180161173005034

	//get code
	//bytecode, err := conn.CodeAt(context.Background(), common.HexToAddress(contractAddress), nil) // nil is latest block
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//isContract := len(bytecode) > 0
	//
	//fmt.Printf("is contract: %v\n", isContract) // is contract: true

	//header, err := conn.HeaderByNumber(context.Background(), big.NewInt(0), nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("test HeaderByNumber ", header) // 5671744

	// get keystore accounts
	file := "./keystore/UTC--2018-11-05T07-13-33.829662100Z--5c2f960a954be76c71b890287463ec81be020e43"
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	password := "123"
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3

	//if err := os.Remove(file); err != nil {
	//	log.Fatal(err)
	//}

	value := big.NewInt(100)            // in wei (1 eth)
	gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	gasLimit := uint64(21000)           // in units
	fromAddress := common.HexToAddress("0x5C2f960a954bE76C71b890287463Ec81BE020e43")
	//toAddress := common.HexToAddress("0x592490348b165b85d878735ee66c8439084d267a")
	//NewTransaction(nonce uint64, from common.Address, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, chainID uint64)

	ks.Unlock(account, "123")

	var total int
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < *txsPerRound/10; i++ {

				tx := types.NewTransaction(nonce, fromAddress, getAddres(), value, gasLimit, gasPrice, nil, 0)

				signedTx, err := ks.SignTx(account, tx, nil)
				if err != nil {
					fmt.Println("signtx error", err)
				}

				err = conn.SendTransaction(context.Background(), signedTx)
				if err != nil {
					log.Fatal(err)
				}
				if !*silent {
					fmt.Printf("tx sent: %s %v\n", signedTx.Hash().Hex(), signedTx.Nonce()) // tx sent: 0x77006fcb3938f648e2cc65bafd27dec30b9bfbe9df41f78498b9c8b7322a249e
				}
				nonce = nonce + 1
			}
			total += *txsPerRound / 10
			fmt.Println("total tx sent ", total)
		}
	}
}
