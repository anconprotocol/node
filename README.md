# Ancon Protocol Node v0.5.0


### Hybrid Smart Contracts protocol for secure offchain data economy

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
  "cid": {
    "/": "baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
  },
  "proof": {
    "/": "baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
  }
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
  "cid": {
    "/": "baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
  },
  "proof": {
    "/": "baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
  }
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
| `pin` | `bool` | ipfs pin |




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
   "proof": {
     ...
   },
   "signature": "...",
   "commitHash": "...",
   "issuer": {
     "/": "baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
   },
   "content": {
     "/": "baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
   },
   "timestamp": 101212365,
   "key": "/anconprotocol/baguqeeraui7hue3i2smgzmzdqmrxrnicqpoggayqkoocqdcjf3q5n66smdlq"
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




## `POST /v0/code`

> Uploads a hybrid smart contract


### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `code` | `string` | hex encoded Ancon Protocol Rust Smart Contract |
| `from` | `string` | DID |
| `signature` | `string` | hex encoded signature of code digest |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` |  content type stream|

example of the returned object:

```
<...data...>
```


## `POST /rpc`

> Trusted offchain JSON-RPC 2.0 gateway


## `ancon_call`

> Executes hybrid smart contract



### Parameters


| Name | Type | Description |
| ---- | ---- | ----------- |
| `to` | `string` | smart contract cid address |
| `from` | `string` | DID |
| `signature` | `string` | signature |
| `data` | `string` | data |



### Returns

| Type | Description |
| -------- | -------- |
| `Promise<Response>` |  content type stream|


### Example

```typescript
import dotenv from 'dotenv'
import { arrayify, hexlify } from '@ethersproject/bytes'
import fs from 'fs'
import { promisify } from 'util'
import { ethers } from 'ethers'
import web3 from 'web3'
import { Secp256k1 } from '@cosmjs/crypto'
import { fetchJson } from '@ethersproject/web'
import { toUtf8Bytes, toUtf8String } from 'ethers/lib/utils'
dotenv.config()

export class ExecuteContractExample {
  constructor(
    private privateKey: string,
    private apiUrl: string = 'http://localhost:7788',
  ) {}

  async execute({ to, from, query }: any): Promise<any> {
    let data = {
      to,
      from,
      data: query,
      // signature: '',
    } as any
    const bz = ethers.utils.toUtf8Bytes(data.to + data.from + data.data)
    const hash = ethers.utils.keccak256(bz)
    const sig = await Secp256k1.createSignature(
      arrayify(hash),
      Buffer.from(this.privateKey, 'hex'),
    )
    data.signature = hexlify(sig.toFixedLength())

    const body = {
      jsonrpc: '2.0',
      method: 'ancon_call',
      params: [data.to, data.from, data.signature, data.data],
      id: 1,
    }
    
    const result = await fetchJson(
      {
        url: `${this.apiUrl}/rpc`,
      },
      JSON.stringify(body),
    )

    return result
  }
}

;(async () => {
  const contract = new ExecuteContractExample(
    process.env.PRIVATE_KEY || '',
    process.env.URL,
  )

  //@ts-ignore
  let metadata = 'baguqeerak5hfvtgsvaaxtm5cbce2khnh7y4ijfnbrgeygwu4wixrz4z52vja'
  const result = await contract.execute({
    to: process.env.CONTRACT,
    from: process.env.DID_USER,
    query: `query   { metadata(cid: \"${metadata}\", path: \"/\"){name} }`
  })

  console.log(result)
  const j = (web3.utils.hexToString(result.result))
  console.log(JSON.parse(j))
  
})()
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

# Trusted offchain gateways

Ancon Protocol node can be used to integrate onchain and offchain sources using EIP-3668 Durin or also called Trusted Offchain gateway. Further in `Subgraph networks` chapter, we'll revisit this feature as we replace REST with Graphsync.

## What is trustless and trusted

A trustless setting onchain means the consensus of a blockchain is enough to validate a transaction is valid and has no bad behavior.

