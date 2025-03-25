package core

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type RootConfig struct {
	Server ServerConfig  `yaml:"server"`
	Chains []ChainConfig `yaml:"chains"`
}

type ServerConfig struct {
	AllowCors bool `yaml:"allow_cors"`
	Http      struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
	Grpc struct {
		Port int `yaml:"port"`
	} `yaml:"grpc"`
	Limit struct {
		Enabled bool  `yaml:"enabled"`
		Period  int64 `yaml:"period"`
	} `yaml:"limit"`
}

type ChainConfig struct {
	Name          string  `yaml:"name"`
	ChainId       string  `yaml:"chain_id"`
	RpcEndpoint   string  `yaml:"rpc_endpoint"`
	AccountPrefix string  `yaml:"account_prefix"`
	GasAdjustment float64 `yaml:"gas_adjustment"`
	GasPrice      string  `yaml:"gas_price"`
	Key           string  `yaml:"key"`
	DropCoin      string  `yaml:"drop_coin"`
	DropMaxCoin string `yaml:"drop_max_coin"`
}

func (cc *ChainConfig) KeyName() string {
	return fmt.Sprintf("__cosmfaucet_%s_root", cc.ChainId)
}

// LoadConfig loads config from file
func LoadConfig(filepath string) *RootConfig {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read file from path: %s, err: %v", filepath, err)
	}

	var config RootConfig
	if err := yaml.Unmarshal(contents, &config); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return &config
}
