package state

import (
	"context"
	"crypto/ecdh"
	"crypto/ed25519"
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"log/slog"
	"net/netip"
	"time"
)

type NyModule interface {
	Init(s *State) error
	Cleanup(s *State) error
}

type State struct {
	*Env
	TrustedNodes map[Node]ed25519.PublicKey
	Modules      map[string]NyModule
}

type Env struct {
	DispatchChannel chan<- func(s *State) error
	LinkChannel     chan<- CtlLink
	CentralCfg
	NodeCfg
	Context context.Context
	Cancel  context.CancelCauseFunc
	PingBuf *ttlcache.Cache[uint64, LinkPing]
	Log     *slog.Logger
}

type DpEndpoint struct {
	Name      string
	DpAddr    *netip.AddrPort
	ProbeAddr *netip.AddrPort
}

// NodeCfg represents local node-level configuration
type NodeCfg struct {
	// Node Private Key
	Key EdPrivateKey
	// Data plane (WireGuard) Private key
	WgKey *EcPrivateKey
	// x509 certificate signed by the root CA
	Cert Cert
	Id   Node
	// Address and port that the control plane listens on
	CtlBind string
	// Address that the data plane can be accessed by
	DpBind    netip.AddrPort
	ProbeBind netip.AddrPort
}

func (k *EcPublicKey) Bytes() []byte {
	return (*ecdh.PublicKey)(k).Bytes()
}

func (k *EcPrivateKey) Bytes() []byte {
	return (*ecdh.PrivateKey)(k).Bytes()
}

func (k *EcPrivateKey) Pubkey() *EcPublicKey {
	return (*EcPublicKey)(((*ecdh.PrivateKey)(k).Public()).(*ecdh.PublicKey))
}

func (k EdPrivateKey) Pubkey() EdPublicKey {
	return EdPublicKey(((ed25519.PrivateKey)(k).Public()).(ed25519.PublicKey))
}

func (n NodeCfg) GeneratePubCfg() PubNodeCfg {
	cfg := PubNodeCfg{
		Id:      n.Id,
		CtlAddr: []string{n.CtlBind},
		DpAddr: []DpEndpoint{
			{fmt.Sprintf("%s-local", n.Id), &n.DpBind, &n.ProbeBind},
		},
	}
	if n.WgKey != nil {
		cfg.DpPubKey = n.WgKey.Pubkey()
	}
	if n.Key != nil {
		cfg.PubKey = n.Key.Pubkey()
	}
	return cfg
}

// PubNodeCfg represents a central representation of a node
type PubNodeCfg struct {
	Id       Node
	PubKey   EdPublicKey
	DpPubKey *EcPublicKey
	CtlAddr  []string
	DpAddr   []DpEndpoint
}

type CentralCfg struct {
	// the public key of the root CA
	RootCa      Cert
	Nodes       []PubNodeCfg
	Edges       []Pair[Node, Node]
	MockWeights []Triple[Node, Node, *time.Duration]
	Version     uint64
}

func (e Env) GetPeers() []Node {
	nodes := make([]Node, 0)
	for _, edge := range e.Edges {
		var neighNode Node
		if edge.V1 == e.Id {
			neighNode = edge.V2
		}
		if edge.V2 == e.Id {
			neighNode = edge.V1
		}
		if neighNode != e.Id && neighNode != "" {
			nodes = append(nodes, neighNode)
		}
	}
	return nodes
}

// Dispatch Dispatches the function to run on the main thread without waiting for it to complete
func (e Env) Dispatch(fun func(*State) error) {
	e.DispatchChannel <- fun
}

// DispatchWait Dispatches the function to run on the main thread and wait for it to complete
func (e Env) DispatchWait(fun func(*State) (any, error)) (any, error) {
	ret := make(chan Pair[any, error])
	e.DispatchChannel <- func(s *State) error {
		res, err := fun(s)
		ret <- Pair[any, error]{res, err}
		return err
	}
	select {
	case res := <-ret:
		return res.V1, res.V2
	case <-e.Context.Done():
		return nil, e.Context.Err()
	}
}

func (e Env) scheduledTask(fun func(*State) error, delay time.Duration) {
	time.Sleep(delay)
	e.Dispatch(fun)
}

func (e Env) ScheduleTask(fun func(*State) error, delay time.Duration) {
	go e.scheduledTask(fun, delay)
}

func (e Env) repeatedTask(fun func(*State) error, delay time.Duration) {
	for e.Context.Err() == nil {
		e.Dispatch(fun)
		time.Sleep(delay)
	}
}

func (e Env) RepeatTask(fun func(*State) error, delay time.Duration) {
	go e.repeatedTask(fun, delay)
}
