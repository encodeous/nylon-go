//go:generate protoc -I . --go_out=. ./protocol/nylon.proto
package main

import (
	"github.com/encodeous/nylon/cmd"
)

//func mock() (*state.CentralCfg, *state.NodeCfg, error) {
//	_, nodeKey, err := ed25519.GenerateKey(nil)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	certTemplate := x509.Certificate{
//		PublicKey: nodeKey.Public(),
//		Subject: pkix.Name{
//			CommonName: "dummyNode",
//		},
//		IsCA:         false,
//		SubjectKeyId: nil,
//		NotBefore:    time.Now(),
//		NotAfter:     time.Now().AddDate(10, 0, 0),
//		SerialNumber: big.NewInt(time.Now().Unix()),
//	}
//
//	ss, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, certTemplate.PublicKey, nodeKey)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	dpKey, err := ecdh.X25519().GenerateKey(rand.Reader)
//	mockNode := state.NodeCfg{
//		Id:    "currentNode",
//		Key:   state.EdPrivateKey(nodeKey),
//		WgKey: (*state.EcPrivateKey)(dpKey),
//		Cert:  state.Cert(ss),
//	}
//
//	mockPubNode := mockNode.GeneratePubCfg()
//
//	mockCentralCfg := state.CentralCfg{
//		RootPubKey: ss,
//		Nodes: []state.PubNodeCfg{
//			mockPubNode,
//		},
//		Version: 0,
//	}
//
//	return &mockCentralCfg, &mockNode, nil
//}

func main() {
	cmd.Execute()
}
