package recipient

import (
	"math/rand"

	"github.com/comatrix/go-comatrix/common"
	"github.com/comatrix/go-comatrix/crypto"
)

const recipientAmount = 1000

var recipients []common.Address

func toAddress(strings []string) []common.Address {
	var result []common.Address

	for _, v := range strings {
		result = append(result, common.HexToAddress(v))
	}
	return result
}

func initTo() {
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
