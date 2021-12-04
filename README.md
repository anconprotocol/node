# Ancon Protocol Node (beta / testnet)

## Initial Beta Release

**Replaces CosmosSDK chain with a Graphsync network sync**

## Usage

1. Download latest release
2. Run `anconsync` with `anconsync -peeraddr <seed peer multiaddress> -addr <host multiaddress> -apiaddr <host API address> -data <data directory>`
3. Configure ports and firewall rules
4. Enjoy

## Features

### State of the art IPLD API engine

Uses the latest `go-ipld-prime` API

### DAG Contracts 

DAG Contracts enables a whole set of use cases that regular smart contracts might be too costly or expensive. In this first release, DAG contracts will work with GraphQL queries with immutable CID data.

Further down the road, GraphQL mutations will be used to integrate with onchain platforms, in an agnostic way.

An use case that will be particularly important is for mass or batch NFT metadata updates or creation. With a DAG Contract, you can patch with IPLD Selectors an existing CID, patch with parent CID and new owner, create the new metadata CID, and then send the event to an onchain function, that takes care of minting, or any other business logic. 

### Streaming files

Ancon IPLD Router Sync uses `go-ipld-prime` **fsstore** which is a streaming enable blockstore that is well performant. 

### GraphQL

We eventually decided to pick GraphQL for our DAG Contracts technology. With that, expect the following features:

#### DAG Contracts Directives

```graphql
query {
  me {
    @walk(resumeLink)   # Will traverse the link 
  }
}

query {
  me {
    @focus(resumeLink) 
  }
}
```

#### DAG Contracts Subscriptions

Any IPLD events will be mapped to GraphQL Subscriptions


#### DAG Contracts Onchain Adapters (GraphQL Mutations)

You will have a list of onchain adapters, with EVM adapters being the first to come out.




## Getting started
``` bash
go mod tidy
```
``` bash
go build ./server.go
```
``` bash
./server
```

## Generating schemas and Swagger

### Server GraphQL schemas

`~/Code/ancon-ipld-router-sync/x/anconsync/codegen$ go run github.com/99designs/gqlgen generate`

### Client GraphQL schemas

`go run github.com/Khan/genqlient --init`

### Swagger

`swag init`

## API

- Swagger: `https://ancon.did.pa/api/swagger/index.html`

## Examples

### Create DAG blocks

```
POST https://ancon.did.pa/api/v0/dagjson
Content-Type: application/json
```

```json
{
  "path": "/",
  "data": {
    "name": "XDV metadata sample",
    "owner": "alice",
    "description": "testing sample",
    "image": "https://explore.ipld.io/#/explore/QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D"
  }
}
```

```
POST https://ancon.did.pa/api/v0/dagcbor
Content-Type: application/json
```

```json
{
  "path": "/",
  "data": "<cbor content base64 encoded>"
}
```

### Reading blocks

```
GET https://ancon.did.pa/api/v0/dagcbor/baguqeeraouhd5jr7ktftgkbo5ufmihs3e2yqonajjejhomvbtsrwynlqdxba/
```

```
GET https://ancon.did.pa/api/v0/dagjson/baguqeeraouhd5jr7ktftgkbo5ufmihs3e2yqonajjejhomvbtsrwynlqdxba/file.json
```


### GraphQL DAG Designer

- Playground: `https://ancon.did.pa/api/v0/query`

**Query a dag-json by CID**

