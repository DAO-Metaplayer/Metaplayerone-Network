package server

import (
	"github.com/vishnushankarsg/metad/chain"
	"github.com/vishnushankarsg/metad/consensus"
	consensusDev "github.com/vishnushankarsg/metad/consensus/dev"
	consensusDummy "github.com/vishnushankarsg/metad/consensus/dummy"
	consensusIBFT "github.com/vishnushankarsg/metad/consensus/ibft"
	consensusPolyBFT "github.com/vishnushankarsg/metad/consensus/polybft"
	"github.com/vishnushankarsg/metad/secrets"
	"github.com/vishnushankarsg/metad/secrets/awsssm"
	"github.com/vishnushankarsg/metad/secrets/gcpssm"
	"github.com/vishnushankarsg/metad/secrets/hashicorpvault"
	"github.com/vishnushankarsg/metad/secrets/local"
	"github.com/vishnushankarsg/metad/state"
)

type GenesisFactoryHook func(config *chain.Chain, engineName string) func(*state.Transition) error

type ConsensusType string

const (
	DevConsensus     ConsensusType = "dev"
	IBFTConsensus    ConsensusType = "ibft"
	PolyBFTConsensus ConsensusType = consensusPolyBFT.ConsensusName
	DummyConsensus   ConsensusType = "dummy"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	DevConsensus:     consensusDev.Factory,
	IBFTConsensus:    consensusIBFT.Factory,
	PolyBFTConsensus: consensusPolyBFT.Factory,
	DummyConsensus:   consensusDummy.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
	secrets.GCPSSM:         gcpssm.SecretsManagerFactory,
}

var genesisCreationFactory = map[ConsensusType]GenesisFactoryHook{
	PolyBFTConsensus: consensusPolyBFT.GenesisPostHookFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
