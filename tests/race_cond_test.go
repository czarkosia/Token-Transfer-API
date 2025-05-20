package tests

import (
	"context"
	"log"
	"sync"
	"testing"
	"token-transfer-api/db"
	"token-transfer-api/graph"
	"token-transfer-api/models"

	"github.com/stretchr/testify/assert"
)

func TestRaceCondition(t *testing.T) {
	db.DBconnect()
	resolver := &graph.Resolver{}

	from_address := "0x000from_address"
	to_address := [3]string{"0x000to_address_1", "0x000to_address_2", "0x000to_address_3"}

	// assert.Error(t, db.DB.Where("address = ?", from_address).First(&models.Wallet{}).Error)
	from_wallet := models.Wallet{Address: from_address, Balance: 10}
	assert.NoError(t, db.DB.Create(&from_wallet).Error)
	var to_wallet [3]models.Wallet
	for i := range 3 {
		// assert.Error(t, db.DB.Where("address = ?", to_address[i]).First(&models.Wallet{}).Error)
		to_wallet[i] = models.Wallet{Address: to_address[i], Balance: 10}
		assert.NoError(t, db.DB.Create(&to_wallet[i]).Error)
	}

	var wg *sync.WaitGroup = new(sync.WaitGroup)
	wg.Add(3)

	go func() {
		defer wg.Done()
		_, err := resolver.Mutation().Transfer(context.TODO(), to_address[0], from_address, 1)
		if err == nil {
			log.Printf("Transfer +1 accepted")
		} else {
			log.Printf("Transfer +1 rejected")
		}
	}()

	go func() {
		defer wg.Done()
		_, err := resolver.Mutation().Transfer(context.TODO(), from_address, to_address[1], 4)
		if err == nil {
			log.Printf("Transfer -4 accepted")
		} else {
			log.Printf("Transfer -4 rejected")
		}
	}()

	go func() {
		defer wg.Done()
		_, err := resolver.Mutation().Transfer(context.TODO(), from_address, to_address[2], 7)
		if err == nil {
			log.Printf("Transfer -7 accepted")
		} else {
			log.Printf("Transfer -7 rejected")
		}
	}()

	wg.Wait()

	var updated_from models.Wallet
	assert.NoError(t, db.DB.First(&updated_from, "address = ?", from_address).Error)
	assert.True(t, updated_from.Balance >= int32(0))

	db.DB.Where("address = ?", from_address).Delete(&models.Wallet{})
	// assert.Error(t, db.DB.Where("address = ?", from_address).First(&models.Wallet{}).Error)
	for i := range 3 {
		db.DB.Where("address = ?", to_address[i]).Delete(&models.Wallet{})
		// assert.Error(t, db.DB.Where("address = ?", to_address[i]).First(&models.Wallet{}).Error)
	}
}
