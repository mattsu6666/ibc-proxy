package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/modules/core/exported"
)

var (
	_, _, _, _ sdk.Msg                            = (*MsgProxyClientState)(nil), (*MsgProxyConnectionOpenTry)(nil), (*MsgProxyConnectionOpenAck)(nil), (*MsgProxyConnectionOpenConfirm)(nil)
	_, _, _, _ codectypes.UnpackInterfacesMessage = (*MsgProxyClientState)(nil), (*MsgProxyConnectionOpenTry)(nil), (*MsgProxyConnectionOpenAck)(nil), (*MsgProxyConnectionOpenConfirm)(nil)

	_, _, _ sdk.Msg = (*MsgProxyChannelOpenTry)(nil), (*MsgProxyChannelOpenAck)(nil), (*MsgProxyChannelOpenConfirm)(nil)
	_, _    sdk.Msg = (*MsgProxyRecvPacket)(nil), (*MsgProxyAcknowledgePacket)(nil)
)

func NewMsgProxyClientState(
	upstreamClientID string,
	upstreamPrefix commitmenttypes.MerklePrefix,
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
	proofClient []byte,
	proofConsensus []byte,
	proofHeight clienttypes.Height,
	consensusHeight clienttypes.Height,
	signer string,
) (*MsgProxyClientState, error) {
	anyClient, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}
	anyConsensus, err := clienttypes.PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}

	return &MsgProxyClientState{
		UpstreamClientId: upstreamClientID,
		UpstreamPrefix:   upstreamPrefix,
		ClientState:      anyClient,
		ConsensusState:   anyConsensus,
		ProofClient:      proofClient,
		ProofConsensus:   proofConsensus,
		ProofHeight:      proofHeight,
		ConsensusHeight:  consensusHeight,
		Signer:           signer,
	}, nil
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyClientState) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgProxyClientState) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgProxyClientState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var clientState exported.ClientState
	err := unpacker.UnpackAny(msg.ClientState, &clientState)
	if err != nil {
		return err
	}

	var consensusState exported.ConsensusState
	return unpacker.UnpackAny(msg.ConsensusState, &consensusState)
}

func NewMsgProxyConnectionOpenTry(
	connectionID string,
	upstreamClientID string,
	upstreamPrefix commitmenttypes.MerklePrefix,
	connection connectiontypes.ConnectionEnd,
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
	proofInit []byte,
	proofClient []byte,
	proofConsensus []byte,
	proofHeight clienttypes.Height,
	consensusHeight clienttypes.Height,
	signer string,
) (*MsgProxyConnectionOpenTry, error) {
	anyClient, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}
	anyConsensus, err := clienttypes.PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}
	return &MsgProxyConnectionOpenTry{
		ConnectionId:     connectionID,
		UpstreamClientId: upstreamClientID,
		UpstreamPrefix:   upstreamPrefix,
		Connection:       connection,
		ClientState:      anyClient,
		ConsensusState:   anyConsensus,
		ProofInit:        proofInit,
		ProofClient:      proofClient,
		ProofConsensus:   proofConsensus,
		ProofHeight:      proofHeight,
		ConsensusHeight:  consensusHeight,
		Signer:           signer,
	}, nil
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyConnectionOpenTry) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgProxyConnectionOpenTry) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgProxyConnectionOpenTry) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var clientState exported.ClientState
	err := unpacker.UnpackAny(msg.ClientState, &clientState)
	if err != nil {
		return err
	}

	var consensusState exported.ConsensusState
	return unpacker.UnpackAny(msg.ConsensusState, &consensusState)
}

func NewMsgProxyConnectionOpenAck(
	connectionID string,
	upstreamClientID string,
	upstreamPrefix commitmenttypes.MerklePrefix,
	connection connectiontypes.ConnectionEnd,
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
	version *connectiontypes.Version,
	proofTry []byte,
	proofClient []byte,
	proofConsensus []byte,
	proofHeight clienttypes.Height,
	consensusHeight clienttypes.Height,
	signer string,
) (*MsgProxyConnectionOpenAck, error) {
	anyClient, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}
	anyConsensus, err := clienttypes.PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}
	return &MsgProxyConnectionOpenAck{
		ConnectionId:     connectionID,
		UpstreamClientId: upstreamClientID,
		UpstreamPrefix:   upstreamPrefix,
		Connection:       connection,
		ClientState:      anyClient,
		ConsensusState:   anyConsensus,
		Version:          version,
		ProofTry:         proofTry,
		ProofClient:      proofClient,
		ProofConsensus:   proofConsensus,
		ProofHeight:      proofHeight,
		ConsensusHeight:  consensusHeight,
		Signer:           signer,
	}, nil
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyConnectionOpenAck) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgProxyConnectionOpenAck) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgProxyConnectionOpenAck) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var clientState exported.ClientState
	err := unpacker.UnpackAny(msg.ClientState, &clientState)
	if err != nil {
		return err
	}

	var consensusState exported.ConsensusState
	return unpacker.UnpackAny(msg.ConsensusState, &consensusState)
}

func NewMsgProxyConnectionOpenConfirm(
	connectionID string,
	upstreamClientID string,
	upstreamPrefix commitmenttypes.MerklePrefix,
	connection connectiontypes.ConnectionEnd,
	proofAck []byte,
	proofHeight clienttypes.Height,
	signer string,
) (*MsgProxyConnectionOpenConfirm, error) {
	return &MsgProxyConnectionOpenConfirm{
		ConnectionId:     connectionID,
		UpstreamClientId: upstreamClientID,
		UpstreamPrefix:   upstreamPrefix,
		Connection:       connection,
		ProofAck:         proofAck,
		ProofHeight:      proofHeight,
		Signer:           signer,
	}, nil
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyConnectionOpenConfirm) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgProxyConnectionOpenConfirm) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgProxyConnectionOpenConfirm) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return nil
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyChannelOpenTry) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgProxyChannelOpenTry) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyChannelOpenAck) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgProxyChannelOpenAck) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyChannelOpenConfirm) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgProxyChannelOpenConfirm) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyRecvPacket) ValidateBasic() error {
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgProxyRecvPacket) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgProxyAcknowledgePacket) ValidateBasic() error {
	return nil
}

func (msg MsgProxyAcknowledgePacket) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}
