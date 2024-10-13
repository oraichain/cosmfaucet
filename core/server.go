package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"strings"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/scalalang2/cosmfaucet/gen/proto/faucetpb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	mux     sync.Mutex
	log     *zap.Logger
	limiter *Limiter

	faucetpb.FaucetServiceServer
	config  *RootConfig
	clients ChainClients
	faucet  *Faucet
}

func NewServer(log *zap.Logger, config *RootConfig, clients ChainClients) *Server {
	chains := make([]ChainId, 0)
	for _, chainCfg := range config.Chains {
		chains = append(chains, chainCfg.ChainId)
	}

	var limiter *Limiter
	if config.Server.Limit.Enabled {
		limiter = NewLimiter(chains, config.Server.Limit.Period)
	}

	faucet := NewFaucet(log, clients, 100)
	faucet.run()

	return &Server{
		log:     log,
		limiter: limiter,
		config:  config,
		clients: clients,
		faucet:  faucet,
	}
}

// RemoteAddr returns the remote address of the request
func RemoteAddr(ctx context.Context) string {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		xForwardFor := headers.Get("x-forwarded-for")
		if len(xForwardFor) > 0 && xForwardFor[0] != "" {
			ips := strings.Split(xForwardFor[0], ",")
			if len(ips) > 0 {
				clientIp := ips[0]
				return clientIp
			}
		}
	}
	return ""
}

// GiveMe sends a `BankMsg` transaction to the chain to send some tokens to the given address
// It blocks the request if the user is given the token in the last 24 hours.
func (s *Server) GiveMe(ctx context.Context, request *faucetpb.GiveMeRequest) (*faucetpb.GiveMeResponse, error) {
	remoteAddr := RemoteAddr(ctx)
	client, ok := s.clients[request.ChainId]
	if !ok {
		return nil, status.Error(codes.NotFound, "chain not supported")
	}

	// find config from RootConfig
	var chainConfig *ChainConfig
	for _, chain := range s.config.Chains {
		if chain.ChainId == request.ChainId {
			chainConfig = &chain
			break
		}
	}

	var address []byte
	var err error
	// Here we accept both an Ethereum-formatted address or a cosmos-sdk bech32 address
	// Each address type maps to the same underlying address that will holding a single
	// balance fungible between the base cosmos chain and the Ethereum VM
	if common.IsHexAddress(request.Address) {
		address = common.HexToAddress(request.Address).Bytes()
	} else {
		// validate address format
		address, err = sdk.GetFromBech32(request.Address, chainConfig.AccountPrefix)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid address")
		}
	}

	ki, err := client.Keybase.Key(chainConfig.KeyName())
	if err != nil {
		return nil, status.Error(codes.Internal, "could not get faucet key")
	}

	coin, err := sdk.ParseCoinNormalized(chainConfig.DropCoin)
	if err != nil {
		s.log.Error("invalid coin format", zap.Error(err))
		return nil, status.Error(codes.Internal, "invalid coin format | this is unexpected error, please inform to the admin.")
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	if s.limiter != nil {
		if !s.limiter.IsAllowed(request.ChainId, remoteAddr) {
			return nil, status.Error(codes.PermissionDenied, "user cannot request token more than once during specific period of time")
		}
	}

	s.faucet.sendTask(request.ChainId, &work{
		chainId: request.ChainId,
		detail: &transferWork{
			fromAddress: client.MustEncodeAccAddr(ki.GetAddress()),
			toAddress:   client.MustEncodeAccAddr(address),
			amount:      []sdk.Coin{coin},
		},
	})

	if s.limiter != nil {
		s.limiter.AddRequest(request.ChainId, remoteAddr)
	}

	return &faucetpb.GiveMeResponse{}, nil
}

// Chains returns all supported chains
func (s *Server) Chains(ctx context.Context, request *faucetpb.GetChainsRequest) (*faucetpb.GetChainsResponse, error) {
	res := make([]*faucetpb.Chain, 0)
	for _, chain := range s.config.Chains {
		res = append(res, &faucetpb.Chain{
			Name:    chain.Name,
			ChainId: chain.ChainId,
		})
	}

	return &faucetpb.GetChainsResponse{
		Chains: res,
	}, nil
}
