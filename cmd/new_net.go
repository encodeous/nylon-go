package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"github.com/encodeous/nylon/state"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

		rootPubkey, rootPrivkey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}

		centralConfig := state.CentralCfg{
			RootPubKey: state.EdPrivateKey(rootPubkey),
			Nodes: []state.PubNodeCfg{
				promptGenPubCfg(nodeCfg),
			},
			Edges: []state.Pair[state.Node, state.Node]{
				{nodeCfg.Id, "other-node"},
			},
			Version: 0,
		}

		fmt.Println("\n\nCentral Network Configuration")

		centralConfigPath = safeSaveFile(centralConfigPath, "Central Config")
		centralKeyPath = safeSaveFile(centralKeyPath, "Central Key")
		ccfg, err := yaml.Marshal(&centralConfig)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(centralConfigPath, ccfg, 0700)
		if err != nil {
			panic(err)
		}

		key, err := state.EdPrivateKey(rootPrivkey).MarshalText()
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// netCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// netCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
