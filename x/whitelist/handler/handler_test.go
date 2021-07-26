package handler_test

import (
	whitelist "github.com/Sifchain/sifnode/x/whitelist/handler"
	"github.com/Sifchain/sifnode/x/whitelist/test"
	"github.com/Sifchain/sifnode/x/whitelist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHandler(t *testing.T) {
	app, ctx, admin := test.CreateTestApp(false)
	handler := whitelist.NewHandler(app.WhitelistKeeper)
	tests := []struct {
		name           string
		msg            types.MsgUpdateWhitelist
		errorAssertion assert.ErrorAssertionFunc
		valueAssertion assert.ValueAssertionFunc
	}{
		{
			name: "Valid Test",
			msg: types.MsgUpdateWhitelist{
				From:     admin,
				Denom:    "TestDenom",
				Decimals: 18,
			},
			errorAssertion: assert.NoError,
			valueAssertion: assert.NotNil,
		},
		{
			name: "Non Admin Account",
			msg: types.MsgUpdateWhitelist{
				From:     sdk.AccAddress("addr2_______________").String(),
				Denom:    "TestDenom",
				Decimals: 18,
			},
			errorAssertion: assert.Error,
			valueAssertion: assert.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := handler(ctx, &tt.msg)
			tt.errorAssertion(t, err)
			tt.valueAssertion(t, res)
		})
	}
}
