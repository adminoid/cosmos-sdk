package testutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adminoid/cosmos-sdk/crypto/hd"
	"github.com/adminoid/cosmos-sdk/crypto/keyring"
	"github.com/adminoid/cosmos-sdk/types"
)

type TestAccount struct {
	Name    string
	Address types.AccAddress
}

func CreateKeyringAccounts(t *testing.T, kr keyring.Keyring, num int) []TestAccount {
	accounts := make([]TestAccount, num)
	for i := range accounts {
		record, _, err := kr.NewMnemonic(
			fmt.Sprintf("key-%d", i),
			keyring.English,
			types.FullFundraiserPath,
			keyring.DefaultBIP39Passphrase,
			hd.Secp256k1)
		assert.NoError(t, err)

		addr, err := record.GetAddress()
		assert.NoError(t, err)

		accounts[i] = TestAccount{Name: record.Name, Address: addr}
	}

	return accounts
}
