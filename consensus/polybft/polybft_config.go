package polybft

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/vishnushankarsg/metad/chain"
	"github.com/vishnushankarsg/metad/consensus/polybft/contractsapi"
	bls "github.com/vishnushankarsg/metad/consensus/polybft/signer"
	"github.com/vishnushankarsg/metad/helper/common"
	"github.com/vishnushankarsg/metad/types"
)

const ConsensusName = "polybft"

// PolyBFTConfig is the configuration file for the Polybft consensus protocol.
type PolyBFTConfig struct {
	// InitialValidatorSet are the genesis validators
	InitialValidatorSet []*Validator `json:"initialValidatorSet"`

	// Bridge is the rootchain bridge configuration
	Bridge *BridgeConfig `json:"bridge"`

	// EpochSize is size of epoch
	EpochSize uint64 `json:"epochSize"`

	// EpochReward is assigned to validators for blocks sealing
	EpochReward uint64 `json:"epochReward"`

	// SprintSize is size of sprint
	SprintSize uint64 `json:"sprintSize"`

	// BlockTime is target frequency of blocks production
	BlockTime common.Duration `json:"blockTime"`

	// Governance is the initial governance address
	Governance types.Address `json:"governance"`

	// MintableNativeToken denotes whether mintable native token is used
	MintableNativeToken bool `json:"mintableNative"`

	// NativeTokenConfig defines name, symbol and decimal count of the native token
	NativeTokenConfig *TokenConfig `json:"nativeTokenConfig"`

	InitialTrieRoot types.Hash `json:"initialTrieRoot"`
}

// LoadPolyBFTConfig loads chain config from provided path and unmarshals PolyBFTConfig
func LoadPolyBFTConfig(chainConfigFile string) (PolyBFTConfig, error) {
	chainCfg, err := chain.ImportFromFile(chainConfigFile)
	if err != nil {
		return PolyBFTConfig{}, err
	}

	return GetPolyBFTConfig(chainCfg)
}

// GetPolyBFTConfig deserializes provided chain config and returns PolyBFTConfig
func GetPolyBFTConfig(chainConfig *chain.Chain) (PolyBFTConfig, error) {
	consensusConfigJSON, err := json.Marshal(chainConfig.Params.Engine[ConsensusName])
	if err != nil {
		return PolyBFTConfig{}, err
	}

	var polyBFTConfig PolyBFTConfig
	if err = json.Unmarshal(consensusConfigJSON, &polyBFTConfig); err != nil {
		return PolyBFTConfig{}, err
	}

	return polyBFTConfig, nil
}

// BridgeConfig is the rootchain configuration, needed for bridging
type BridgeConfig struct {
	StateSenderAddr          types.Address `json:"stateSenderAddress"`
	CheckpointManagerAddr    types.Address `json:"checkpointManagerAddress"`
	ExitHelperAddr           types.Address `json:"exitHelperAddress"`
	RootERC20PredicateAddr   types.Address `json:"erc20PredicateAddress"`
	RootNativeERC20Addr      types.Address `json:"nativeERC20Address"`
	RootERC721PredicateAddr  types.Address `json:"erc721PredicateAddress"`
	RootERC1155PredicateAddr types.Address `json:"erc1155PredicateAddress"`

	JSONRPCEndpoint         string                   `json:"jsonRPCEndpoint"`
	EventTrackerStartBlocks map[types.Address]uint64 `json:"eventTrackerStartBlocks"`
}

func (p *PolyBFTConfig) IsBridgeEnabled() bool {
	return p.Bridge != nil
}

// Validator represents public information about validator accounts which are the part of genesis
type Validator struct {
	Address       types.Address
	BlsPrivateKey *bls.PrivateKey
	BlsKey        string
	BlsSignature  string
	Balance       *big.Int
	Stake         *big.Int
	MultiAddr     string
}

type validatorRaw struct {
	Address      types.Address `json:"address"`
	BlsKey       string        `json:"blsKey"`
	BlsSignature string        `json:"blsSignature"`
	Balance      *string       `json:"balance"`
	Stake        *string       `json:"stake"`
	MultiAddr    string        `json:"multiAddr"`
}

func (v *Validator) MarshalJSON() ([]byte, error) {
	raw := &validatorRaw{Address: v.Address, BlsKey: v.BlsKey, MultiAddr: v.MultiAddr, BlsSignature: v.BlsSignature}
	raw.Balance = types.EncodeBigInt(v.Balance)
	raw.Stake = types.EncodeBigInt(v.Stake)

	return json.Marshal(raw)
}

