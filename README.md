# Ancon IPLD Router Sync

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
## API

### POST /file

File upload, supports streaming.

#### Request

`file`: a file to upload

#### Response

`cid`: a cid hash


### GET /file/:cid

Fetches a file

#### Request

`cid`: a file to download

#### Response

A streaming response, set `Content-Type` to get the blob of the file.


### POST /dagjson

Uploads JSON data to be stored as DAG nodes.

#### Request

- `did`: (Reserved) User DID
- `path`: Path
- `data`: JSON data as Base64 encoded

#### Response

`cid`: a cid hash


### GET /dagcbor/:cid/*path

Fetches JSON

#### Request

- `cid`: file cid
- `path`: Path


### POST /dagcbor

Uploads CBOR data to be stored as DAG nodes.

#### Request

- `did`: (Reserved) User DID
- `path`: Path
- `data`: CBOR data as Base64 encoded

#### Response

`cid`: a cid hash



### GET /dagcbor/:cid/*path

Fetches CBOR

#### Request

- `cid`: file cid
- `path`: Path

## Examples
### Sample payload
``` json
{
  "data": {
    "user": [
      {
        "name": "Alicia"
      },
      {
        "name": "Bob"
      }
    ]
  }
}
```
Encode as base 64
### Sample schema
``` bash
{"schema": "type Query {me: User} type User {id: ID name: String }"}
```
Encode as base 64

### Write a base 64 encoded json schema as dagjson

``` bash
curl -X POST http://localhost:7788/dagjson -H "Content-Type: application/x-www-form-urlencoded" -d "did=did:web:ancon.did.pa:user:ifesa&path=/&data=eyJzY2hlbWEiOiAidHlwZSBRdWVyeSB7bWU6IFVzZXJ9IHR5cGUgVXNlciB7aWQ6IElEIG5hbWU6IFN0cmluZyB9In0="
```

### Read the schema dagjson by schema cid
``` bash
curl -X GET http://localhost:7788/dagjson/baguqeeraz4xl6e32jaq7xvnuirtszmhoburjgtjp4thrkrjkx4z53eul37rq/
```

### Write a base 64 encoded json payload

``` bash
curl -X POST http://localhost:7788/dagjson -H "Content-Type: application/x-www-form-urlencoded" -d "did=did:web:ancon.did.pa:user:ifesa&path=/&data=ewogICJkYXRhIjogewogICAgInVzZXIiOiBbCiAgICAgIHsKICAgICAgICAibmFtZSI6ICJBbGljaWEiCiAgICAgIH0sCiAgICAgIHsKICAgICAgICAibmFtZSI6ICJCb2IiCiAgICAgIH0KICAgIF0KICB9Cn0=="
```

### Read the payload dagjson by payload cid 
``` bash
curl -X GET http://localhost:7788/dagjson/baguqeerakjuzcjvxkwpagivz7bnf45j4zsr7cbmu5263hk7j3467cq7wyyja/
```

### Write a graph ql query via schema cid, payload cid & query

``` bash
curl -X POST http://localhost:7788/gql -H "Content-Type: application/x-www-form-urlencoded" -d "op=&variables=&schemacid=baguqeeraz4xl6e32jaq7xvnuirtszmhoburjgtjp4thrkrjkx4z53eul37rq&payloadcid=baguqeerakjuzcjvxkwpagivz7bnf45j4zsr7cbmu5263hk7j3467cq7wyyja&query=query{me{name}}"
```

### Read the graphql dagjson query response via cid

``` bash
curl -X GET http://localhost:7788/dagjson/baguqeerardwpsizg6h7y3i6yd2zy4guekkg6kzq3p3bqxck3vawegemdqcxq/
```
