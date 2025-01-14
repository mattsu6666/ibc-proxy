package ibctesting

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/modules/core/04-channel/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/modules/core/24-host"
	"github.com/cosmos/ibc-go/modules/core/exported"
	proxytypes "github.com/datachainlab/ibc-proxy/modules/proxy/types"
	"github.com/datachainlab/ibc-proxy/testing/simapp"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (coord *Coordinator) InitProxy(
	source, counterparty *TestChain,
	clientType string, useMultiV bool,
) (clientID string, err error) {
	if useMultiV {
		clientID, err = coord.CreateMultiVClient(source, counterparty, clientType)
	} else {
		clientID, err = coord.CreateClient(source, counterparty, clientType)
	}
	if err != nil {
		return "", err
	}
	err = source.App.(*simapp.SimApp).IBCProxyKeeper.EnableProxy(source.GetContext(), clientID)
	if err != nil {
		return clientID, err
	}
	coord.CommitBlock(source)
	return clientID, err
}

// steps:
// - setup connection
// - setup channel (port=proxy)
// - relay the proxy request packet
func (coord *Coordinator) SetupProxy(
	downstream, proxy *TestChain,
	upstreamClientID string,
) (string, error) {
	// A: downstream, B: proxy

	clientA, clientB := coord.SetupClients(downstream, proxy, exported.Tendermint)
	connA, connB := coord.CreateConnection(downstream, proxy, clientA, clientB, proxytypes.Version)
	chanA, _ := coord.CreateChannel(downstream, proxy, connA, connB, proxytypes.PortID, proxytypes.PortID, types.UNORDERED)

	// begin proxy handshake

	packet, proxyClient, err := downstream.App.(*simapp.SimApp).IBCProxyKeeper.SendProxyRequest(
		downstream.GetContext(),
		chanA.PortID, chanA.ID,
		upstreamClientID,
		clienttypes.NewHeight(0, 110), 0,
	)
	if err != nil {
		return "", err
	}

	ack := proxy.App.(*simapp.SimApp).IBCProxyKeeper.OnRecvPacket(proxy.GetContext(), *packet)
	err = downstream.App.(*simapp.SimApp).IBCProxyKeeper.OnAcknowledgementPacket(
		downstream.GetContext(), *packet, ack.(channeltypes.Acknowledgement),
	)
	if err != nil {
		return proxyClient, err
	}

	coord.CommitBlock(downstream, proxy)

	return proxyClient, nil
}

func (chain *TestChain) UpdateProxyClient(proxy *TestChain, proxyClientID string) error {
	header, err := chain.ConstructUpdateTMClientHeader(proxy, proxyClientID)
	if err != nil {
		return err
	}
	msg, err := clienttypes.NewMsgUpdateClient(proxyClientID, header, chain.SenderAccount.GetAddress().String())
	if err != nil {
		return err
	}
	return chain.sendMsgs(msg)
}

func (coord *Coordinator) CreateConnectionWithProxy(
	chainA, chainB *TestChain,
	clientA, clientB string,
	nextChannelVersion string,
	proxies ProxyPair,
) (*TestConnection, *TestConnection) {

	connA, connB, err := coord.ConnOpenInitWithProxy(chainA, chainB, clientA, clientB, nextChannelVersion, proxies)
	require.NoError(coord.t, err)

	err = coord.ConnOpenTryWithProxy(chainB, chainA, connB, connA, proxies.Swap())
	require.NoError(coord.t, err)

	err = coord.ConnOpenAckWithProxy(chainA, chainB, connA, connB, proxies)
	require.NoError(coord.t, err)

	err = coord.ConnOpenConfirmWithProxy(chainB, chainA, connB, connA, proxies.Swap())
	require.NoError(coord.t, err)

	return connA, connB
}