In cross chain use cases, there are many to accomplish this, one is with atomic swaps, which we'll use in parts, other is with ZK technology and other with protocols that are based on Merkle Proofs.


## Design and Architecture of a Hybrid Smart Contract

`Hybrid Smart Contracts` is the term used for integrating both offchain and onchain seamlessly in a secure way.

Ancon Protocol Node SDK uses a set of technologies, the developer should have a good grasp of the following:

- Go language
- IPLD
- GraphQL
- Rust

## Create Rust GraphQL Query and Mutations

Download the `github.com/anconprotocol/contracts` and label it with the name of your project.

The current source code has an example of a onchain DID ownership trasfer for ERC721 tokens.

```rust
use crate::sdk::focused_transform_patch_str;
use crate::sdk::read_dag;
use crate::sdk::submit_proof;
use crate::sdk::{generate_proof, get_proof, read_dag_block, write_dag_block};
use juniper::FieldResult;

extern crate juniper;

use juniper::{
    graphql_object, EmptyMutation, EmptySubscription, FieldError, GraphQLEnum, GraphQLValue,
    RootNode, Variables,
};
use serde::{Deserialize, Serialize};
use serde_json::json;

use std::collections::HashMap;

use std::str;
use std::vec::*;

pub struct Context {
    pub metadata: HashMap<String, Ancon721Metadata>,
    pub transfer: HashMap<String, MetadataPacket>,
}

impl juniper::Context for Context {}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct MetadataPacket {
    pub cid: String,
    pub from_owner: String,
    pub result_cid: String,
    pub to_owner: String,
    pub to_address: String,
    pub token_id: String,
    pub proof: String,
}

#[graphql_object(context = Context)]
impl MetadataPacket {
    fn cid(&self) -> &str {
        &self.cid
    }

    fn from_owner(&self) -> &str {
        &self.from_owner
    }

    fn result_cid(&self) -> &str {
        &self.result_cid
    }
    fn to_owner(&self) -> &str {
        &self.to_owner
    }

    fn to_address(&self) -> &str {
        &self.to_address
    }

    fn token_id(&self) -> &str {
        &self.token_id
    }
    fn proof(&self) -> &str {
        &self.proof
    }
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Ancon721Metadata {
    pub name: String,
    pub description: String,
    pub image: String,
    pub parent: String,
    pub owner: String,
    pub sources: Vec<String>,
}

#[graphql_object(context = Context)]
impl Ancon721Metadata {
    fn name(&self) -> &str {
        &self.name
    }

    fn description(&self) -> &str {
        &self.description
    }

    fn image(&self) -> &str {
        &self.image
    }
    fn parent(&self) -> &str {
        &self.parent
    }

    fn owner(&self) -> &str {
        &self.owner
    }

    async fn sources(&self) -> &Vec<String> {
        &self.sources
    }
}

#[derive(Clone, Copy, Debug)]
pub struct Query;

#[graphql_object(context = Context)]
impl Query {
    fn api_version() -> &'static str {
        "0.1"
    }

    pub fn metadata(context: &Context, cid: String, path: String) -> Ancon721Metadata {
        let v = read_dag(&cid);
        let res = serde_json::from_slice(&v);
        res.unwrap()
    }
}

#[derive(Clone, Copy, Debug)]
pub struct Mutation;

#[graphql_object(context = Context)]
impl Mutation {
    //Dagblock mutation
    fn transfer(context: &Context, input: MetadataTransactionInput) -> Vec<MetadataPacket> {
        let v = read_dag(&input.cid);
        let res = serde_json::from_slice(&v);
        let metadata: Ancon721Metadata = res.unwrap();

        //generate current metadata proof packet
        let proof = generate_proof(&input.cid);

        let updated_cid =
            focused_transform_patch_str(&input.cid, "owner", &metadata.owner, &input.new_owner);
        let updated =
            focused_transform_patch_str(&updated_cid, "parent", &metadata.parent, &input.cid);

        //generate updated metadata proof packet
        let proof_cid = apply_request_with_proof(input.clone(), &proof, &updated);
        let v = read_dag(&proof_cid);
        let res = serde_json::from_slice(&v);
        let packet: MetadataPacket = res.unwrap();
        let current_packet = MetadataPacket {
            cid: input.cid,
            from_owner: input.owner,
            result_cid: updated,
            to_address: "".to_string(),
            to_owner: input.new_owner,
            token_id: "".to_string(),
            proof: proof,
        };
        let result = vec![current_packet, packet];
        result
    }
}

#[derive(Clone, Debug, GraphQLInputObject, Serialize, Deserialize)]
struct MetadataTransactionInput {
    path: String,
    cid: String,
    owner: String,
    new_owner: String,
}

type Schema = RootNode<'static, Query, Mutation, EmptySubscription<Context>>;

pub fn schema() -> Schema {
    Schema::new(Query, Mutation, EmptySubscription::<Context>::new())
}

fn apply_request_with_proof(
    input: MetadataTransactionInput,
    prev_proof: &str,
    new_cid: &str,
) -> String {
    // Must combined proofs (prev and new) in host function
    // then send to chain and return result
    let js = json!({
        "previous": prev_proof,
        "next_cid": new_cid,
        "input": input
    });
    submit_proof(&js.to_string())
}
```


