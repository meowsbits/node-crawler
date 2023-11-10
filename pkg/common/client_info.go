package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
)

type ClientInfo struct {
	ClientType      string
	SoftwareVersion uint64
	Capabilities    []p2p.Cap
	NetworkID       uint64
	ForkID          forkid.ID
	Blockheight     string
	TotalDifficulty *big.Int
	HeadHash        common.Hash
}

type forkIDFilterer struct {
	name   string
	filter forkid.Filter
}

var forkIDFilters = []forkIDFilterer{
	{name: "mainnet", filter: forkid.NewStaticFilter(params.MainnetChainConfig, params.MainnetGenesisHash)},
	{name: "goerli", filter: forkid.NewStaticFilter(params.GoerliChainConfig, params.GoerliGenesisHash)},
	{name: "sepolia", filter: forkid.NewStaticFilter(params.SepoliaChainConfig, params.SepoliaGenesisHash)},
	{name: "classic",
		filter: forkid.NewStaticFilter(params.ClassicChainConfig, params.MainnetGenesisHash)},
	{name: "mordor",
		filter: forkid.NewStaticFilter(params.MordorChainConfig, params.MordorGenesisHash)},
}

func ForkIDName(fid forkid.ID) string {
	for _, target := range forkIDFilters {
		if target.filter(fid) == nil {
			// No error returned; match.
			return target.name
		}
	}
	return ""
}
