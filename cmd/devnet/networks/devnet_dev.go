package networks

import (
	"github.com/ledgerwatch/erigon-lib/chain/networkname"
	"github.com/ledgerwatch/erigon/cmd/devnet/accounts"
	"github.com/ledgerwatch/erigon/cmd/devnet/args"
	"github.com/ledgerwatch/erigon/cmd/devnet/devnet"
	account_services "github.com/ledgerwatch/erigon/cmd/devnet/services/accounts"
	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/log/v3"
)

func NewDevDevnet(
	dataDir string,
	baseRpcHost string,
	baseRpcPort int,
	producerCount int,
	logger log.Logger,
) devnet.Devnet {
	faucetSource := accounts.NewAccount("faucet-source")

	var nodes []devnet.Node

	if producerCount == 0 {
		producerCount++
	}

	for i := 0; i < producerCount; i++ {
		nodes = append(nodes, &args.BlockProducer{
			NodeArgs: args.NodeArgs{
				ConsoleVerbosity: "0",
				DirVerbosity:     "5",
			},
			AccountSlots: 200,
		})
	}

	network := devnet.Network{
		DataDir:            dataDir,
		Chain:              networkname.DevChainName,
		Logger:             logger,
		BasePrivateApiAddr: "localhost:10090",
		BaseRPCHost:        baseRpcHost,
		BaseRPCPort:        baseRpcPort,
		Alloc: types.GenesisAlloc{
			faucetSource.Address: {Balance: accounts.EtherAmount(200_000)},
		},
		Services: []devnet.Service{
			account_services.NewFaucet(networkname.DevChainName, faucetSource),
		},
		MaxNumberOfEmptyBlockChecks: 30,
		Nodes: append(nodes,
			&args.BlockConsumer{
				NodeArgs: args.NodeArgs{
					ConsoleVerbosity: "0",
					DirVerbosity:     "5",
				},
			}),
	}

	return devnet.Devnet{&network}
}
