import XDVNFT from "../contracts/XDVNFT.cdc";

// This transaction configures a user's account
// to use the NFT contract by creating a new empty collection,
// storing it in their account storage, and publishing a capability
transaction {
  prepare(acct: AuthAccount) {
    // Create a new empty collection
    let collection <- XDVNFT.createEmptyCollection()

    // store the empty NFT Collection in account storage
    acct.save<@XDVNFT.Collection>(<-collection, to: /storage/XDVNFTCollection)

    log("Collection created for account 1")

    // create a public capability for the Collection
    acct.link<&{XDVNFT.NFTReceiver}>(/public/XDVNFTCollection, target: /storage/XDVNFTCollection)

    log("Capability created")
  }
}