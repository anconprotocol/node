import NonFungibleToken from "./NonFungibleToken.cdc";

// Note: Contract inspired by the official NFT implementation:
// https://github.com/onflow/flow-nft/blob/master/contracts/ExampleNFT.cdc
pub contract XDVNFT: NonFungibleToken {
  pub var totalSupply: UInt64;

  pub event ContractInitialized();
  pub event Withdraw(id: UInt64, from: Address?);
  pub event Deposit(id: UInt64, to: Address?);

  pub resource NFT: NonFungibleToken.INFT {
    pub let id: UInt64;
    pub let metadata: String;

    init(initId: UInt64, metadata: String) {
      self.id = initId;
      self.metadata = metadata;
    }
  }

  pub resource Collection: NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic {
    pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT};

    init () {
      self.ownedNFTs <- {};
    }

    destroy () {
      destroy self.ownedNFTs;
    }

    pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
      let token <- self.ownedNFTs.remove(key: withdrawID)
        ?? panic("missing NFT");

      emit Withdraw(id: token.id, from: self.owner?.address);
      return <- token;
    }

    pub fun deposit(token: @NonFungibleToken.NFT) {
      let id = token.id;
      let oldToken <- self.ownedNFTs[id] <- token;

      emit Deposit(id: id, to: self.owner?.address);

      destroy oldToken;
    }

    pub fun getIDs(): [UInt64] {
      return self.ownedNFTs.keys;
    }

    pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
      return &self.ownedNFTs[id] as! &NonFungibleToken.NFT;
    }
  }

  pub resource interface IDocumentMinter {
    pub fun mintNFT(recipient: &{NonFungibleToken.CollectionPublic}, metadata: String);
  }

  pub resource DocumentMinter: IDocumentMinter {
    pub fun mintNFT(recipient: &{NonFungibleToken.CollectionPublic}, metadata: String) {
      var newDocument <- create NFT(
        initId: XDVNFT.totalSupply,
        metadata: metadata
      );
      recipient.deposit(token: <- newDocument);
      XDVNFT.totalSupply = XDVNFT.totalSupply + (1 as UInt64);
    }
  }

  init () {
    self.totalSupply = 0;

    let collection <- create Collection();
    self.account.save(<- collection, to: /storage/XDVNFTCollection)

    self.account.link<&{NonFungibleToken.CollectionPublic}>(
      /public/XDVNFTCollection,
      target: /storage/XDVNFTCollection
    );

    let minter <- create DocumentMinter();
    self.account.save(<- minter, to:/storage/XDVNFTMinter)
    // NOTE: Minter has public access!
    self.account.link<&{IDocumentMinter}>(
      /public/XDVNFTMinter,
      target: /storage/XDVNFTMinter
    );

    emit ContractInitialized();
  }

  pub fun createEmptyCollection(): @NonFungibleToken.Collection {
    return <- create Collection();
  }
}
