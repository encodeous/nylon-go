package cmd

import (
	"github.com/encodeous/nylon/core"
	"github.com/encodeous/nylon/state"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run nylon",
	Long:  `This will run nylon on the current host. Ensure it has enough permissions to interact with in-kernel WireGuard.`,
	Run: func(cmd *cobra.Command, args []string) {
		var centralCfg state.CentralCfg
		file, err := os.ReadFile(centralConfigPath)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(file, &centralCfg)
		if err != nil {
			panic(err)
		}

		var nodeCfg state.NodeCfg
		file, err = os.ReadFile(nodeConfigPath)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(file, &nodeCfg)
		if err != nil {
			panic(err)
		}

		err = state.CentralConfigValidator(&centralCfg)
		if err != nil {
			panic(err)
		}
		err = state.NodeConfigValidator(&nodeCfg)
		if err != nil {
			panic(err)
		}

		level := slog.LevelInfo
		if ok, _ := cmd.Flags().GetBool("verbose"); ok {
			level = slog.LevelDebug
		}

		err = core.Start(centralCfg, nodeCfg, level)
		if err != nil {
			panic(err)
		}
	},
	GroupID: "ny",
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
