## Aguaclara - Ancon Protocol
### Cosmos - Near NFT Bridge


### Steps 

## Registry chain - Ancon Protocol

1. Once ChangeOwnership is commited, query for tx proof using ancon_getTransactionByHash

## Chain request

- NFT Token must exist in both chains
- Verifier contracts with [root hash installed]
- AnconMetadataOwnable
- onReceiver in NFT token

## Source EVM Chain - Verifier

1. Root Hash (Deployment)
2. Verify and Execute (changeOwnerWithProof)

## Source Execution onReceiver example

Chain id Rinkeby, XDVNFT, 0x1

## Destination EVM Chain - Verifier

1. Root Hash (Deployment)
2. Verify and Execute (changeOwnerWithProof)

## Destination Execution onReceiver example

Chain id Rinkeby, XDVNFT, 0x1

## Optional duplicate ownership 

XDVNFT BSC & Rinkeby ---> same onchain metadata (must have current owners)