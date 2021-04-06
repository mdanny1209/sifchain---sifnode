package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sifchain/sifnode/x/clp/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryGetPool(c context.Context, req *types.PoolReq) (*types.PoolRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pool, err := k.GetPool(ctx, req.Symbol)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "validator %s not found", req.Symbol)
	}

	return &types.PoolRes{
		Pool:             &pool,
		Height:           ctx.BlockHeight(),
		ClpModuleAddress: types.GetCLPModuleAddress().String(),
	}, nil
}

func (k Keeper) QueryGetPools(c context.Context, req *types.PoolsReq) (*types.PoolsRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pool := k.GetPools(ctx)

	return &types.PoolsRes{
		Pools:            pool,
		Height:           ctx.BlockHeight(),
		ClpModuleAddress: types.GetCLPModuleAddress().String(),
	}, nil
}

func (k Keeper) LiquidityProvider(c context.Context, req *types.LiquidityProviderReq) (*types.LiquidityProviderRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	lp, err := k.GetLiquidityProvider(ctx, req.Symbol, req.LpAddress)
	if err != nil {
		return nil, err
	}
	pool, err := k.GetPool(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	native, external, _, _ := CalculateAllAssetsForLP(pool, lp)
	lpResponse := types.NewLiquidityProviderResponse(lp, ctx.BlockHeight(), native.String(), external.String())
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return &lpResponse, nil
}

func (k Keeper) GetAssetList(c context.Context, req *types.AssetListReq) (*types.AssetListRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	addr, err := sdk.AccAddressFromBech32(req.LpAddress)
	if err != nil {
		return nil, err
	}

	assetList := k.GetAssetsForLiquidityProvider(ctx, addr)

	var al []*types.Asset

	for _, asset := range assetList {
		al = append(al, &asset)
	}

	return &types.AssetListRes{
		Assets: al,
	}, nil
}

func (k Keeper) GetLiquidityProviderList(c context.Context, req *types.LiquidityProviderListReq) (*types.LiquidityProviderListRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	searchingAsset := types.NewAsset(req.Symbol)
	lpList := k.GetLiquidityProvidersForAsset(ctx, searchingAsset)

	var lpl []*types.LiquidityProvider
	for _, lp := range lpList {
		lpl = append(lpl, &lp)
	}
	return &types.LiquidityProviderListRes{
		LiquidityProviders: lpl,
		Height:             ctx.BlockHeight(),
	}, nil
}

func (k Keeper) QueryGetLiquidityProviders(c context.Context, req *types.LiquidityProvidersReq) (*types.LiquidityProvidersRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var lpl []*types.LiquidityProvider
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, types.LiquidityProviderPrefix)

	pageRes, err := query.FilteredPaginate(valStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var lp types.LiquidityProvider
		err := k.cdc.UnmarshalBinaryBare(value, &lp)
		if err != nil {
			return false, err
		}

		if accumulate {
			lpl = append(lpl, &lp)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.LiquidityProvidersRes{
		LiquidityProviders: lpl,
		Height:             ctx.BlockHeight(),
		Pagination:         pageRes,
	}, nil
}
