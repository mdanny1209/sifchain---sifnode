package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func NewMsgCreateDistribution(signer sdk.AccAddress, DistributionName string, DistributionType DistributionType, input []types.Input, output []types.Output) MsgCreateDistribution {

	return MsgCreateDistribution{
		Signer: signer.String(),
		Distribution: &Distribution{
			DistributionName: DistributionName,
			DistributionType: DistributionType,
		},
		Input:  input,
		Output: output,
	}
}

func (m MsgCreateDistribution) Route() string {
	return RouterKey
}

func (m MsgCreateDistribution) Type() string {
	return MsgTypeCreateDistribution
}

func (m MsgCreateDistribution) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	if m.Distribution.DistributionName == "" {
		return sdkerrors.Wrap(ErrInvalid, "Name cannot be empty")
	}

	err = types.ValidateInputsOutputs(m.Input, m.Output)
	if err != nil {
		return err
	}

	return nil
}

func NewMsgCreateUserClaim(signer sdk.AccAddress, claimType DistributionType) MsgCreateUserClaim {
	return MsgCreateUserClaim{
		Signer:        signer.String(),
		UserClaimType: claimType,
	}
}
func (m MsgCreateDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreateDistribution) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

func (m MsgCreateUserClaim) Route() string {
	return RouterKey
}

func (m MsgCreateUserClaim) Type() string {
	return MsgTypeCreateClaim
}

func (m MsgCreateUserClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	return nil
}

func (m MsgCreateUserClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreateUserClaim) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
