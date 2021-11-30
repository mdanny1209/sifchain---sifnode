package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const (
	ModuleName = "margin"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route
	QuerierRoute = ModuleName

	// RouterKey is the msg router key
	RouterKey = ModuleName
)

func (mtp MTP) Validate() error {
	if mtp.Asset == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no asset specified")
	}
	if mtp.Address == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no address specified")
	}

	return nil
}