## DAG Operations or Mutation

Now we are going to apply a GraphQL mutation, which is an update to an immutable object, means it will discard the old CID and create a new CID from the latest update in the DAG block. 

Why are we abstracting on top of GraphQL? The main reason is to provide a better and more expedite approach to software engineering differents pieces of technologies like IPFS and blockchain. By enforcing the schemas with code generation in server side, we also get a similar developer experience as when you do smart contract development. 

In this example, we'll use one of the easiest IPLD Operator, which is the `focused transform`, where we 

- Pinpoint or **select** a path inside a root node
- Patch or mutate that selection with a function call. In our case, a diff patch, eg if previous node matches previous node requested a change, then apply requested change to node.

In Ancon Protocol Contracts SDK, you can use focused transform with `focused_transform_patch_str`

```rust

#[derive(Clone, Copy, Debug)]
pub struct Mutation;

#[graphql_object(context = Context)]
impl Mutation {
    //Dagblock mutation
    fn transfer(context: &Context, input: MetadataTransactionInput) -> Vec<MetadataPacket> {
        let v = read_dag(&input.cid);
        let res = serde_json::from_slice(&v);
        let metadata: Ancon721Metadata = res.unwrap();

        //generate current metadata proof packet
        let proof = generate_proof(&input.cid);

        let updated_cid =
            focused_transform_patch_str(&input.cid, "owner", &metadata.owner, &input.new_owner);
        let updated =
            focused_transform_patch_str(&updated_cid, "parent", &metadata.parent, &input.cid);

        //generate updated metadata proof packet
        let proof_cid = apply_request_with_proof(input.clone(), &proof, &updated);
        let v = read_dag(&proof_cid);
        let res = serde_json::from_slice(&v);
        let packet: MetadataPacket = res.unwrap();
        let current_packet = MetadataPacket {
            cid: input.cid,
            from_owner: input.owner,
            result_cid: updated,
            to_address: "".to_string(),
            to_owner: input.new_owner,
            token_id: "".to_string(),
            proof: proof,
        };
        let result = vec![current_packet, packet];
        result
    }
}

#[derive(Clone, Debug, GraphQLInputObject, Serialize, Deserialize)]
struct MetadataTransactionInput {
    path: String,
    cid: String,
    owner: String,
    new_owner: String,
}


```

The result must always be the previous and next packets.

## Smart Contract APIs

```rust
pub fn focused_transform_patch_str(cid: &str, path: &str, prev: &str, next: &str) -> String
```

Applies an IPLD focused transform using a patch design pattern for string node values

```rust
pub fn read_dag(cid: &str) -> Vec<u8>
```

Reads a cid from dag store

```rust
pub fn submit_proof(data: &str) -> String
```

Submits proof (offchain)

```rust
pub fn get_proof(cid: &str) -> String
```

Retrieves proof (offchain)

```rust
pub fn generate_proof(cid: &str) -> String
```

Generates proof (offchain)

> Copyright IFESA 2021, 2022
