package protocol

import (
	"github.com/anconprotocol/sdk"
	"github.com/anconprotocol/sdk/proofsignature"
)

type RPCApi struct {
	Namespace string
	Version   string
	Service   *RPCService
	Public    bool
}

type RPCService struct {
	Storage *sdk.Storage
	Proof   *proofsignature.IavlProofAPI
}

func NewRPCApi(storage *sdk.Storage, proof *proofsignature.IavlProofAPI) *RPCApi {
	return &RPCApi{
		Namespace: "dag",
		Version:   "1.0",
		Service: &RPCService{
			Storage: storage,
			Proof:   proof,
		},
		Public: true,
	}
}

// func (s *ProtocolService) Add(to string, from string, sig []byte, data string) string {

// 	return hexutil.Encode(res.([]byte))
// }

// func (s *ProtocolService) Mutate(to string, from string, sig []byte, data string) string {

// 	return hexutil.Encode(res.([]byte))
// }

// func (s *ProtocolService) Get(to string, from string, sig []byte, data string) string {

// 	return hexutil.Encode(res.([]byte))
// }
