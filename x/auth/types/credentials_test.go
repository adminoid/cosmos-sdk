package types_test

import (
	"testing"

	cryptotypes "github.com/adminoid/cosmos-sdk/crypto/types"
	sdk "github.com/adminoid/cosmos-sdk/types"
	authtypes "github.com/adminoid/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

func TestNewModuleCrendentials(t *testing.T) {
	expected := sdk.MustAccAddressFromBech32("cosmos1fpn0w0yf4x300llf5r66jnfhgj4ul6cfahrvqsskwkhsw6sv84wsmz359y")

	credential := authtypes.NewModuleCredential("group", [][]byte{{0x20}, {0x0}})
	require.NoError(t, sdk.VerifyAddressFormat(credential.Address().Bytes()))
	addr, err := sdk.AccAddressFromHexUnsafe(credential.Address().String())
	require.NoError(t, err)
	require.Equal(t, expected.String(), addr.String())

	require.True(t, credential.Equals(authtypes.NewModuleCredential("group", [][]byte{{0x20}, {0x0}})))
	require.False(t, credential.Equals(authtypes.NewModuleCredential("group", [][]byte{{0x20}, {0x1}})))
	require.False(t, credential.Equals(authtypes.NewModuleCredential("group", [][]byte{{0x20}})))
}

func TestNewBaseAccountWithPubKey(t *testing.T) {
	expected := sdk.MustAccAddressFromBech32("cosmos1fpn0w0yf4x300llf5r66jnfhgj4ul6cfahrvqsskwkhsw6sv84wsmz359y")

	credential := authtypes.NewModuleCredential("group", [][]byte{{0x20}, {0x0}})
	account, err := authtypes.NewBaseAccountWithPubKey(credential)
	require.NoError(t, err)
	require.Equal(t, expected, account.GetAddress())
	require.Equal(t, credential, account.GetPubKey())
}

func TestNewBaseAccountWithPubKey_WrongCredentials(t *testing.T) {
	_, err := authtypes.NewBaseAccountWithPubKey(cryptotypes.PubKey(nil))
	require.Error(t, err)
}
