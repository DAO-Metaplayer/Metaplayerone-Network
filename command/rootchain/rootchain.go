package rootchain

import (
	"github.com/spf13/cobra"

	"github.com/vishnushankarsg/metad/command/rootchain/deploy"
)

// GetCommand creates "rootchain" helper command
func GetCommand() *cobra.Command {
	rootchainCmd := &cobra.Command{
		Use:   "rootchain",
		Short: "Top level rootchain helper command.",
	}

	rootchainCmd.AddCommand(
		// rootchain deploy
		deploy.GetCommand(),
	)

	return rootchainCmd
}
