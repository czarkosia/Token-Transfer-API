package tests

import (
	"context"
	"testing"
	"token-transfer-api/db"
	"token-transfer-api/graph"
	"token-transfer-api/models"

	"github.com/stretchr/testify/assert"
)

func TestSuccessfulTransfer(t *testing.T) {
	db.DBconnect()
	from_address := "0xtransfertestfrom123"
	to_address := "0xtransfertestto321"

	assert.Error(t, db.DB.Where("address = ?", from_address).First(&models.Wallet{}).Error)
	assert.Error(t, db.DB.Where("address = ?", to_address).First(&models.Wallet{}).Error)

	from_wallet := models.Wallet{Address: from_address, Balance: 100000}
	to_wallet := models.Wallet{Address: to_address, Balance: 500}
	assert.NoError(t, db.DB.Create(&from_wallet).Error)
	assert.NoError(t, db.DB.Create(&to_wallet).Error)

	assert.NoError(t, db.DB.Where("address = ?", from_address).First(&models.Wallet{}).Error)
	assert.NoError(t, db.DB.Where("address = ?", to_address).First(&models.Wallet{}).Error)

	resolver := &graph.Resolver{}
	_, err := resolver.Mutation().Transfer(context.TODO(), from_address, to_address, 1000)
	assert.NoError(t, err)

	var updated_from, updated_to models.Wallet
	db.DB.First(&updated_from, "address = ?", from_address)
	db.DB.First(&updated_to, "address = ?", to_address)

	assert.Equal(t, updated_from.Balance, int32(99000))
	assert.Equal(t, updated_to.Balance, int32(1500))

	db.DB.Where("address = ?", from_address).Delete(&models.Wallet{})
	db.DB.Where("address = ?", to_address).Delete(&models.Wallet{})

	assert.Error(t, db.DB.Where("address = ?", from_address).First(&models.Wallet{}).Error)
	assert.Error(t, db.DB.Where("address = ?", to_address).First(&models.Wallet{}).Error)
}

func TestFailedTransfer(t *testing.T) {
	db.DBconnect()
	from_address := "0xtransfertestfrom123"
	to_address := "0xtransfertestto321"
	resolver := &graph.Resolver{}
	var err error

	assert.Error(t, db.DB.Where("address = ?", from_address).First(&models.Wallet{}).Error)
	assert.Error(t, db.DB.Where("address = ?", to_address).First(&models.Wallet{}).Error)

	_, err = resolver.Mutation().Transfer(context.TODO(), from_address, to_address, 7)
	assert.Error(t, err)

	from_wallet := models.Wallet{Address: from_address, Balance: 10}
	to_wallet := models.Wallet{Address: to_address, Balance: 500}
	assert.NoError(t, db.DB.Create(&from_wallet).Error)
	assert.NoError(t, db.DB.Create(&to_wallet).Error)

	assert.NoError(t, db.DB.Where("address = ?", from_address).First(&models.Wallet{}).Error)
	assert.NoError(t, db.DB.Where("address = ?", to_address).First(&models.Wallet{}).Error)

	_, err = resolver.Mutation().Transfer(context.TODO(), from_address, to_address, 100)
	assert.Error(t, err)
	_, err = resolver.Mutation().Transfer(context.TODO(), from_address, to_address, -100)
	assert.Error(t, err)

	var updated_from, updated_to models.Wallet
	db.DB.First(&updated_from, "address = ?", from_address)
	db.DB.First(&updated_to, "address = ?", to_address)

	assert.Equal(t, updated_from.Balance, int32(10))
	assert.Equal(t, updated_to.Balance, int32(500))

	db.DB.Where("address = ?", from_address).Delete(&models.Wallet{})
	db.DB.Where("address = ?", to_address).Delete(&models.Wallet{})

	assert.Error(t, db.DB.Where("address = ?", from_address).First(&models.Wallet{}).Error)
	assert.Error(t, db.DB.Where("address = ?", to_address).First(&models.Wallet{}).Error)
}