func (coord *Coordinator) CreateChannelWithProxy(
	chainA, chainB *TestChain,
	connA, connB *TestConnection,
	sourcePortID, counterpartyPortID string,
	order channeltypes.Order,
	proxies ProxyPair,
) (*TestChannel, *TestChannel) {
	channelA, channelB, err := coord.ChanOpenInitWithProxy(chainA, chainB, connA, connB, sourcePortID, counterpartyPortID, order, proxies)
	require.NoError(coord.t, err)

	err = coord.ChanOpenTryWithProxy(chainB, chainA, channelB, channelA, connB, connA, order, proxies.Swap())
	require.NoError(coord.t, err)

	err = coord.ChanOpenAckWithProxy(chainA, chainB, channelA, channelB, connA, connB, order, proxies)
	require.NoError(coord.t, err)

	err = coord.ChanOpenConfirmWithProxy(chainB, chainA, channelB, channelA, connB, connA, order, proxies.Swap())
	require.NoError(coord.t, err)

	return &channelA, &channelB
}

type ProxyInfo struct {
	Chain            *TestChain
	ClientID         string
	UpstreamClientID string
	UpstreamPrefix   exported.Prefix
}

type ProxyPair [2]*ProxyInfo

func (pair ProxyPair) Swap() ProxyPair {
	pair[0], pair[1] = pair[1], pair[0]
	return pair
}

func (coord *Coordinator) ConnOpenInitWithProxy(
	source, counterparty *TestChain,
	clientID, counterpartyClientID, nextChannelVersion string, proxies ProxyPair,
) (*TestConnection, *TestConnection, error) {
	if proxies[1] == nil {
		return coord.ConnOpenInit(source, counterparty, clientID, counterpartyClientID, nextChannelVersion)
	}

	sourceConnection := source.AddTestConnection(clientID, counterpartyClientID, nextChannelVersion)
	counterpartyConnection := counterparty.AddTestConnection(counterpartyClientID, clientID, nextChannelVersion)

	// initialize connection on source
	if err := source.ConnectionOpenInit(counterparty, sourceConnection, counterpartyConnection); err != nil {
		return sourceConnection, counterpartyConnection, err
	}
	coord.IncrementTime()

	// update source client on counterparty connection
	if err := coord.UpdateClient(
		proxies[1].Chain, source,
		proxies[1].UpstreamClientID, exported.Tendermint,
	); err != nil {
		return sourceConnection, counterpartyConnection, err
	}

	return sourceConnection, counterpartyConnection, nil
}

