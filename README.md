# ancon-ipld-router-sync


## Ancon Offchain Store and Graphsync node

Trusted storage used by `Ancon Protocol`. Requires CID to be anchor to onchain protocol to be used with `Ancon Data Contracts`

### API

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

'did': (Reserved) User DID
'path': Path
`data`: JSON data as Base64 encoded

#### Response

`cid`: a cid hash


### GET /dagcbor/:cid/*path

Fetches JSON

#### Request

`cid`: file cid
`path`: Path


### POST /dagcbor

Uploads CBOR data to be stored as DAG nodes.

#### Request

'did': (Reserved) User DID
'path': Path
`data`: CBOR data as Base64 encoded

#### Response

`cid`: a cid hash



### GET /dagcbor/:cid/*path

Fetches CBOR

#### Request

`cid`: file cid
`path`: Path

