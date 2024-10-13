package core

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"os"
	"strings"

	lens "github.com/strangelove-ventures/lens/client"
	"go.uber.org/zap"
)

// chain initializer
var newChainClient = newLensClient

// newLensClient create a client for the cosmos blockchain
func newLensClient(logger *zap.Logger, config ChainConfig, homePath string) (*lens.ChainClient, error) {
	if !strings.HasPrefix(config.RpcEndpoint, "http") {
		return nil, errInvalidEndpoint{rpc: config.RpcEndpoint}
	}

	keyName := config.KeyName()
	cfg := lens.ChainClientConfig{
		Key:            keyName,
		ChainID:        config.ChainId,
		RPCAddr:        config.RpcEndpoint,
		AccountPrefix:  config.AccountPrefix,
		KeyringBackend: "test",
		GasAdjustment:  config.GasAdjustment,
		GasPrices:      config.GasPrice,
		KeyDirectory:   homePath,
		Debug:          false,
		Timeout:        "20s",
		OutputFormat:   "json",
		SignModeStr:    "direct",
		Modules:        lens.ModuleBasics,
	}

	cc, err := lens.NewChainClient(logger, &cfg, homePath, os.Stdin, os.Stdout)
	if err != nil {
		return nil, err
	}

	// First delete the canonical key if it exists to prevent RestoreKey from throwing an error
	err = cc.Keybase.Delete(keyName)
	if err != nil && !sdkerrors.ErrKeyNotFound.Is(err) {
		return nil, err
	}

	addr, err := cc.RestoreKey(keyName, config.Key, 118)
	if err != nil {
		return nil, err
	}
	logger.Info("master wallet is restored", zap.String("address", addr))

	return cc, nil
}
