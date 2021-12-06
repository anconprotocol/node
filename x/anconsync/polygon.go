package anconsync

import (
	"github.com/0xPolygon/polygon-sdk/state"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// 	// event AddOnchainMetadata(string memory name, string memory description, string indexed memory image, string memory owner, string memory parent, bytes memory sources)

// 	// MetadataTransferOwnershipEvent represent the signature of
// 	// `event InitiateMetadataTransferOwnership(address fromOwner, address toOwner, string memory metadataUri)`

func MetadataTransferOwnershipEvent() abi.Event {

	addressType, _ := abi.NewType("address", "", nil)
	stringType, _ := abi.NewType("string", "", nil)
	return abi.NewEvent(
		"InitiateMetadataTransferOwnership",
		"InitiateMetadataTransferOwnership",
		false,
		abi.Arguments{abi.Argument{
			Name:    "fromOwner",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "toOwner",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "metadataUri",
			Type:    stringType,
			Indexed: false,
		}},
	)
}

func PostTxProcessing(s Storage, t *state.Transition) error {
	for _, log := range t.Txn().Logs() {
		for _, topic := range log.Topics {
			if common.Hash(topic) != MetadataTransferOwnershipEvent().ID {
				continue
			}
		}
		if len(log.Topics) == 0 {
			continue
		}

		// if !ContractAllowed(log.Address) {
		// 	// Check the contract whitelist to prevent accidental native call.
		// 	continue
		// }
		_, err := MetadataTransferOwnershipEvent().Inputs.Unpack(log.Data)

		if err != nil {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func GetHooks(s Storage) func(t *state.Transition) {
	return func(t *state.Transition) {
		PostTxProcessing(s, t)
		// logs := t.Txn().Logs()

		// // for _, log := range logs {
		// // 	log.Topics
		// // }
	}
}
