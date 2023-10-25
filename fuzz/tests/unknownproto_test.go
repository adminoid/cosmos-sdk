//go:build gofuzz || go1.18

package tests

import (
	"testing"

	"github.com/adminoid/cosmos-sdk/codec/unknownproto"
	"github.com/adminoid/cosmos-sdk/testutil/testdata"
)

func FuzzUnknownProto(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		msg := new(testdata.TestVersion2)
		resolver := new(unknownproto.DefaultAnyResolver)
		unknownproto.RejectUnknownFields(b, msg, true, resolver)
		unknownproto.RejectUnknownFields(b, msg, false, resolver)
	})
}
