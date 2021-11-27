import XDVNFT from "../contracts/XDVNFT.cdc";
import AtlanticoRouter from "../contracts/AtlanticoRouter.cdc";

transaction(
  leafOpUint: [UInt8],
  prefix: [UInt8],
  existenceProofInnerOp: [[[UInt8]]],
  existenceProofInnerOpHash: UInt8,
  existenceProofKey: [UInt8],
  existenceProofValue: [UInt8],
  key: [UInt8],
  value: [UInt8],
  tokenId: UInt64,
  recipient: Address
){

  let document: @XDVNFT.NFT;

  prepare(sender: AuthAccount) {
    let allowed = AtlanticoRouter.canSendMetadataOwnership(
      leafOpUint: leafOpUint,
      prefix: prefix,
      existenceProofInnerOp: existenceProofInnerOp,
      existenceProofInnerOpHash: existenceProofInnerOpHash,
      existenceProofKey: existenceProofKey,
      existenceProofValue: existenceProofValue,
      key: key,
      value: value
    );

    if (!allowed) {
      panic("Invalid proof for key")
    }

    let collectionRef: &XDVNFT.Collection = sender.borrow<&{XDVNFT.Collection}>(/storage/XDVNFTCollection)
      ?? panic("Sender does not have an XDV Collection");
    self.document <- collectionRef.withdraw(withdrawID: tokenId) as! @XDVNFT.NFT;
  }

  execute {
    let recipientAccount = getAccount(recipient);
    let recipientCollectionCapability = recipientAccount.getCapability(/public/XDVNFTCollection);
    let recipientCollection: &XDVNFT.Collection = recipientCollectionCapability.borrow<&{XDVNFT.Collection}>()
      ?? panic("Could not borrow Recipient's collection capability");

    recipientCollection.deposit(token: <- self.document);
  }
}