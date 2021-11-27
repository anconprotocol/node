# Atlantico - Ancon Protocol

Atlantico is the Flow implementation of the Ancon protocol, written in Cadence. It allows Document NFTs to be transferred to the Flow Blockchain.

## To Test

1. `npm i`
2. `npm test`

All tests should be green just after installing.

## To deploy (to emulator)

The `./flow.json` is already ready for deployment to the Flow Emulator. The configuration for Testnet and Mainnet are still being prepared.

1. [Install the Flow Emulator](https://docs.onflow.org/flow-cli/install/)
2. On a console, run `flow emulator --storage-limit=false` and let it running.
3. On another console, run `flow project deploy --network=emulator`.

Deployment should succeed without errors. And now you can run the scripts and transactions from console against the emulator.

## Contract Explanations

Most of the bridge's functionality is implemented in the contracts included in `./contracts`.

### AnconVerifier Contract
This is a translation from the Solidity `AnconVerifier` to Cadence. The logic should be the same in both environments. The main difference is the `QueryRoot` resource. It is published as a public capability, so that anyone can call the `queryRoot` function anywhere in Flow.

### AtlanticoRouter
This is a replacement of `AnconRouter.sol` and contains part of its logic. It cannot be an exact translation of the Solidity contract because Flow and Cadence require a different arquitecture. See `./transactions/AtlanticoRouterTransfer.cdc` for the remaining pieces.

### NonFungibleToken Contract
This is an exact copy of the Cadence's Standard. In testnet and mainnet it will be replaced with the version already deployed via "aliases". See `./flow.json`'s `contracts.NonFungibleToken.aliases` for the relevant configuration.

### XDVNFT Contract
This is pretty similar to the Solidity Version, just translated to Cadence's Resource paradigm. For details on on how it works, please see the [official tutorial](https://docs.onflow.org/cadence/tutorial/04-non-fungible-tokens/).  Documents are represented as a standard NFT Resources and can be dealt with via Transactions.

**Important**: The Minter is a Public Capability! This is part of the XDV Design.

## Transactions
Some extra required functionality is implemented as Cadence Transactions
### AtlanticoRouterTransfer

This is expected to replace Solidity's `AguaclaraRouter.sendMetadataOwnership` and has the same signature. This executes the Query Root calculation on the `AtlanticoRouter` contract, and if the verification succees, then it transfers the NFT Token from the Sender to the Recipient.

### SetupAccount

This transaction must be executed at least once by anyone who wants to receive XDV NFT Tokens. This is implemented from the Official NFT Tutorial.

## Expected Flow:

1. "Recipient" must execute the `SetupAccount.cdc` to build the XDV NFT Collection in its account. This must only be done the first time Recipient interacts with the platform.
1. "Sender" must already have the XDV NFT.
1. "Sender" executes the `AtlanticoRouterTransfer.cdc` transaction. It is designed to be a direct replacement of Solidity's `AguaclaraRouter.sendMetadataOwnership` and the parameters should be the same.
