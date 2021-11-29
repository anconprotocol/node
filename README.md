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

## Generating schemas

`~/Code/ancon-ipld-router-sync/x/anconsync/codegen$ go run github.com/99designs/gqlgen generate`


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

