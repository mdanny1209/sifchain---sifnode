//nolint:dupl
package ethbridge

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
)

const errorMessageKey = "errorMessage"

// NewHandler returns a handler for "ethbridge" type messages.
func NewHandler(k Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgCreateEthBridgeClaim:
			res, err := msgServer.CreateEthBridgeClaim(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBurn:
			res, err := msgServer.Burn(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgLock:
			res, err := msgServer.Lock(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case MsgUpdateWhiteListValidator:
			return handleMsgUpdateWhiteListValidator(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		case MsgUpdateCethReceiverAccount:
			return handleMsgUpdateCethReceiverAccount(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		case MsgRescueCeth:
			return handleMsgRescueCeth(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized ethbridge message type: %v", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgUpdateWhiteListValidator(
	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
	bridgeKeeper Keeper, msg MsgUpdateWhiteListValidator,
) (*sdk.Result, error) {
	logger := bridgeKeeper.Logger(ctx)
	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender.String())

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	if err := bridgeKeeper.ProcessUpdateWhiteListValidator(ctx, msg.CosmosSender, msg.Validator, msg.OperationType); err != nil {
		logger.Error("bridge keeper failed to process update validator.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit update whitelist validators event.",
		"CosmosSender", msg.CosmosSender.String(),
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"Validator", msg.Validator.String(),
		"OperationType", msg.OperationType)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyOperationType, msg.OperationType),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUpdateCethReceiverAccount(
	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
	bridgeKeeper Keeper, msg MsgUpdateCethReceiverAccount,
) (*sdk.Result, error) {
	logger := bridgeKeeper.Logger(ctx)
	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender.String())

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	if err := bridgeKeeper.ProcessUpdateCethReceiverAccount(ctx, msg.CosmosSender, msg.CethReceiverAccount); err != nil {
		logger.Error("keeper failed to process update ceth receiver account.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit update ceth receiver account event.",
		"CosmosSender", msg.CosmosSender.String(),
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"CethReceiverAccount", msg.CethReceiverAccount.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(types.AttributeKeyCethReceiverAccount, msg.CethReceiverAccount.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRescueCeth(
	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
	bridgeKeeper Keeper, msg MsgRescueCeth) (*sdk.Result, error) {
	logger := bridgeKeeper.Logger(ctx)
	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender.String())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}
	if err := bridgeKeeper.ProcessRescueCeth(ctx, msg); err != nil {
		logger.Error("keeper failed to process rescue ceth message.", errorMessageKey, err.Error())
		return nil, err
	}
	logger.Info("sifnode emit rescue ceth event.",
		"CosmosSender", msg.CosmosSender.String(),
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"CosmosReceiver", msg.CosmosReceiver.String(),
		"CethAmount", msg.CethAmount.String())

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
