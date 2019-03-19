package recipient

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const recipientAmount = 1000

var recipients []common.Address

func initTo() {
	generateRecipientAddress()
}

// GetRecipient get one recipient
func GetRecipient() common.Address {
	if recipients == nil {
		initTo()
	}
	return recipients[rand.Intn(len(recipients))]

}

func generateRecipientAddress() {

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