func (coord *Coordinator) ConnOpenTryWithProxy(
	source, counterparty *TestChain,
	sourceConnection, counterpartyConnection *TestConnection,
	proxies ProxyPair,
) error {
	if proxies[0] == nil {
		if err := source.ConnectionOpenTryWithProxy(counterparty, sourceConnection, counterpartyConnection, *proxies[1]); err != nil {
			return err
		}
		coord.IncrementTime()
	} else {
		// source: downstream, counterparty: upstream
		proxy := proxies[0].Chain

		proxyConnection := connectiontypes.NewConnectionEnd(
			connectiontypes.INIT,
			counterpartyConnection.ClientID,
			connectiontypes.Counterparty{
				ClientId:     proxies[0].ClientID,
				ConnectionId: "",
				Prefix:       counterparty.GetPrefix(),
			},
			[]*connectiontypes.Version{ConnectionVersion},
			DefaultDelayPeriod,
		)

		var (
			counterpartyClient exported.ClientState
			proofClient        []byte
			consensusState     exported.ConsensusState
			proofConsensus     []byte
			consensusHeight    clienttypes.Height
		)

		if proxies[1] == nil {
			var found bool
			counterpartyClient, proofClient = counterparty.QueryClientStateProof(counterpartyConnection.ClientID)
			proofConsensus, consensusHeight = counterparty.QueryConsensusStateProof(counterpartyConnection.ClientID)
			consensusState, found = counterparty.GetConsensusState(counterpartyConnection.ClientID, consensusHeight)
			if !found {
				return fmt.Errorf("consensusState '%v-%v' not found", counterpartyConnection.ClientID, consensusHeight)
			}
		} else {
			head := counterparty.QueryMultiVBranchProof(counterpartyConnection.ClientID)
			counterpartyProxy := *proxies[1]
			counterpartyClient, proofClient = counterpartyProxy.Chain.QueryMultiVLeafClientProof(head, counterpartyProxy.UpstreamClientID)
			consensusState, proofConsensus, consensusHeight = counterpartyProxy.Chain.QueryMultiVLeafConsensusProof(head, counterpartyProxy.UpstreamClientID)
		}

		proofInit, proofHeight := counterparty.QueryProof(host.ConnectionKey(counterpartyConnection.ID))
		msg, err := proxytypes.NewMsgProxyConnectionOpenTry(
			counterpartyConnection.ID,
			proxies[0].UpstreamClientID,
			proxies[0].UpstreamPrefix.(commitmenttypes.MerklePrefix),
			proxyConnection,
			counterpartyClient,
			consensusState,
			proofInit, proofClient, proofConsensus, proofHeight, consensusHeight, proxy.SenderAccount.GetAddress().String(),
		)
		if err != nil {
			return err
		}
		if _, err := proxy.SendMsgs(msg); err != nil {
			return err
		}
		coord.CommitBlock(proxy)

		if err := source.UpdateProxyClient(proxy, proxies[0].ClientID); err != nil {
			return err
		}
		coord.CommitBlock(source)

		{
			client, proofClient := proxy.QueryProxyClientStateProof(counterpartyConnection.ClientID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)
			proofInit, proofHeight := proxy.QueryProxyConnectionStateProof(counterpartyConnection.ID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)
			proofConsensus, consensusHeight := proxy.QueryProxyConsensusStateProof(counterpartyConnection.ClientID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)

			msg := connectiontypes.NewMsgConnectionOpenTry(
				"", sourceConnection.ClientID, counterpartyConnection.ID, counterpartyConnection.ClientID, client, // testing doesn't use flexible selection
				counterparty.GetPrefix(), []*connectiontypes.Version{ConnectionVersion}, DefaultDelayPeriod,
				proofInit, proofClient, proofConsensus,
				proofHeight, consensusHeight,
				source.SenderAccount.GetAddress().String(),
			)
			if _, err := source.SendMsgs(msg); err != nil {
				return err
			}
			coord.CommitBlock(source)
		}
	}

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

func (coord *Coordinator) ConnOpenAckWithProxy(
	source, counterparty *TestChain,
	sourceConnection, counterpartyConnection *TestConnection,
	proxies ProxyPair,
) error {
	if proxies[0] == nil {
		if err := source.ConnectionOpenAckWithProxy(counterparty, sourceConnection, counterpartyConnection, *proxies[1]); err != nil {
			return err
		}
		coord.IncrementTime()
	} else {
		// source: downstream, counterparty: upstream
		proxy := proxies[0].Chain

		proxyConnection := connectiontypes.NewConnectionEnd(
			connectiontypes.TRYOPEN,
			counterpartyConnection.ClientID,
			connectiontypes.Counterparty{
				ClientId:     sourceConnection.ClientID,
				ConnectionId: sourceConnection.ID,
				Prefix:       counterparty.GetPrefix(),
			},
			[]*connectiontypes.Version{ConnectionVersion},
			DefaultDelayPeriod,
		)

		{
			var (
				counterpartyClient exported.ClientState
				proofClient        []byte
				consensusState     exported.ConsensusState
				proofConsensus     []byte
				consensusHeight    clienttypes.Height
			)

			if proxies[1] == nil {
				var found bool
				counterpartyClient, proofClient = counterparty.QueryClientStateProof(counterpartyConnection.ClientID)
				proofConsensus, consensusHeight = counterparty.QueryConsensusStateProof(counterpartyConnection.ClientID)
				consensusState, found = counterparty.GetConsensusState(counterpartyConnection.ClientID, consensusHeight)
				if !found {
					return fmt.Errorf("consensusState '%v-%v' not found", counterpartyConnection.ClientID, consensusHeight)
				}
			} else {
				head := counterparty.QueryMultiVBranchProof(counterpartyConnection.ClientID)
				counterpartyProxy := *proxies[1]
				counterpartyClient, proofClient = counterpartyProxy.Chain.QueryMultiVLeafClientProof(head, counterpartyProxy.UpstreamClientID)
				consensusState, proofConsensus, consensusHeight = counterpartyProxy.Chain.QueryMultiVLeafConsensusProof(head, counterpartyProxy.UpstreamClientID)
			}

			proofTry, proofHeight := counterparty.QueryProof(host.ConnectionKey(counterpartyConnection.ID))

			msg, err := proxytypes.NewMsgProxyConnectionOpenAck(
				counterpartyConnection.ID,
				proxies[0].UpstreamClientID,
				proxies[0].UpstreamPrefix.(commitmenttypes.MerklePrefix),
				proxyConnection,
				counterpartyClient,
				consensusState,
				ConnectionVersion,
				proofTry,
				proofClient,
				proofConsensus,
				proofHeight,
				consensusHeight,
				proxy.SenderAccount.GetAddress().String(),
			)
			if err != nil {
				return err
			}
			if _, err := proxy.SendMsgs(msg); err != nil {
				return err
			}
			coord.CommitBlock(proxy)
		}

		if err := source.UpdateProxyClient(proxy, proxies[0].ClientID); err != nil {
			return err
		}
		coord.CommitBlock(source)

		{ // callerA calls connOpenAck with proxied proof
			client, proofClient := proxy.QueryProxyClientStateProof(counterpartyConnection.ClientID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)
			proofTry, proofHeight := proxy.QueryProxyConnectionStateProof(counterpartyConnection.ID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)
			proofConsensus, consensusHeight := proxy.QueryProxyConsensusStateProof(counterpartyConnection.ClientID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)

			msg := connectiontypes.NewMsgConnectionOpenAck(
				sourceConnection.ID, counterpartyConnection.ID, client, // testing doesn't use flexible selection
				proofTry, proofClient, proofConsensus,
				proofHeight, consensusHeight,
				proxyConnection.Versions[0],
				source.SenderAccount.GetAddress().String(),
			)
			if _, err := source.SendMsgs(msg); err != nil {
				return err
			}
			coord.CommitBlock(source)
		}
	}

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

func (coord *Coordinator) ConnOpenConfirmWithProxy(
	source, counterparty *TestChain,
	sourceConnection, counterpartyConnection *TestConnection,
	proxies ProxyPair,
) error {
	if proxies[0] == nil {
		if err := source.ConnectionOpenConfirm(counterparty, sourceConnection, counterpartyConnection); err != nil {
			return err
		}
		coord.IncrementTime()
	} else {
		// source: downstream, counterparty: upstream
		proxy := proxies[0].Chain

		proxyConnection := connectiontypes.NewConnectionEnd(
			connectiontypes.OPEN,
			counterpartyConnection.ClientID,
			connectiontypes.Counterparty{
				ClientId:     sourceConnection.ClientID,
				ConnectionId: sourceConnection.ID,
				Prefix:       counterparty.GetPrefix(),
			},
			[]*connectiontypes.Version{ConnectionVersion},
			DefaultDelayPeriod,
		)

		connectionKey := host.ConnectionKey(counterpartyConnection.ID)
		proofAck, proofHeight := counterparty.QueryProof(connectionKey)

		msg, err := proxytypes.NewMsgProxyConnectionOpenConfirm(
			counterpartyConnection.ID,
			proxies[0].UpstreamClientID,
			proxies[0].UpstreamPrefix.(commitmenttypes.MerklePrefix),
			proxyConnection,
			proofAck,
			proofHeight,
			proxy.SenderAccount.GetAddress().String(),
		)
		if err != nil {
			return err
		}
		if _, err := proxy.SendMsgs(msg); err != nil {
			return err
		}
		coord.CommitBlock(proxy)

		if err := source.UpdateProxyClient(proxy, proxies[0].ClientID); err != nil {
			return err
		}
		coord.CommitBlock(source)

		{
			proofAck, proofHeight := proxy.QueryProxyConnectionStateProof(counterpartyConnection.ID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)

			msg := connectiontypes.NewMsgConnectionOpenConfirm(
				sourceConnection.ID, proofAck, proofHeight,
				source.SenderAccount.GetAddress().String(),
			)
			if _, err := source.SendMsgs(msg); err != nil {
				return err
			}
			coord.CommitBlock(source)
		}
	}

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

func (coord *Coordinator) ChanOpenInitWithProxy(
	source, counterparty *TestChain,
	connection, counterpartyConnection *TestConnection,
	sourcePortID, counterpartyPortID string,
	order channeltypes.Order,
	proxies ProxyPair,
) (TestChannel, TestChannel, error) {
	if proxies[1] == nil {
		return coord.ChanOpenInit(source, counterparty, connection, counterpartyConnection, sourcePortID, counterpartyPortID, order)
	}

	sourceChannel := source.AddTestChannel(connection, sourcePortID)
	counterpartyChannel := counterparty.AddTestChannel(counterpartyConnection, counterpartyPortID)

	// NOTE: only creation of a capability for a transfer or mock port is supported
	// Other applications must bind to the port in InitGenesis or modify this code.
	source.CreatePortCapability(sourceChannel.PortID)
	coord.IncrementTime()

	// initialize channel on source
	if err := source.ChanOpenInit(sourceChannel, counterpartyChannel, order, connection.ID); err != nil {
		return sourceChannel, counterpartyChannel, err
	}
	coord.IncrementTime()

	// update source client on counterparty connection
	if err := coord.UpdateClient(
		proxies[1].Chain, source,
		proxies[1].UpstreamClientID, exported.Tendermint,
	); err != nil {
		return sourceChannel, counterpartyChannel, err
	}

	return sourceChannel, counterpartyChannel, nil
}

func (coord *Coordinator) ChanOpenTryWithProxy(
	source, counterparty *TestChain,
	sourceChannel, counterpartyChannel TestChannel,
	sourceConnection, counterpartyConnection *TestConnection,
	order channeltypes.Order,
	proxies ProxyPair,
) error {
	if proxies[0] == nil {
		if err := source.ChanOpenTry(counterparty, sourceChannel, counterpartyChannel, order, sourceConnection.ID); err != nil {
			return err
		}
		coord.IncrementTime()
	} else {
		// source: downstream, counterparty: upstream
		proxy := proxies[0].Chain

		proofInit, proofHeight := counterparty.QueryProof(host.ChannelKey(counterpartyChannel.PortID, counterpartyChannel.ID))

		msg := &proxytypes.MsgProxyChannelOpenTry{
			UpstreamClientId: proxies[0].UpstreamClientID,
			UpstreamPrefix:   proxies[0].UpstreamPrefix.(commitmenttypes.MerklePrefix),
			Order:            order,
			ConnectionHops:   []string{counterpartyConnection.ID},
			PortId:           counterpartyChannel.PortID,
			ChannelId:        counterpartyChannel.ID,
			DownstreamPortId: sourceChannel.PortID,
			Version:          counterpartyChannel.Version,
			ProofInit:        proofInit,
			ProofHeight:      proofHeight,
			Signer:           proxy.SenderAccount.GetAddress().String(),
		}
		if _, err := proxy.SendMsgs(msg); err != nil {
			return err
		}
		coord.CommitBlock(proxy)

		if err := source.UpdateProxyClient(proxy, proxies[0].ClientID); err != nil {
			return err
		}
		coord.CommitBlock(source)

		{
			proof, proofHeight := proxy.QueryProxyChannelStateProof(counterpartyChannel.PortID, counterpartyChannel.ID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)
			msg := channeltypes.NewMsgChannelOpenTry(
				sourceChannel.PortID,
				"",
				sourceChannel.Version, order,
				[]string{sourceConnection.ID},
				counterpartyChannel.PortID, counterpartyChannel.ID, counterpartyChannel.Version,
				proof, proofHeight,
				source.SenderAccount.GetAddress().String(),
			)
			if _, err := source.SendMsgs(msg); err != nil {
				return err
			}
			coord.CommitBlock(source)
		}
	}

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

func (coord *Coordinator) ChanOpenAckWithProxy(
	source, counterparty *TestChain,
	sourceChannel, counterpartyChannel TestChannel,
	sourceConnection, counterpartyConnection *TestConnection,
	order channeltypes.Order,
	proxies ProxyPair,
) error {
	if proxies[0] == nil {
		if err := source.ChanOpenAck(counterparty, sourceChannel, counterpartyChannel); err != nil {
			return err
		}
		coord.IncrementTime()
	} else {
		// source: downstream, counterparty: upstream
		proxy := proxies[0].Chain

		proofTry, proofHeight := counterparty.QueryProof(host.ChannelKey(counterpartyChannel.PortID, counterpartyChannel.ID))

		msg := &proxytypes.MsgProxyChannelOpenAck{
			UpstreamClientId:    proxies[0].UpstreamClientID,
			UpstreamPrefix:      proxies[0].UpstreamPrefix.(commitmenttypes.MerklePrefix),
			Order:               order,
			ConnectionHops:      []string{counterpartyConnection.ID},
			PortId:              counterpartyChannel.PortID,
			ChannelId:           counterpartyChannel.ID,
			DownstreamPortId:    sourceChannel.PortID,
			DownstreamChannelId: sourceChannel.ID,
			Version:             counterpartyChannel.Version,
			ProofTry:            proofTry,
			ProofHeight:         proofHeight,
			Signer:              proxy.SenderAccount.GetAddress().String(),
		}
		if _, err := proxy.SendMsgs(msg); err != nil {
			return err
		}
		coord.CommitBlock(proxy)

		if err := source.UpdateProxyClient(proxy, proxies[0].ClientID); err != nil {
			return err
		}
		coord.CommitBlock(source)

		{
			proof, proofHeight := proxy.QueryProxyChannelStateProof(counterpartyChannel.PortID, counterpartyChannel.ID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)
			msg := channeltypes.NewMsgChannelOpenAck(
				sourceChannel.PortID, sourceChannel.ID,
				counterpartyChannel.ID, counterpartyChannel.Version,
				proof, proofHeight,
				source.SenderAccount.GetAddress().String(),
			)
			if _, err := source.SendMsgs(msg); err != nil {
				return err
			}
			coord.CommitBlock(source)
		}
	}

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

func (coord *Coordinator) ChanOpenConfirmWithProxy(
	source, counterparty *TestChain,
	sourceChannel, counterpartyChannel TestChannel,
	sourceConnection, counterpartyConnection *TestConnection,
	order channeltypes.Order,
	proxies ProxyPair,
) error {

	if proxies[0] == nil {
		if err := source.ChanOpenConfirm(counterparty, sourceChannel, counterpartyChannel); err != nil {
			return err
		}
		coord.IncrementTime()
	} else {
		// source: downstream, counterparty: upstream
		proxy := proxies[0].Chain

		proofAck, proofHeight := counterparty.QueryProof(host.ChannelKey(counterpartyChannel.PortID, counterpartyChannel.ID))

		msg := &proxytypes.MsgProxyChannelOpenConfirm{
			UpstreamClientId:    proxies[0].UpstreamClientID,
			UpstreamPrefix:      proxies[0].UpstreamPrefix.(commitmenttypes.MerklePrefix),
			PortId:              counterpartyChannel.PortID,
			ChannelId:           counterpartyChannel.ID,
			DownstreamChannelId: sourceChannel.ID,
			ProofAck:            proofAck,
			ProofHeight:         proofHeight,
			Signer:              proxy.SenderAccount.GetAddress().String(),
		}
		if _, err := proxy.SendMsgs(msg); err != nil {
			return err
		}
		coord.CommitBlock(proxy)

		if err := source.UpdateProxyClient(proxy, proxies[0].ClientID); err != nil {
			return err
		}
		coord.CommitBlock(source)

		{
			proof, proofHeight := proxy.QueryProxyChannelStateProof(counterpartyChannel.PortID, counterpartyChannel.ID, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)

			msg := channeltypes.NewMsgChannelOpenConfirm(
				sourceChannel.PortID, sourceChannel.ID,
				proof, proofHeight,
				source.SenderAccount.GetAddress().String(),
			)
			if _, err := source.SendMsgs(msg); err != nil {
				return err
			}
			coord.CommitBlock(source)
		}
	}

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

func (coord *Coordinator) SendPacketWithProxy(
	source, counterparty *TestChain, // source: packet sender, counterparty: packet receiver
	sourceConnection, counterpartyConnection *TestConnection,
	proxies ProxyPair,
	msgs ...sdk.Msg,
) error {
	if _, err := source.SendMsgs(msgs...); err != nil {
		return err
	}
	coord.CommitBlock(source)

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

func (coord *Coordinator) RecvPacketWithProxy(
	source, counterparty *TestChain, // source: packet receiver, counterparty: packet sender
	sourceConnection, counterpartyConnection *TestConnection,
	packet channeltypes.Packet, proxies ProxyPair,
) error {
	if proxies[0] == nil {
		if err := counterparty.recvPacket(coord, source, counterpartyConnection.ClientID, packet); err != nil {
			return err
		}
	} else {
		// source: downstream, counterparty: upstream
		proxy := proxies[0].Chain
		proof, proofHeight := counterparty.QueryProof(host.PacketCommitmentKey(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence()))
		msg := &proxytypes.MsgProxyRecvPacket{
			UpstreamClientId: proxies[0].UpstreamClientID,
			UpstreamPrefix:   proxies[0].UpstreamPrefix.(commitmenttypes.MerklePrefix),
			Packet:           packet,
			Proof:            proof,
			ProofHeight:      proofHeight,
			Signer:           proxy.SenderAccount.GetAddress().String(),
		}
		if _, err := proxy.SendMsgs(msg); err != nil {
			return err
		}
		coord.CommitBlock(proxy)

		if err := source.UpdateProxyClient(proxy, proxies[0].ClientID); err != nil {
			return err
		}
		coord.CommitBlock(source)

		// relay the packet to the source chain
		{
			proof, proofHeight := proxy.QueryProxyPacketCommitmentProof(packet.SourcePort, packet.SourceChannel, packet.Sequence, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)
			recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, source.SenderAccount.GetAddress().String())
			if _, err := source.SendMsgs(recvMsg); err != nil {
				return err
			}
			coord.CommitBlock(source)
		}
	}

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

func (coord *Coordinator) AcknowledgePacketWithProxy(
	source, counterparty *TestChain, // source: ack receiver, counterparty: ack sender
	sourceChannel, counterpartyChannel *TestChannel,
	sourceConnection, counterpartyConnection *TestConnection,
	packet channeltypes.Packet, ack []byte, proxies ProxyPair,
) error {
	if proxies[0] == nil {
		if err := source.acknowledgePacket(coord, counterparty, sourceConnection.ClientID, packet, ack); err != nil {
			return err
		}
	} else {
		// source: downstream, counterparty: upstream
		proxy := proxies[0].Chain
		proof, proofHeight := counterparty.QueryProof(host.PacketAcknowledgementKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence()))
		msg := &proxytypes.MsgProxyAcknowledgePacket{
			UpstreamClientId: proxies[0].UpstreamClientID,
			UpstreamPrefix:   proxies[0].UpstreamPrefix.(commitmenttypes.MerklePrefix),
			Packet:           packet,
			Acknowledgement:  ack,
			Proof:            proof,
			ProofHeight:      proofHeight,
			Signer:           proxy.SenderAccount.GetAddress().String(),
		}
		if _, err := proxy.SendMsgs(msg); err != nil {
			return err
		}
		coord.CommitBlock(proxy)

		if err := source.UpdateProxyClient(proxy, proxies[0].ClientID); err != nil {
			return err
		}
		coord.CommitBlock(source)

		{
			proof, proofHeight := proxy.QueryProxyAcknowledgementProof(packet.DestinationPort, packet.DestinationChannel, packet.Sequence, proxies[0].UpstreamPrefix, proxies[0].UpstreamClientID)
			ackMsg := channeltypes.NewMsgAcknowledgement(packet, ack, proof, proofHeight, source.SenderAccount.GetAddress().String())
			if _, err := source.SendMsgs(ackMsg); err != nil {
				return err
			}
			coord.CommitBlock(source)
		}
	}

	if proxies[1] == nil {
		return coord.UpdateClient(
			counterparty, source, counterpartyConnection.ClientID, exported.Tendermint,
		)
	} else {
		return coord.UpdateClient(
			proxies[1].Chain, source, proxies[1].UpstreamClientID, exported.Tendermint,
		)
	}
}

// source: packet sender, counterparty: packet receiver
func (chain *TestChain) recvPacket(coord *Coordinator, counterparty *TestChain, sourceClient string, packet channeltypes.Packet) error {
	// get proof of packet commitment on source
	packetKey := host.PacketCommitmentKey(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	proof, proofHeight := chain.QueryProof(packetKey)

	// Increment time and commit block so that 5 second delay period passes between send and receive
	coord.IncrementTime()
	coord.CommitBlock(chain, counterparty)

	recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, counterparty.SenderAccount.GetAddress().String())
	if err := counterparty.sendMsgs(recvMsg); err != nil {
		return err
	}
	coord.IncrementTime()
	return nil
}

// source: packet sender / ack receiver, counterparty: packet receiver / ack sender
func (chain *TestChain) acknowledgePacket(coord *Coordinator, counterparty *TestChain, sourceClient string, packet channeltypes.Packet, ack []byte) error {
	// get proof of acknowledgement on counterparty
	packetKey := host.PacketAcknowledgementKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
	proof, proofHeight := counterparty.QueryProof(packetKey)

	coord.IncrementTime()
	coord.CommitBlock(chain, counterparty)

	ackMsg := channeltypes.NewMsgAcknowledgement(packet, ack, proof, proofHeight, chain.SenderAccount.GetAddress().String())
	if err := chain.sendMsgs(ackMsg); err != nil {
		return err
	}
	coord.IncrementTime()
	return nil
}

func (chain *TestChain) QueryProxyClientStateProof(clientID string, upstreamPrefix exported.Prefix, upstreamClientID string) (exported.ClientState, []byte) {
	// retrieve client state to provide proof for
	clientState, found := chain.App.(*simapp.SimApp).IBCProxyKeeper.GetClientStateCommitment(
		chain.GetContext(),
		upstreamPrefix,
		clientID,
		upstreamClientID,
	)
	require.True(chain.t, found)
	proofClient, _ := chain.QueryProxyProof(proxytypes.ProxyClientStateKey(upstreamPrefix, upstreamClientID, clientID))
	return clientState, proofClient
}

func (chain *TestChain) QueryProxyConsensusStateProof(clientID string, upstreamPrefix exported.Prefix, upstreamClientID string) ([]byte, clienttypes.Height) {
	clientState, found := chain.App.(*simapp.SimApp).IBCProxyKeeper.GetClientStateCommitment(
		chain.GetContext(),
		upstreamPrefix,
		clientID,
		upstreamClientID,
	)
	require.True(chain.t, found)
	consensusHeight := clientState.GetLatestHeight().(clienttypes.Height)
	proofConsensus, _ := chain.QueryProxyProof(proxytypes.ProxyConsensusStateKey(upstreamPrefix, upstreamClientID, clientID, consensusHeight))
	return proofConsensus, consensusHeight
}

func (chain *TestChain) QueryProxyConnectionStateProof(connectionID string, upstreamPrefix exported.Prefix, upstreamClientID string) ([]byte, clienttypes.Height) {
	return chain.QueryProxyProof(proxytypes.ProxyConnectionKey(upstreamPrefix, upstreamClientID, connectionID))
}

func (chain *TestChain) QueryProxyChannelStateProof(portID string, channelID string, upstreamPrefix exported.Prefix, upstreamClientID string) ([]byte, clienttypes.Height) {
	return chain.QueryProxyProof(proxytypes.ProxyChannelKey(upstreamPrefix, upstreamClientID, portID, channelID))
}

func (chain *TestChain) QueryProxyPacketCommitmentProof(sourcePort, sourceChannel string, packetSequence uint64, upstreamPrefix exported.Prefix, upstreamClientID string) ([]byte, clienttypes.Height) {
	return chain.QueryProxyProof(proxytypes.ProxyPacketCommitmentKey(upstreamPrefix, upstreamClientID, sourcePort, sourceChannel, packetSequence))
}

func (chain *TestChain) QueryProxyAcknowledgementProof(destPort, destChannel string, packetSequence uint64, upstreamPrefix exported.Prefix, upstreamClientID string) ([]byte, clienttypes.Height) {
	return chain.QueryProxyProof(proxytypes.ProxyAcknowledgementKey(upstreamPrefix, upstreamClientID, destPort, destChannel, packetSequence))
}

// QueryProof performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryProxyProof(key []byte) ([]byte, clienttypes.Height) {
	res := chain.App.Query(abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", proxytypes.StoreKey),
		Height: chain.App.LastBlockHeight() - 1,
		Data:   key,
		Prove:  true,
	})

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	require.NoError(chain.t, err)

	proof, err := chain.App.AppCodec().Marshal(&merkleProof)
	require.NoError(chain.t, err)

	revision := clienttypes.ParseChainID(chain.ChainID)

	// proof height + 1 is returned as the proof created corresponds to the height the proof
	// was created in the IAVL tree. Tendermint and subsequently the clients that rely on it
	// have heights 1 above the IAVL tree. Thus we return proof height + 1
	return proof, clienttypes.NewHeight(revision, uint64(res.Height)+1)
}
