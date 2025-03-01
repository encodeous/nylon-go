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
		centralPath := cmd.Flag("config").Value.String()
		nodePath := cmd.Flag("node").Value.String()
	start:
		var centralCfg state.CentralCfg
		file, err := os.ReadFile(centralPath)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(file, &centralCfg)
		if err != nil {
			panic(err)
		}

		var nodeCfg state.LocalCfg
		file, err = os.ReadFile(nodePath)
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

		restart, err := core.Start(centralCfg, nodeCfg, level, centralPath)
		if err != nil {
			panic(err)
		}
		if restart {
			goto start
		}
	},
	GroupID: "ny",
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	runCmd.Flags().BoolVarP(&state.DBG_log_probe, "lprobe", "p", false, "Write probes to console")
	runCmd.Flags().BoolVarP(&state.DBG_log_router, "lroute", "r", false, "Write router updates to console")
	runCmd.Flags().BoolVarP(&state.DBG_log_wireguard, "lwg", "w", false, "Outputs wireguard logs to the console")
	runCmd.Flags().BoolVarP(&state.DBG_log_route_table, "ltable", "t", false, "Outputs route table to the console")
	runCmd.Flags().BoolVarP(&state.DBG_log_route_changes, "lrchange", "g", false, "Outputs route changes to the console")
	runCmd.Flags().BoolVarP(&state.DBG_log_repo_updates, "lrepo", "", false, "Outputs repo updates to the console")
	runCmd.Flags().StringP("config", "c", DefaultConfigPath, "Path to the config file")
	runCmd.Flags().StringP("node", "n", DefaultNodeConfigPath, "Path to the node config file")
}
