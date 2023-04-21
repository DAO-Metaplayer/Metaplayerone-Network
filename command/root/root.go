package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vishnushankarsg/metad/command/backup"
	"github.com/vishnushankarsg/metad/command/bridge"
	"github.com/vishnushankarsg/metad/command/genesis"
	"github.com/vishnushankarsg/metad/command/helper"
	"github.com/vishnushankarsg/metad/command/license"
	"github.com/vishnushankarsg/metad/command/monitor"
	"github.com/vishnushankarsg/metad/command/peers"
	"github.com/vishnushankarsg/metad/command/polybft"
	"github.com/vishnushankarsg/metad/command/polybftsecrets"
	"github.com/vishnushankarsg/metad/command/regenesis"
	"github.com/vishnushankarsg/metad/command/rootchain"
	"github.com/vishnushankarsg/metad/command/secrets"
	"github.com/vishnushankarsg/metad/command/server"
	"github.com/vishnushankarsg/metad/command/status"
	"github.com/vishnushankarsg/metad/command/txpool"
	"github.com/vishnushankarsg/metad/command/version"
	"github.com/vishnushankarsg/metad/command/whitelist"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "Metachain  is a framework for building Ethereum-compatible Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		rootchain.GetCommand(),
		monitor.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		whitelist.GetCommand(),
		license.GetCommand(),
		polybftsecrets.GetCommand(),
		polybft.GetCommand(),
		bridge.GetCommand(),
		regenesis.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
