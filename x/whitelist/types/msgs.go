package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgUpdateWhitelist{}

func (m *MsgUpdateWhitelist) Route() string {
	return RouterKey
}

func (m *MsgUpdateWhitelist) Type() string {
	return "update"
}

func (m *MsgUpdateWhitelist) ValidateBasic() error {
	if m.Denom == "" {
		return errors.New("no denom specified")
	}
	coin := sdk.Coin{
		Denom:  m.Denom,
		Amount: sdk.OneInt(),
	}
	if !coin.IsValid() {
		return errors.New("Denom is not valid")
	}

	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid from address")
	}
	if m.Decimals < 0 {
		return errors.New("Decimals cannot be negative")
	}

	return nil
}

func (m *MsgUpdateWhitelist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgUpdateWhitelist) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}
