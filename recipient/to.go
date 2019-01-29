package recipient

import (
	"math/rand"

	"github.com/comatrix/go-comatrix/common"
	"github.com/comatrix/go-comatrix/crypto"
)

const recipientAmount = 1000

var recipientAddressString = []string{
	"0x9246eebcc9e71e5f69ca48c9fd1f39a5fd9ad3e8",
	"0x5c2f960a954be76c71b890287463ec81be020e43",
	"0x80371043454fd85c609860a8545f9456e6caef9d",
	"0x000b45d515b6a0098787571eb407caf8ff7a670a",
	"0x592490348b165b85d878735ee66c8439084d267a",
}

var recipients []common.Address

func toAddress(strings []string) []common.Address {
	var result []common.Address

	for _, v := range strings {
		result = append(result, common.HexToAddress(v))
	}
	return result
}

func initTo() {
	// from fixed address
	//recipients = toAddress(recipientAddressString)

	//from generated address
	generateAddress()
}

// GetRecipient get one recipient
func GetRecipient() common.Address {
	if recipients == nil {
		initTo()
	}

	return recipients[rand.Intn(len(recipients))]

}

func generateAddress() {

	for i := 0; i < recipientAmount; i++ {
		// Create an account
		key, err := crypto.GenerateKey()
		if err != nil {
			panic("GenerateKey failed ")
		}

		// Get the address
		address := crypto.PubkeyToAddress(key.PublicKey)

		recipients = append(recipients, address)
	}
}
