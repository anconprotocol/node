# Ancon Protocol Node v0.5.0


### Protocol for secure offchain data economy

Ancon protocol is a new kind of SDK and technology that can be used to implement secure offchain data integrations using best of breed offchain protocols like ipfs and any blockchain with smart contracts support.


![AnconProtocolProducts](https://user-images.githubusercontent.com/1248071/147708647-f0e25a24-8c54-4a62-923e-5a73bb0c9e60.png)


## Ancon Protocol Node - L2 Gateway

Node manages offchain data integrations and trusted offchain gateways.  It has DID web and DID key, Graphsync, and dag-json / dag-cbor technology support.



## Usage

1. Download latest release
2. Create rootkey with `anconsync --init`
2. Production settings (recommended) `anconsync  --rootkey=<your new rootkey> --peeraddr /ip4/127.0.0.1/tcp/4001/p2p/<peer-id> --cors true  --origins=http://localhost:3000 --quic true --tlscert=/etc/letsencrypt/live/mynode/fullchain.pem --tlskey=/etc/letsencrypt/live/mynode/privkey.pem
~                                                                                                                                `
3. Configure ports and firewall rules
4. Enjoy

## Getting started
``` bash
go mod tidy
```
``` bash
go build ./main.go
```
``` bash
./main
```

If you have problems with the system buffer size this can help:
https://github.com/lucas-clemente/quic-go/wiki/UDP-Receive-Buffer-Size

# API Reference


## `POST /v0/did/key`

> Creates a new did-key


### Parameters

None


### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | An object that contains the CID |

example of the returned object:

```json
{
  "commitHash": "/AzWS9kE67z+wRs8htT3G+bRYDLy8V/Jg/cGUBprV/s=",
  "content": {
    "/": "baguqeerafkyyjhrgfai6x6djd23ot2d6vytaf35uvg6s2egc7llqkuc7nfua"
  },
  "height": 4892,
  "issuer": "0xeeC58E89996496640c8b5898A7e0218E9b6E90cB",
  "key": "L2FuY29ucHJvdG9jb2wvYmFmeXJlaWJxaXFiY2FmbnptanFtdjNpeTd1emppaW1uZWlxMmNxc3AzYm1odGNqYnJ3eXF3dnl3YmkvdXNlci9iYWd1cWVlcmFma3l5amhyZ2ZhaTZ4NmRqZDIzb3QyZDZ2eXRhZjM1dXZnNnMyZWdjN2xscWt1YzduZnVh",
  "parent": "/anconprotocol/bafyreibqiqbcafnzmjqmv3iy7uzjiimneiq2cqsp3bmhtcjbrwyqwvywbi/user",
  "signature": "0x971d3282785c390336860c5f5e5e1c7058f028738da7a003b8d81da7182cd6880798f8608e74d381c77d071f88adfa20e528bed1afba05f3c564bc6b59ec2dc61c",
  "timestamp": 1642350132
}
```



## `POST /v0/did/web`

> Creates a new did-web


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `domainName` | `string` | Subdomain eg alice.ipfs.pa |
| `pub` | `string` | (hex) public key  |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | An object that contains the CID |

example of the returned object:

```json
  {
    "commitHash": "/AzWS9kE67z+wRs8htT3G+bRYDLy8V/Jg/cGUBprV/s=",
  "content": {
    "/": "baguqeerafkyyjhrgfai6x6djd23ot2d6vytaf35uvg6s2egc7llqkuc7nfua"
  },
  "height": 4892,
  "issuer": "0xeeC58E89996496640c8b5898A7e0218E9b6E90cB",
  "key": "L2FuY29ucHJvdG9jb2wvYmFmeXJlaWJxaXFiY2FmbnptanFtdjNpeTd1emppaW1uZWlxMmNxc3AzYm1odGNqYnJ3eXF3dnl3YmkvdXNlci9iYWd1cWVlcmFma3l5amhyZ2ZhaTZ4NmRqZDIzb3QyZDZ2eXRhZjM1dXZnNnMyZWdjN2xscWt1YzduZnVh",
  "parent": "/anconprotocol/bafyreibqiqbcafnzmjqmv3iy7uzjiimneiq2cqsp3bmhtcjbrwyqwvywbi/user",
  "signature": "0x971d3282785c390336860c5f5e5e1c7058f028738da7a003b8d81da7182cd6880798f8608e74d381c77d071f88adfa20e528bed1afba05f3c564bc6b59ec2dc61c",
  "timestamp": 1642350132
  }
```


## `GET /v0/did/:did`

> Returns did document as json


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `did` | `string` | DID Doc id |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | A did document |

example of the returned object:

```javascript
{
  "@context": ["https://www.w3.org/ns/did/v1"],
  "authentication": [
    "",
    {
      "controller": "did:key:z3rc1YQMG366ttmuHeX2qodNeZAEhU6ktdjJEdLMRGX9gtpjaHitW6eu4BvZMEF",
      "id": "did:key:z3rc1YQMG366ttmuHeX2qodNeZAEhU6ktdjJEdLMRGX9gtpjaHitW6eu4BvZMEF#",
      "publicKeyBase58": "J8v1rHsHjrBP9khKJdiZBrfq4u2Ame2aUsv8fACmKjaF",
      "type": "Ed25519VerificationKey2018"
    }
  ],
  "created": "2021-12-04T07:57:33.030203855-05:00",
  "id": "did:key:z3rc1YQMG366ttmuHeX2qodNeZAEhU6ktdjJEdLMRGX9gtpjaHitW6eu4BvZMEF",
  "updated": "2021-12-04T07:57:33.030203855-05:00",
  "verificationMethod": [
    {
      "controller": "did:key:z3rc1YQMG366ttmuHeX2qodNeZAEhU6ktdjJEdLMRGX9gtpjaHitW6eu4BvZMEF",
      "id": "did:key:z3rc1YQMG366ttmuHeX2qodNeZAEhU6ktdjJEdLMRGX9gtpjaHitW6eu4BvZMEF#",
      "publicKeyBase58": "J8v1rHsHjrBP9khKJdiZBrfq4u2Ame2aUsv8fACmKjaF",
      "type": "Ed25519VerificationKey2018"
    }
  ]
}

```

## `GET /proof/:key?height=n&exportAs=qr`

> Gets a proof given a key and height


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `key` | `string` | Key(base 64) |
| `height` | `int` | Height (int) |
| `exportAs` | `string` | Export as: qr |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | An object that contains the hash |


## `POST /proofs/qr`

> Translate  a QR to a proof


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `file` | Image blob | Image QR generated by Ancon |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | An object that contains the hash |


## `GET /proofs/lasthash`

> Reads current last hash


### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | An object that contains the hash |

## `GET /user/:did/did.json`

> Reads a did-web


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `did` | `string` | did web domain name   |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | An object that contains the CID |

example of the returned object:

```json
{
  "@context": ["https://www.w3.org/ns/did/v1"],
  "authentication": [
    "",
    {
      "controller": "did:web:ipfs:user:rogelio",
      "id": "did:web:ipfs:user:rogelio",
      "publicKeyBase58": "ER5jUmbiApGWtR4QVHjG7nHpaMGhZmf4BRMSLw4tBeEmT8RZhUKwppqsjHihwXp4RpVjVXaChRZFyKj1s41uGJs",
      "type": "Secp256k1VerificationKey2018"
    }
  ],
  "created": "2021-12-04T08:20:35.623500975-05:00",
  "id": "did:web:ipfs:user:rogelio",
  "updated": "2021-12-04T08:20:35.623500975-05:00",
  "verificationMethod": [
    {
      "controller": "did:web:ipfs:user:rogelio",
      "id": "did:web:ipfs:user:rogelio",
      "publicKeyBase58": "ER5jUmbiApGWtR4QVHjG7nHpaMGhZmf4BRMSLw4tBeEmT8RZhUKwppqsjHihwXp4RpVjVXaChRZFyKj1s41uGJs",
      "type": "Secp256k1VerificationKey2018"
    }
  ]
}

```

## `PUT /v0/dagjson`
## `PUT /v0/dag`

> Mutates a dag-json in users path. Must have registerd DID and messasge must be signed with signature matching DID.


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `from` | `string` | DID identifier |
| `signature` | `string` | signature as hex |
| `data` | `object` | Mutations |
| `cid` | `string` | cid to mutate |


### Mutations

Mutations only apply to the current cid, and is executed sequentially. Each sequence cid stored in `parent` property.

```json
{ 
  ...,
  "data": [
    {
      "path": "content/royalty",
      "previousValue": 0.1,
      "nextValue": 1,
    },
    {
      "path": "content/owner",
      "previousValue": "alice",
      "nextValue": "bob",
    },
    {
      // will add a new node
      "path": "tag",
      "previousValue": null,
      "nextValue": "nft_from_panama",
    },
  ],
  ...,
}
```


### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | An object that contains the CID |

example of the returned object:

```json
{
  "cid": {
    "/": "baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
  }
}
```


## `POST /v0/dagjson`
## `POST /v0/dag`

> Stores json as dag-json in users path. Must have registerd DID and messasge must be signed with signature matching DID.


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `path` | `string` | path |
| `from` | `string` | DID identifier |
| `signature` | `string` | signature as hex |
| `data` | `object` | object to store |
| `encrypt` | `bool` | enables JOSE Web Encryption |
| `authorizedRecipients` | `string array` | comma delimited Ethereum address



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` | An object that contains the CID |

example of the returned object:

```json
{
  "cid": {
    "/": "baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
  }
}
```




## `GET /v0/dagjson/:cid/*path`
## `GET /v0/dag/:cid/*path`

> Reads a dag-json block


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `cid` | `string` | cid |
| `path` | `string` | path |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<DagResponse>` |  json object |



```json
 {
  "commitHash": "/AzWS9kE67z+wRs8htT3G+bRYDLy8V/Jg/cGUBprV/s=",
  "content": {
    "/": "baguqeerafkyyjhrgfai6x6djd23ot2d6vytaf35uvg6s2egc7llqkuc7nfua"
  },
  "height": 4892,
  "issuer": "0xeeC58E89996496640c8b5898A7e0218E9b6E90cB",
  "key": "L2FuY29ucHJvdG9jb2wvYmFmeXJlaWJxaXFiY2FmbnptanFtdjNpeTd1emppaW1uZWlxMmNxc3AzYm1odGNqYnJ3eXF3dnl3YmkvdXNlci9iYWd1cWVlcmFma3l5amhyZ2ZhaTZ4NmRqZDIzb3QyZDZ2eXRhZjM1dXZnNnMyZWdjN2xscWt1YzduZnVh",
  "parent": "/anconprotocol/bafyreibqiqbcafnzmjqmv3iy7uzjiimneiq2cqsp3bmhtcjbrwyqwvywbi/user",
  "signature": "0x971d3282785c390336860c5f5e5e1c7058f028738da7a003b8d81da7182cd6880798f8608e74d381c77d071f88adfa20e528bed1afba05f3c564bc6b59ec2dc61c",
  "timestamp": 1642350132
}
```


## `GET /v0/file/:cid/*path`

> Reads a dag-json block


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `cid` | `string` | cid |
| `path` | `string` | path |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` |  content type stream|

example of the returned object:

```
<...data...>
```


## `GET /swagger`

> REST documentation


# CLI Reference

- `peeraddr`: Connects to subgraph node with IPFS
- `addr`:  Host libp2p address
- `apiaddr`:  Host API address
- `data`: Storage directory
- `cors`: Set to true to enable CORS requests
- `origins`: Comma separated list of addresses
- `init`: Initializes the proof storage by creating a genesis block
- `keys`: Generates Secp256k1 keys
- `hostname`: Node identifier
- `rootkey`: Rootkey to validate
- `sync`: Syncs with peers
- `peers`:  List of peers to sync
- `quic`: Enables QUIC
- `tlscert`: TLS certificate for QUIC
- `tlskey`: TLS key for QUIC
- `ipfshost`: IPFS Host address for DAG Pinning

# Trusted offchain gateways

Ancon Protocol node can be used to integrate onchain and offchain sources using EIP-3668 Durin or also called Trusted Offchain gateway. Further in `Subgraph networks` chapter, we'll revisit this feature as we replace REST with Graphsync.

## What is trustless and trusted

A trustless setting onchain means the consensus of a blockchain is enough to validate a transaction is valid and has no bad behavior.

In cross chain use cases, there are many to accomplish this, one is with atomic swaps, which we'll use in parts, other is with ZK technology and other with protocols that are based on Merkle Proofs.



> Copyright IFESA 2021, 2022
