package ethereum

import (
	"encoding/base64"

	"github.com/buger/jsonparser"

	ics23 "github.com/confio/ics23/go"
)

type EncodePackedExistenceProof struct {
	LeafOp        []int32
	InnerOpSuffix []byte
	InnerOpPrefix []byte
	InnerOpHashOp int32
	Prefix        []byte
	Key           string
	Value         string
}

var template = `{
	"proof": {
		"proofs": [
			{
				"Proof": {
					"exist": {
						"key": "YmFndXFlZXJhcDVwZHlmend6ZDV4NnR2cGJibjRsZno0bWg2cnd2dGhvbmJ3ZzRrbnBkcGt5dDNucmhrcQ==",
						"value":"...",
						"leaf": {
							"hash": 1,
							"prehash_value": 1,
							"length": 1,
							"prefix": "AAIC"
						},
						"path": [
							{
								"hash": 1,
								"prefix": "AgQCIGhqEPIiQrR2tMcmliOUwD/Yq+51sHW7EIDc5BAgCtIpIA=="
							}
						]
					}
				}
			}
		]
	},
	"value": "..."
}`

func (a *OnchainAdapter) MarshalProof(v []byte) *EncodePackedExistenceProof {
	// "key": "YmFndXFlZXJhcDVwZHlmend6ZDV4NnR2cGJibjRsZno0bWg2cnd2dGhvbmJ3ZzRrbnBkcGt5dDNucmhrcQ==",
	// "value":"",
	// "leaf": {
	// 	"hash": 1,
	// 	"prehash_value": 1,
	// 	"length": 1,
	// 	"prefix": "AAIC"
	// },
	// "path": [
	// 	{
	// 		"hash": 1,
	// 		"prefix": "AgQCIGhqEPIiQrR2tMcmliOUwD/Yq+51sHW7EIDc5BAgCtIpIA=="
	// 	}
	// ]

	existPayload, _, _, err := jsonparser.Get(v, "proof", "proofs", "[0]", "Proof", "exist")
	if err != nil {
		return nil
	}
	k, _, _, _ := jsonparser.Get(existPayload, "key")
	val, _, _, _ := jsonparser.Get(existPayload, "value")

	hash, _ := jsonparser.GetInt(existPayload, "leaf", "hash")
	prehashValue, _ := jsonparser.GetInt(existPayload, "leaf", "prehash_value")
	length, _ := jsonparser.GetInt(existPayload, "leaf", "length")
	pre, _ := jsonparser.GetString(existPayload, "leaf", "prefix")

	ophash, _ := jsonparser.GetInt(existPayload, "path", "[0]", "hash")
	opprefix, _ := jsonparser.GetString(existPayload, "path", "[0]", "prefix")
	opsuffix, _ := jsonparser.GetString(existPayload, "path", "[0]", "suffix")
	p := &ics23.ExistenceProof{
		Key:   k,
		Value: val,
		Leaf: &ics23.LeafOp{
			Hash:         ics23.HashOp(hash),
			PrehashKey:   ics23.HashOp(hash),
			PrehashValue: ics23.HashOp(prehashValue),
			Length:       ics23.LengthOp(length),
			Prefix:       []byte(pre),
		},
		Path: []*ics23.InnerOp{{
			Hash:   ics23.HashOp(ophash),
			Prefix: []byte(opprefix),
			Suffix: []byte(opsuffix),
		}},
	}
	leafOp := make([]int32, 4)

	oppref, err := base64.RawStdEncoding.DecodeString(opprefix)
	opsuff, err := base64.RawStdEncoding.DecodeString(opsuffix)

	leafOp[0] = ics23.HashOp_value[p.Leaf.Hash.String()]
	leafOp[1] = ics23.HashOp_value[p.Leaf.PrehashKey.String()]
	leafOp[2] = ics23.HashOp_value[p.Leaf.PrehashValue.String()]
	leafOp[3] = ics23.LengthOp_value[p.Leaf.Length.String()]

	var innerOpHash int32

	key := string(p.Key)
	value := (string(p.Value))
	prefix, err := base64.RawStdEncoding.DecodeString(string(p.Leaf.Prefix))
	innerOpHash = ics23.HashOp_value[p.Path[0].Hash.String()]

	return &EncodePackedExistenceProof{
		LeafOp:        leafOp,
		InnerOpPrefix: oppref,
		InnerOpSuffix: opsuff,
		InnerOpHashOp: innerOpHash,
		Prefix:        prefix,
		Key:           key,
		Value:         value,
	}

}
