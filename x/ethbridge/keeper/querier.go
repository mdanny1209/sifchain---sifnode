package keeper

import (
	"fmt"
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

// TODO: move to x/oracle

// NewLegacyQuerier is the module level router for state queries
func NewLegacyQuerier(keeper Keeper, cdc *codec.LegacyAmino) sdk.Querier {

	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryEthProphecy:
			return legacyQueryEthProphecy(ctx, cdc, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown ethbridge query endpoint")
		}
	}
}

func legacyQueryEthProphecy(ctx sdk.Context, cdc *codec.LegacyAmino, query abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var req types.QueryEthProphecyRequest
	
	if err := cdc.UnmarshalJSON(query.Data, &req); err != nil {
		return nil, sdkerrors.Wrap(types.ErrJSONMarshalling, fmt.Sprintf("failed to parse req: %s", err.Error()))
	}

	id := strconv.FormatInt(req.EthereumChainId, 10) + strconv.FormatInt(req.Nonce, 10) + req.EthereumSender

	prophecy, found := keeper.oracleKeeper.GetProphecy(ctx, id)
	if !found {
		return nil, sdkerrors.Wrap(oracletypes.ErrProphecyNotFound, id)
	}

	bridgeClaims, err := types.MapOracleClaimsToEthBridgeClaims(
		req.EthereumChainId,
		types.NewEthereumAddress(req.BridgeContractAddress),
		req.Nonce,
		req.Symbol,
		types.NewEthereumAddress(req.TokenContractAddress),
		types.NewEthereumAddress(req.EthereumSender),
		prophecy.ValidatorClaims,
		types.CreateEthClaimFromOracleString,
	)
	if err != nil {
		return nil, err
	}

	response := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	return cdc.MarshalJSONIndent(response, "", "  ")
}
