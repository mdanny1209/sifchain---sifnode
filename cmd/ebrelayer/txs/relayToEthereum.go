package txs

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
)

const (
	// GasLimit the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit            = uint64(500000)
	transactionInterval = 1 * time.Second
)

var GasPriceMinimum *big.Int = big.NewInt(120000000000)
var nonce uint64 = 0

func getNonce(address common.Address, client *ethclient.Client) (uint64, error) {
	if nonce > 0 {
		return nonce, nil
	}

	nonceValue, err := client.PendingNonceAt(context.Background(), address)
	if nonceValue > 0 {
		nonce = nonceValue
	}

	return nonce, err
}

func incrementNonce() {
	nonce++
}

// func recursiveRetry(provider string, contractAddress common.Address, event types.Event, claim ProphecyClaim,
// 	key *ecdsa.PrivateKey, sugaredLogger *zap.SugaredLogger, cosmosBridge *cosmosbridge.CosmosBridge,
// 	client *ethclient.Client, auth *bind.TransactOpts, retryNum uint16) error {
	
// 	if retryNum > 30 {
// 		err := fmt.Errorf("retried over %d times. stopping retry", retryNum)
// 		return err
// 	}

// 	_, err := cosmosBridge.NewProphecyClaim(auth, uint8(claim.ClaimType),
// 		claim.CosmosSender, claim.CosmosSenderSequence, claim.EthereumReceiver, claim.Symbol, claim.Amount.BigInt())
// 	if err != nil {
// 		return recursiveRetry(provider, contractAddress, event, claim, key, sugaredLogger, cosmosBridge, client, auth, retryNum + 1)
// 	}

// 	return nil
// }

// RelayProphecyClaimToEthereum relays the provided ProphecyClaim to CosmosBridge contract on the Ethereum network
func RelayProphecyClaimToEthereum(provider string, contractAddress common.Address, event types.Event,
	claim ProphecyClaim, key *ecdsa.PrivateKey, sugaredLogger *zap.SugaredLogger) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target, err := initRelayConfig(provider, contractAddress, event, key, sugaredLogger)
	if err != nil {
		sugaredLogger.Errorw("failed in init relay config.",
			errorMessageKey, err.Error())
		return err
	}

	// Initialize CosmosBridge instance
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmosBridge instance.",
			errorMessageKey, err.Error())
		return err
	}

	// Send transaction
	sugaredLogger.Infow("Sending new ProphecyClaim to CosmosBridge.",
		"CosmosSender", claim.CosmosSender,
		"CosmosSenderSequence", claim.CosmosSenderSequence)

	tx, err := cosmosBridgeInstance.NewProphecyClaim(auth, uint8(claim.ClaimType),
		claim.CosmosSender, claim.CosmosSenderSequence, claim.EthereumReceiver, claim.Symbol, claim.Amount.BigInt())

	if err.Error() == "not found" {
		// if there is a not found error while broadcasting to infura, retry until there isn't an error
		i := 0
		for i < 30 {
			sugaredLogger.Errorw(
				"could not find tx after broadcasting",
				"retry", i,
				"err", err.Error(),
			)

			tx, err = cosmosBridgeInstance.NewProphecyClaim(auth, uint8(claim.ClaimType),
				claim.CosmosSender, claim.CosmosSenderSequence, claim.EthereumReceiver, claim.Symbol, claim.Amount.BigInt())
			if err == nil {
				break
			}

			// sleep for 1 second until retrying
			time.Sleep(transactionInterval)
			i++
		}

		if i == 30 {
			sugaredLogger.Errorw("failed to send ProphecyClaim to CosmosBridge.",
				errorMessageKey, err.Error())
			return err
		}
	}

	sugaredLogger.Infow("get NewProphecyClaim tx hash:", "ProphecyClaimHash", tx.Hash().Hex())

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		i := 0
		for i < 30 {
			sugaredLogger.Errorw(
				"no tx receipt after broadcasting",
				"retry", i,
				"getTxReceiptErr", err.Error(), 
			)

			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err == nil {
				sugaredLogger.Infow(
					"Successfully received transaction receipt after retry",
					"txReceipt", receipt,
				)
				break
			}

			// sleep for 1 second until retrying
			time.Sleep(transactionInterval)
			i++
		}

		if i == 30 {
			sugaredLogger.Errorw("failed to get transaction receipt.",
				errorMessageKey, err.Error())
			return err
		}
	}

	// do not increment the nonce until the transaction receipt has been found
	incrementNonce()

	switch receipt.Status {
	case 0:
		sugaredLogger.Infow("trasaction failed:")
	case 1:
		sugaredLogger.Infow("trasaction is successful:")
	}
	return nil
}

// initRelayConfig set up Ethereum client, validator's transaction auth, and the target contract's address
func initRelayConfig(provider string, registry common.Address, event types.Event, key *ecdsa.PrivateKey,
	sugaredLogger *zap.SugaredLogger) (*ethclient.Client, *bind.TransactOpts, common.Address, error) {
	// Start Ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		sugaredLogger.Errorw("failed to connect ethereum node.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	// Load the validator's address
	sender, err := LoadSender()
	if err != nil {
		sugaredLogger.Errorw("failed to load validator address.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	nonce, err := getNonce(sender, client)
	
	sugaredLogger.Infow("Current eth operator at pending nonce.", "pendingNonce", nonce)

	if err != nil {
		sugaredLogger.Errorw("failed to broadcast transaction.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		sugaredLogger.Errorw("failed to get gas price.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth := bind.NewKeyedTransactor(key)

	sugaredLogger.Infow("ethereum tx current nonce from client api.",
		"nonce", nonce,
		"suggestedGasPrice", gasPrice)

	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))

	quarterGasPrice := big.NewInt(0)
	quarterGasPrice = quarterGasPrice.Div(gasPrice, big.NewInt(4))

	gasPrice.Sub(gasPrice, quarterGasPrice)
	if gasPrice.Cmp(GasPriceMinimum) == -1 {

		sugaredLogger.Errorw(
			"gas price under minimum of 120 gigawei",
			"gasPriceBeforeAdjustment", gasPrice,
			"gasPriceAfterAdjustment", GasPriceMinimum,
		)

		gasPrice = GasPriceMinimum
	}

	sugaredLogger.Infow("final gas price after adjustment.",
		"finalGasPrice", gasPrice)

	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

	sugaredLogger.Infow("nonce before send transaction.",
		"transactOptsAuth.Nonce", transactOptsAuth.Nonce)

	var targetContract ContractRegistry
	switch event {
	// ProphecyClaims are sent to the CosmosBridge contract
	case types.MsgBurn, types.MsgLock:
		targetContract = CosmosBridge
	// OracleClaims are sent to the Oracle contract
	case types.LogNewProphecyClaim:
		targetContract = Oracle
	default:
		panic("invalid target contract address")
	}

	// Get the specific contract's address
	target, err := GetAddressFromBridgeRegistry(client, registry, targetContract, sugaredLogger)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmos bridger contract address from registry.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err

	}
	return client, transactOptsAuth, target, nil
}
