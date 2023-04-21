package polybft

import (
	"github.com/vishnushankarsg/metad/command/sidechain/registration"
	"github.com/vishnushankarsg/metad/command/sidechain/staking"
	"github.com/vishnushankarsg/metad/command/sidechain/unstaking"
	"github.com/vishnushankarsg/metad/command/sidechain/validators"

	"github.com/vishnushankarsg/metad/command/sidechain/whitelist"
	"github.com/vishnushankarsg/metad/command/sidechain/withdraw"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	polybftCmd := &cobra.Command{
		Use:   "validator",
		Short: "validator command",
	}

	polybftCmd.AddCommand(
		staking.GetCommand(),
		unstaking.GetCommand(),
		withdraw.GetCommand(),
		validators.GetCommand(),
		whitelist.GetCommand(),
		registration.GetCommand(),
	)

	return polybftCmd
}
