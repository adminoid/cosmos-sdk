package tx

import (
	"github.com/adminoid/cosmos-sdk/codec/types"
)

// TxExtensionOptionI defines the interface for tx extension options
type ExtensionOptionI interface{}

// unpackTxExtensionOptionsI unpacks Any's to TxExtensionOptionI's.
func unpackTxExtensionOptionsI(unpacker types.AnyUnpacker, anys []*types.Any) error {
	for _, any := range anys {
		var opt ExtensionOptionI
		err := unpacker.UnpackAny(any, &opt)
		if err != nil {
			return err
		}
	}

	return nil
}
