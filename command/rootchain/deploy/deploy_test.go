package deploy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/testutil"

	"github.com/vishnushankarsg/metad/command"
	"github.com/vishnushankarsg/metad/command/rootchain/helper"
	"github.com/vishnushankarsg/metad/consensus/polybft"
	"github.com/vishnushankarsg/metad/types"
)

func TestDeployContracts_NoPanics(t *testing.T) {
	t.Parallel()

	server := testutil.DeployTestServer(t, nil)
	t.Cleanup(func() {
		if err := os.RemoveAll(params.genesisPath); err != nil {
			t.Fatal(err)
		}
	})

	client, err := jsonrpc.NewClient(server.HTTPAddr())
	require.NoError(t, err)

	testKey, err := helper.GetRootchainPrivateKey("")
	require.NoError(t, err)

	receipt, err := server.Fund(testKey.Address())
	require.NoError(t, err)
	require.Equal(t, uint64(types.ReceiptSuccess), receipt.Status)

	outputter := command.InitializeOutputter(GetCommand())

	require.NotPanics(t, func() {
		_, err = deployContracts(outputter, client, 10, []*polybft.Validator{})
	})
	require.NoError(t, err)
}
