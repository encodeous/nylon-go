package cmd

import (
	"fmt"
	"github.com/encodeous/nylon/state"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"net/netip"
	"os"
)

// netCmd represents the net command
var netCmd = &cobra.Command{
	Use:   "new-net",
	Short: "Create a new nylon network with central configuration",
	Run: func(cmd *cobra.Command, args []string) {
		nodeCfg := promptCreateNode()

		ncfg, err := yaml.Marshal(&nodeCfg)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(nodeConfigPath, ncfg, 0700)
		if err != nil {
			panic(err)
		}

		pkey := state.GenerateKey()

		centralConfig := state.CentralCfg{
			RootPubKey: pkey.XPubkey(),
			Nodes: []state.PubNodeCfg{
				{
					Id: "sample node",
					Prefixes: []netip.Prefix{
						netip.MustParsePrefix("10.0.0.1/32"),
						netip.MustParsePrefix("10.0.0.2/32"),
						netip.MustParsePrefix("10.1.0.0/16"),
					},
					PubKey: state.NyPublicKey{},
					Endpoints: []netip.AddrPort{
						netip.MustParseAddrPort(fmt.Sprintf("8.8.8.8:%d", nodeCfg.Port)),
					},
				},
			},
			Graph: []string{
				"Group1 = node1, node2",
				"Group2 = node5, node6",
				"Group1, Group2, node7",
			},
			Version: 0,
		}

		fmt.Println("Where should the central config be saved?:")
		centralConfigPath = safeSaveFile(centralConfigPath, "Central Config")
		fmt.Println("Where should the central key be saved?:")
		centralKeyPath = safeSaveFile(centralKeyPath, "Central Key")
		ccfg, err := yaml.Marshal(&centralConfig)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(centralConfigPath, ccfg, 0700)
		if err != nil {
			panic(err)
		}

		key, err := pkey.MarshalText()
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(centralKeyPath, key, 0700)
		if err != nil {
			panic(err)
		}

	},
	GroupID: "init",
}

func init() {
	rootCmd.AddCommand(netCmd)
}
