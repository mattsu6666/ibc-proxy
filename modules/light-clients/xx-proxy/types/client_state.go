package types

import (
	"errors"

	ics23 "github.com/confio/ics23/go"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/modules/core/exported"
)

const ProxyClientType string = "proxyclient"

var _ exported.ClientState = (*ClientState)(nil)
var _ codectypes.UnpackInterfacesMessage = (*ClientState)(nil)

func NewClientState(upstreamClientID string, proxyClientState *codectypes.Any) *ClientState {
	return &ClientState{UpstreamClientId: upstreamClientID, ProxyClientState: proxyClientState}
}

func (cs *ClientState) IsInitialized() bool {
	return cs.ProxyClientState != nil
}

func (cs *ClientState) ClientType() string {
	return ProxyClientType
}

func (cs *ClientState) GetProxyClientState() exported.ClientState {
	state, err := clienttypes.UnpackClientState(cs.ProxyClientState)
	if err != nil {
		panic(err)
	}
	return state
}

func (cs *ClientState) GetUpstreamClientState() exported.ClientState {
	state, err := clienttypes.UnpackClientState(cs.UpstreamClientState)
	if err != nil {
		panic(err)
	}
	return state
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (cs *ClientState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := unpacker.UnpackAny(cs.ProxyClientState, new(exported.ClientState)); err != nil {
		return err
	}
	if err := unpacker.UnpackAny(cs.UpstreamClientState, new(exported.ClientState)); err != nil {
		return err
	}
	return nil
}

func (cs *ClientState) GetLatestHeight() exported.Height {
	if cs.ProxyClientState == nil {
		return clienttypes.NewHeight(0, 0)
	}
	return cs.GetProxyClientState().GetLatestHeight()
}

func (cs *ClientState) Status(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryCodec,
) exported.Status {
	return cs.GetProxyClientState().Status(ctx, clientStore, cdc)
}

func (cs *ClientState) Validate() error {
	return cs.GetProxyClientState().Validate()
}

func (cs *ClientState) GetProofSpecs() []*ics23.ProofSpec {
	return cs.GetProxyClientState().GetProofSpecs()
}

// Initialization function
// Clients must validate the initial consensus state, and may store any client-specific metadata
// necessary for correct light client operation
func (cs *ClientState) Initialize(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore, consState exported.ConsensusState) error {
	if cs.ProxyClientState != nil || cs.Prefix != nil || cs.UpstreamClientState != nil || cs.UpstreamConsensusState != nil {
		return sdkerrors.Wrap(errors.New("invalid clientState"), "each fields of the clientState must be empty")
	} else if consState != nil {
		return sdkerrors.Wrap(errors.New("invalid consensusState"), "each fields of the consensusState must be empty")
	}
	return nil
}

// Genesis function
func (cs *ClientState) ExportMetadata(_ sdk.KVStore) []exported.GenesisMetadata {
	panic("not implemented") // TODO: Implement
}

// Upgrade functions
// NOTE: proof heights are not included as upgrade to a new revision is expected to pass only on the last
// height committed by the current revision. Clients are responsible for ensuring that the planned last
// height of the current revision is somehow encoded in the proof verification process.
// This is to ensure that no premature upgrades occur, since upgrade plans committed to by the counterparty
// may be cancelled or modified before the last planned height.
func (cs *ClientState) VerifyUpgradeAndUpdateState(ctx sdk.Context, cdc codec.BinaryCodec, store sdk.KVStore, newClient exported.ClientState, newConsState exported.ConsensusState, proofUpgradeClient []byte, proofUpgradeConsState []byte) (exported.ClientState, exported.ConsensusState, error) {
	panic("not implemented") // TODO: Implement
}

// Utility function that zeroes out any client customizable fields in client state
// Ledger enforced fields are maintained while all custom fields are zero values
// Used to verify upgrades
func (cs *ClientState) ZeroCustomFields() exported.ClientState {
	panic("not implemented") // TODO: Implement
}

// State verification functions
func (cs *ClientState) VerifyClientState(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height, prefix exported.Prefix, counterpartyClientIdentifier string, proof []byte, clientState exported.ClientState) error {
	return cs.GetProxyClientState().VerifyClientState(NewProxyStore(cdc, store), cdc, height, newPrefix(cs.Prefix, prefix, cs.UpstreamClientId), counterpartyClientIdentifier, proof, clientState)
}

func (cs *ClientState) VerifyClientConsensusState(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height, counterpartyClientIdentifier string, consensusHeight exported.Height, prefix exported.Prefix, proof []byte, consensusState exported.ConsensusState) error {
	return cs.GetProxyClientState().VerifyClientConsensusState(NewProxyStore(cdc, store), cdc, height, counterpartyClientIdentifier, consensusHeight, newPrefix(cs.Prefix, prefix, cs.UpstreamClientId), proof, consensusState)
}

func (cs *ClientState) VerifyConnectionState(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height, prefix exported.Prefix, proof []byte, connectionID string, connectionEnd exported.ConnectionI) error {
	return cs.GetProxyClientState().VerifyConnectionState(NewProxyStore(cdc, store), cdc, height, newPrefix(cs.Prefix, prefix, cs.UpstreamClientId), proof, connectionID, connectionEnd)
}

func (cs *ClientState) VerifyChannelState(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height, prefix exported.Prefix, proof []byte, portID string, channelID string, channel exported.ChannelI) error {
	return cs.GetProxyClientState().VerifyChannelState(NewProxyStore(cdc, store), cdc, height, newPrefix(cs.Prefix, prefix, cs.UpstreamClientId), proof, portID, channelID, channel)
}

func (cs *ClientState) VerifyPacketCommitment(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height, currentTimestamp uint64, delayPeriod uint64, prefix exported.Prefix, proof []byte, portID string, channelID string, sequence uint64, commitmentBytes []byte) error {
	return cs.GetProxyClientState().VerifyPacketCommitment(NewProxyStore(cdc, store), cdc, height, currentTimestamp, delayPeriod, newPrefix(cs.Prefix, prefix, cs.UpstreamClientId), proof, portID, channelID, sequence, commitmentBytes)
}

func (cs *ClientState) VerifyPacketAcknowledgement(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height, currentTimestamp uint64, delayPeriod uint64, prefix exported.Prefix, proof []byte, portID string, channelID string, sequence uint64, acknowledgement []byte) error {
	return cs.GetProxyClientState().VerifyPacketAcknowledgement(NewProxyStore(cdc, store), cdc, height, currentTimestamp, delayPeriod, newPrefix(cs.Prefix, prefix, cs.UpstreamClientId), proof, portID, channelID, sequence, acknowledgement)
}

func (cs *ClientState) VerifyPacketReceiptAbsence(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height, currentTimestamp uint64, delayPeriod uint64, prefix exported.Prefix, proof []byte, portID string, channelID string, sequence uint64) error {
	return cs.GetProxyClientState().VerifyPacketReceiptAbsence(NewProxyStore(cdc, store), cdc, height, currentTimestamp, delayPeriod, newPrefix(cs.Prefix, prefix, cs.UpstreamClientId), proof, portID, channelID, sequence)
}

func (cs *ClientState) VerifyNextSequenceRecv(store sdk.KVStore, cdc codec.BinaryCodec, height exported.Height, currentTimestamp uint64, delayPeriod uint64, prefix exported.Prefix, proof []byte, portID string, channelID string, nextSequenceRecv uint64) error {
	return cs.GetProxyClientState().VerifyNextSequenceRecv(NewProxyStore(cdc, store), cdc, height, currentTimestamp, delayPeriod, newPrefix(cs.Prefix, prefix, cs.UpstreamClientId), proof, portID, channelID, nextSequenceRecv)
}

func newPrefix(proxyPrefix, upstreamPrefix exported.Prefix, upstreamClientID string) exported.Prefix {
	return commitmenttypes.MultiPrefix{
		Prefix:     proxyPrefix,
		PathPrefix: append([]byte(upstreamClientID+"/"), upstreamPrefix.Bytes()...),
	}
}