func (v *Validator) UnmarshalJSON(data []byte) error {
	var (
		raw validatorRaw
		err error
	)

	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	v.Address = raw.Address
	v.BlsKey = raw.BlsKey
	v.BlsSignature = raw.BlsSignature
	v.MultiAddr = raw.MultiAddr

	v.Balance, err = types.ParseUint256orHex(raw.Balance)
	if err != nil {
		return err
	}

	v.Stake, err = types.ParseUint256orHex(raw.Stake)
	if err != nil {
		return err
	}

	return nil
}

// UnmarshalBLSPublicKey unmarshals the hex encoded BLS public key
func (v *Validator) UnmarshalBLSPublicKey() (*bls.PublicKey, error) {
	decoded, err := hex.DecodeString(v.BlsKey)
	if err != nil {
		return nil, err
	}

	return bls.UnmarshalPublicKey(decoded)
}

// UnmarshalBLSSignature unmarshals the hex encoded BLS signature
func (v *Validator) UnmarshalBLSSignature() (*bls.Signature, error) {
	decoded, err := hex.DecodeString(v.BlsSignature)
	if err != nil {
		return nil, err
	}

	return bls.UnmarshalSignature(decoded)
}

// ToValidatorInitAPIBinding converts Validator to instance of contractsapi.ValidatorInit
func (v Validator) ToValidatorInitAPIBinding() (*contractsapi.ValidatorInit, error) {
	blsSignature, err := v.UnmarshalBLSSignature()
	if err != nil {
		return nil, err
	}

	signBigInts, err := blsSignature.ToBigInt()
	if err != nil {
		return nil, err
	}

	pubKey, err := v.UnmarshalBLSPublicKey()
	if err != nil {
		return nil, err
	}

	return &contractsapi.ValidatorInit{
		Addr:      v.Address,
		Pubkey:    pubKey.ToBigInt(),
		Signature: signBigInts,
		Stake:     new(big.Int).Set(v.Stake),
	}, nil
}

// ToValidatorMetadata creates ValidatorMetadata instance
func (v *Validator) ToValidatorMetadata() (*ValidatorMetadata, error) {
	blsKey, err := v.UnmarshalBLSPublicKey()
	if err != nil {
		return nil, err
	}

	metadata := &ValidatorMetadata{
		Address:     v.Address,
		BlsKey:      blsKey,
		VotingPower: new(big.Int).Set(v.Stake),
		IsActive:    true,
	}

	return metadata, nil
}

// String implements fmt.Stringer interface
func (v *Validator) String() string {
	return fmt.Sprintf("Address=%s; Balance=%d; P2P Multi addr=%s; BLS Key=%s;",
		v.Address, v.Balance, v.MultiAddr, v.BlsKey)
}

// RootchainConfig contains rootchain metadata (such as JSON RPC endpoint and contract addresses)
type RootchainConfig struct {
	JSONRPCAddr string

	StateSenderAddress          types.Address
	CheckpointManagerAddress    types.Address
	BLSAddress                  types.Address
	BN256G2Address              types.Address
	ExitHelperAddress           types.Address
	RootERC20PredicateAddress   types.Address
	RootNativeERC20Address      types.Address
	ERC20TemplateAddress        types.Address
	RootERC721PredicateAddress  types.Address
	RootERC721Address           types.Address
	RootERC1155PredicateAddress types.Address
	RootERC1155Address          types.Address
}

// ToBridgeConfig creates BridgeConfig instance
func (r *RootchainConfig) ToBridgeConfig() *BridgeConfig {
	return &BridgeConfig{
		JSONRPCEndpoint: r.JSONRPCAddr,

		StateSenderAddr:          r.StateSenderAddress,
		CheckpointManagerAddr:    r.CheckpointManagerAddress,
		ExitHelperAddr:           r.ExitHelperAddress,
		RootERC20PredicateAddr:   r.RootERC20PredicateAddress,
		RootNativeERC20Addr:      r.RootNativeERC20Address,
		RootERC721PredicateAddr:  r.RootERC721PredicateAddress,
		RootERC1155PredicateAddr: r.RootERC1155PredicateAddress,
	}
}

// TokenConfig is the configuration of native token used by  network
type TokenConfig struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint8  `json:"decimals"`
}
