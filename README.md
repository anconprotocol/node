# Ancon IPLD Router Sync


## Ancon Offchain Store and Graphsync node

Trusted storage used by `Ancon Protocol`. Requires CID to be anchor to onchain protocol to be used with `Ancon Data Contracts`
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
curl -X POST http://localhost:7788/dagjson -H "Content-Type: application/x-www-form-urlencoded" -d "did=did:web:ancon.did.pa:user:ifesa&path=/&data=eyJzY2hlbWEiOiAidHlwZSBRdWVyeSB7bWU6IFVzZXJ9IHR5cGUgVXNlciB7aWQ6IElEIG5hbWU6IFN0cmluZyB9In0=codec=dag-json"
```

### Read the schema dagjson by schema cid
``` bash
curl -X GET http://localhost:7788/dagjson/baguqeeraz4xl6e32jaq7xvnuirtszmhoburjgtjp4thrkrjkx4z53eul37rq/
```

### Write a base 64 encoded json payload

``` bash
curl -X POST http://localhost:7788/dagjson -H "Content-Type: application/x-www-form-urlencoded" -d "did=did:web:ancon.did.pa:user:ifesa&path=/&data=ewogICJkYXRhIjogewogICAgInVzZXIiOiBbCiAgICAgIHsKICAgICAgICAibmFtZSI6ICJBbGljaWEiCiAgICAgIH0sCiAgICAgIHsKICAgICAgICAibmFtZSI6ICJCb2IiCiAgICAgIH0KICAgIF0KICB9Cn0==codec=dag-json"
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