![Query](https://user-images.githubusercontent.com/1248071/143898301-b72b4ad3-9faf-459f-860f-52bd82283914.png)

**Mutates existing metadata CID**

![Mutation](https://user-images.githubusercontent.com/1248071/143898299-50164b0b-4a2a-4582-9c22-aa7e9f374e45.png)


### EVM Adapter for Durin / Secure Offchain

1. Fork off repo and customize `adapters/ethereum_adapter.go`, in our case a NFT ownership transfer. This must be plugged in to the gateway, either REST or JSON-RPC.

```go
func (adapter *EthereumAdapter) ExecuteDagContract(
	metadatadCid string,
	resultCid string,
	fromOwner string,
	toOwner string,
) (*DagTransaction, error) {

	pk, has := os.LookupEnv("ETHEREUM_ADAPTER_KEY")
	if !has {
		return nil, fmt.Errorf("environment key ETHEREUM_ADAPTER_KEY not found")
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, fmt.Errorf("invalid ETHEREUM_ADAPTER_KEY")
	}

	data, err := ExecuteDagContractWithProofAbiMethod().Inputs.Pack(metadatadCid, fromOwner, resultCid, toOwner)
	if err != nil {
		return nil, fmt.Errorf("packing for proof generation failed")
	}

	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("signing failed")
	}

	return &DagTransaction{
		MetadataCid: metadatadCid,
		ResultCid: resultCid,
		FromOwner: fromOwner,
		ToOwner: toOwner,
		Signature:     hexutil.Encode(signature),
	}, nil
}

```

Graphql integration executes the mutation and then the resulting cid is signed by the Durin gateway. 

2. Customize `DagContractTrusted.sol`. The signature proof needs to match the gateway signed proof.

```solidity
//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.4;

import "./IDagContractTrustedReceiver.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/Address.sol";

contract DagContractTrusted is Ownable {
    using ECDSA for bytes32;
    using Address for address;
    string public url;
    address private _signer;
    mapping(bytes32 => bool) executed;
    error OffchainLookup(string url, bytes prefix);
    struct DagContractRequestProof {
        string metadataCid;
        address fromOwner;
        address toOwner;
        string resultCid;
        address toReceiverContractAddress;
        bytes32 signature;
    }

    constructor() {}

    function setUrl(string memory url_) external onlyOwner {
        url = url_;
    }

    function setSigner(address signer_) external onlyOwner {
        _signer = signer_;
    }

    function getSigner() external view returns (address) {
        return _signer;
    }

    /**
     * @dev Requests a DAG contract offchain execution
     */
    function request(address toReceiverContractAddress, uint256 tokenId)
        public
        returns (bytes32)
    {
        revert OffchainLookup(
            url,
            abi.encodeWithSignature(
                "requestWithProof(address toReceiverContractAddress, uint256 tokenId, DagContractRequestProof memory proof)",
                toReceiverContractAddress,
                tokenId
            )
        );
    }

    /**
     * @dev Requests a DAG contract offchain execution with proof
     */
    function requestWithProof(
        address toReceiverContractAddress,
        uint256 tokenId,
        DagContractRequestProof memory proof
    ) external returns (bool) {
        if (executed[proof.signature]) {
            return false;
        } else {
            bytes32 digest = keccak256(
                abi.encodePacked(
                    "\x19Ethereum Signed Message:\n32",
                    keccak256(
                        abi.encodePacked(
                            proof.metadataCid,
                            proof.fromOwner,
                            proof.resultCid,
                            proof.toOwner,
                            toReceiverContractAddress,
                            tokenId
                        )
                    )
                )
            );

            address recovered = digest.recover(digest, proof.signature);

            require(
                _signer == recovered,
                "Signer is not the signer of the token"
            );
            executed[proof.signature] = true;
            bytes memory data = abi.encodePacked(
                toReceiverContractAddress,
                tokenId
            );
            _onDagContractResponseReceived(
                toReceiverContractAddress,
                address(this),
                msg.sender,
                proof.metadataCid,
                proof.resultCid,
                data
            );
            return true;
        }
    }

    /**
     * @dev Receives the DAG contract execution result
     */
    // function onDagContractResponseReceived(
    //     address operator,
    //     address from,
    //     string memory parentCid,
    //     string memory newCid,
    //     bytes calldata data
    // ) external returns (bytes4) {
    //     return IDagContractTrustedReceiver.onDagContractResponseReceived.selector;
    // }

    function _onDagContractResponseReceived(
        address to,
        address operator,
        address from,
        string memory parentCid,
        string memory newCid,
        bytes memory _data
    ) private returns (bool) {
        if (to.isContract()) {
            try
                IDagContractTrustedReceiver(to).onDagContractResponseReceived(
                    operator,
                    from,
                    parentCid,
                    newCid,
                    _data
                )
            returns (bytes4 retval) {
                return
                    retval ==
                    IDagContractTrustedReceiver
                        .onDagContractResponseReceived
                        .selector;
            } catch (bytes memory reason) {
                if (reason.length == 0) {
                    revert("DagContractTrusted: invalid receiver implementer");
                } else {
                    assembly {
                        revert(add(32, reason), mload(reason))
                    }
                }
            }
        } else {
            return true;
        }
    }
}

```

### Pending

- Tests
- Use case / demo


> Copyright IFESA 2021
