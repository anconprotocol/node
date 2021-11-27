
import AnconVerifier from "./AnconVerifier.cdc";

pub contract AtlanticoRouter {
  pub let verifier: Address;
  pub var headers: {String: [UInt8]};

  init(verifierLocation: Address) {
    self.verifier = verifierLocation;
    self.headers = {};
  }

  pub fun updateHeader(key: [UInt8], rootHash: [UInt8]): Bool {
    let encodedKey = String.encodeHex(key);
    self.headers[encodedKey] = rootHash;
    return true;
  }

  pub fun canSendMetadataOwnership(
    leafOpUint: [UInt8],
    prefix: [UInt8],
    existenceProofInnerOp: [[[UInt8]]],
    existenceProofInnerOpHash: UInt8,
    existenceProofKey: [UInt8],
    existenceProofValue: [UInt8],
    key: [UInt8],
    value: [UInt8],
  ): Bool {

    // Get the QueryRoot Resource
    let verifierAccount = getAccount(self.verifier);
    let queryRootCapability = verifierAccount.getCapability(/public/AnconQueryRoot);
    let queryRoot = queryRootCapability.borrow<&{AnconVerifier.IQueryRoot}>()
      ?? panic("QueryRoot resource is not available");

    let calculatedHash = queryRoot.queryRootCalculation(
      leafOpUint: leafOpUint,
      prefix: prefix,
      existenceProofInnerOp: existenceProofInnerOp,
      existenceProofInnerOpHash: existenceProofInnerOpHash,
      existenceProofKey: existenceProofKey,
      existenceProofValue: existenceProofValue
    );

    let currentKeyHash = String.encodeHex(HashAlgorithm.SHA3_256.hash(key));
    let currentHeader = self.headers[currentKeyHash] ?? [];

    let parsedCurrent = String.encodeHex(HashAlgorithm.SHA3_256.hash(currentHeader))
    let parsedCalculated = String.encodeHex(HashAlgorithm.SHA3_256.hash(calculatedHash));
    return parsedCurrent == parsedCalculated;
  }
}