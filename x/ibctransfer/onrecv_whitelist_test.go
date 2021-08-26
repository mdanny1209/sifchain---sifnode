package ibctransfer_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	whitelistmocks "github.com/Sifchain/sifnode/x/tokenregistry/types/mock"
)

func TestIsRecvPacketAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "foochainid"}, false, nil)

	returningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}

	nonReturningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}

	// If token is returning it will have prefix that matches packet source port and channel
	// FROM COSMOS: This is the prefix that would have been prefixed to the denomination
	// on sender chain IF and only if the token originally came from the
	// receiving chain.
	// NOTE: We use SourcePort and SourceChannel here, because the counterparty
	// chain would have prefixed with DestPort and DestChannel when originally
	// receiving this coin as seen in the "sender chain is the source" condition.
	returningDenom := transfertypes.FungibleTokenPacketData{
		// If atom has a prefix when coming in,
		// it has an extra hop in the path received in ibc transfer OnRecvPacket().
		Denom: "transfer/channel-0/rowan",
	}

	whitelistedDenom := transfertypes.FungibleTokenPacketData{
		// When sender chain is the source,
		// it simply sends the base denom without path prefix
		Denom: "atom",
	}
	// If token is returning it will have prefix that matches packet source port and channel
	// FROM COSMOS: This is the prefix that would have been prefixed to the denomination
	// on sender chain IF and only if the token originally came from the
	// receiving chain.
	disallowedDenom := transfertypes.FungibleTokenPacketData{
		// If atom has a prefix when coming in,
		// it has an extra hop in the path received in ibc transfer OnRecvPacket().
		Denom: "transfer/channel-66/atom",
	}

	wl := whitelistmocks.NewMockKeeper(ctrl)

	// Case: Returning: FALSE, Whitelisted: TRUE, Permissions: TRUE
	// Expected Result: TRUE
	wl.EXPECT().
		CheckDenomPermissions(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3", []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT}).
		Return(true)
	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3").
		Return(true)
	got := ibctransfer.IsRecvPacketAllowed(ctx, wl, nonReturningTransferPacket, whitelistedDenom)
	require.Equal(t, got, true)
	// Case: Returning: FALSE, Whitelisted: FALSE, Permissions: TRUE
	// Expected Result: FALSE
	wl.EXPECT().
		CheckDenomPermissions(ctx,
			"ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2", []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT}).
		Return(false)
	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2").
		Return(false)
	got = ibctransfer.IsRecvPacketAllowed(ctx, wl, nonReturningTransferPacket, disallowedDenom)
	require.Equal(t, got, false)
	// Case: Returning: TRUE, Whitelisted: FALSE, Permissions: TRUE
	// Expected Result: TRUE
	wl.EXPECT().
		CheckDenomPermissions(ctx,
			"rowan", []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT}).
		Return(true)
	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"rowan").
		Return(true)
	got = ibctransfer.IsRecvPacketAllowed(ctx, wl, returningTransferPacket, returningDenom)
	require.Equal(t, got, true)
	// Case: Returning: TRUE, Whitelisted: True, Permissions: TRUE
	// Expected Result: TRUE
	wl.EXPECT().
		CheckDenomPermissions(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3", []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT}).
		Return(true)
	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3").
		Return(true)
	got = ibctransfer.IsRecvPacketAllowed(ctx, wl, returningTransferPacket, whitelistedDenom)
	require.Equal(t, got, true)

	// Case: Returning: FALSE, Whitelisted: TRUE, Permissions: FALSE
	// Expected Result: FALSE
	wl.EXPECT().
		CheckDenomPermissions(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3", []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT}).
		Return(false)
	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3").
		Return(true)
	got = ibctransfer.IsRecvPacketAllowed(ctx, wl, nonReturningTransferPacket, whitelistedDenom)
	require.Equal(t, got, false)
	// Case: Returning: FALSE, Whitelisted: FALSE, Permissions: FALSE
	// Expected Result: FALSE
	wl.EXPECT().
		CheckDenomPermissions(ctx,
			"ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2", []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT}).
		Return(false)
	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2").
		Return(false)
	got = ibctransfer.IsRecvPacketAllowed(ctx, wl, nonReturningTransferPacket, disallowedDenom)
	require.Equal(t, got, false)
	// Case: Returning: TRUE, Whitelisted: FALSE, Permissions: FALSE
	// Expected Result: FALSE
	wl.EXPECT().
		CheckDenomPermissions(ctx,
			"rowan", []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT}).
		Return(false)
	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"rowan").
		Return(true)
	got = ibctransfer.IsRecvPacketAllowed(ctx, wl, returningTransferPacket, returningDenom)
	require.Equal(t, got, false)
	// Case: Returning: TRUE, Whitelisted: True, Permissions: FALSE
	// Expected Result: TRUE
	wl.EXPECT().
		CheckDenomPermissions(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3", []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT}).
		Return(false)
	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3").
		Return(true)
	got = ibctransfer.IsRecvPacketAllowed(ctx, wl, returningTransferPacket, whitelistedDenom)
	require.Equal(t, got, false)
}

func TestIsRecvPacketReturning(t *testing.T) {
	packet := channeltypes.Packet{
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
	}

	returningData := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-0/atom",
	}

	nonReturningData := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-11/atom",
	}

	got := ibctransfer.IsRecvPacketReturning(packet, returningData)
	require.Equal(t, got, true)

	got = ibctransfer.IsRecvPacketReturning(packet, nonReturningData)
	require.Equal(t, got, false)
}
