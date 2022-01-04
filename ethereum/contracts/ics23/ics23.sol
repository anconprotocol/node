pragma solidity >=0.8.0;
import "./Bytes.sol";
import "./Ics23Helper.sol";

contract ICS23 {
  using Bytes for bytes;
  // using Ics23Helper for *;
  // Ics23Helper Ics23Helper;
  // Ics23Helper.ProofSpec ProofSpec;
  // Ics23Helper.LeafOp LeafOp;
  // Ics23Helper.HashOp HashOp;
  // Ics23Helper.LengthOp LengthOp;
  // Ics23Helper.InnerOp InnerOp;
  // Ics23Helper.ExistenceProof ExistenceProof;
  // Ics23Helper.InnerSpec InnerSpec;
  // Data structures and helper functions

  function getIavlSpec() public pure returns (Ics23Helper.ProofSpec memory) {
    Ics23Helper.ProofSpec memory spec;

    uint256[] memory childOrder = new uint256[](2);
    childOrder[0] = 0;
    childOrder[1] = 1;

    spec.leafSpec = Ics23Helper.LeafOp(
      true,
      Ics23Helper.HashOp.SHA256,
      Ics23Helper.HashOp.NO_HASH,
      Ics23Helper.HashOp.SHA256,
      Ics23Helper.LengthOp.VAR_PROTO,
      hex"00"
    );

    spec.innerSpec = Ics23Helper.InnerSpec(childOrder, 33, 4, 12, "", Ics23Helper.HashOp.SHA256);
    return spec;
  }

  enum Ordering {
    LT,
    EQ,
    GT
  }

  function checkAgainstSpecLeafOp(Ics23Helper.LeafOp memory op, Ics23Helper.ProofSpec memory spec)
    internal
    pure
  {
    require(op.hash == spec.leafSpec.hash, "Unexpected HashOp");
    require(
      op.prehash_key == spec.leafSpec.prehash_key,
      "Unexpected PrehashKey"
    );
    require(
      op.prehash_value == spec.leafSpec.prehash_value,
      "Unexpected PrehashKey"
    );
    require(op.len == spec.leafSpec.len, "UnexpecteleafSpec LengthOp");
    require(
      Bytes.hasPrefix(op.prefix, spec.leafSpec.prefix),
      "LeafOpLib: wrong prefix"
    );
  }

  function applyValueLeafOp(
    Ics23Helper.LeafOp memory op,
    bytes memory key,
    bytes memory value
  ) internal pure returns (bytes memory) {
    require(key.length > 0, "Leaf op needs key");
    require(value.length > 0, "Leaf op needs value");
    bytes memory data = abi.encodePacked(
      abi.encodePacked(op.prefix, prepareLeafData(op.prehash_key, op.len, key)),
      prepareLeafData(op.prehash_value, op.len, value)
    );
    return doHash(op.hash, data);
  }

  function prepareLeafData(
    Ics23Helper.HashOp hashOp,
    Ics23Helper.LengthOp lengthOp,
    bytes memory data
  ) private pure returns (bytes memory) {
    return doLength(lengthOp, doHashOrNoop(hashOp, data));
  }

  function doHashOrNoop(Ics23Helper.HashOp hashOp, bytes memory data)
    private
    pure
    returns (bytes memory)
  {
    if (hashOp == Ics23Helper.HashOp.NO_HASH) {
      return data;
    }
    return doHash(hashOp, data);
  }

  function checkAgainstSpecInnerOp(Ics23Helper.InnerOp memory op, Ics23Helper.ProofSpec memory spec)
    internal
    pure
  {
    require(op.hash == spec.leafSpec.hash, "Unexpected HashOp");
    require(
      !Bytes.hasPrefix(op.prefix, spec.leafSpec.prefix),
      "InnerOpLib: wrong prefix"
    );
    require(
      op.prefix.length >= uint256(spec.innerSpec.minPrefixLength),
      "InnerOp prefix too short"
    );

    uint256 maxLeftChildLen = (spec.innerSpec.childOrder.length - 1) *
      uint256(spec.innerSpec.childSize);
    require(
      op.prefix.length <=
        uint256(spec.innerSpec.maxPrefixLength) + maxLeftChildLen,
      "InnerOp prefix too short"
    );
  }

  function applyValueInnerOp(Ics23Helper.InnerOp memory op, bytes memory child)
    internal
    pure
    returns (bytes memory)
  {
    require(child.length > 0, "Inner op needs child value");
    return
      doHash(
        op.hash,
        abi.encodePacked(abi.encodePacked(op.prefix, child), op.suffix)
      );
  }

  function doHash(Ics23Helper.HashOp hashOp, bytes memory data)
    internal
    pure
    returns (bytes memory)
  {
    if (hashOp == Ics23Helper.HashOp.SHA256) {
      return Bytes.fromBytes32(sha256(data));
    }

    if (hashOp == Ics23Helper.HashOp.SHA512) {
      //TODO: implement sha512
      revert("SHA512 not implemented");
    }

    if (hashOp == Ics23Helper.HashOp.RIPEMD160) {
      return Bytes.fromBytes32(ripemd160(data));
    }

    if (hashOp == Ics23Helper.HashOp.BITCOIN) {
      bytes32 hash = sha256(data);
      return Bytes.fromBytes32(ripemd160(Bytes.fromBytes32(hash)));
    }
    revert("Unsupported hashop");
  }

  function doLength(Ics23Helper.LengthOp lengthOp, bytes memory data)
    internal
    pure
    returns (bytes memory)
  {
    if (lengthOp == Ics23Helper.LengthOp.NO_PREFIX) {
      return data;
    }
    if (lengthOp == Ics23Helper.LengthOp.VAR_PROTO) {
      return abi.encodePacked(encodeVarintProto(uint64(data.length)), data);
    }
    if (lengthOp == Ics23Helper.LengthOp.REQUIRE_32_BYTES) {
      require(data.length == 32, "Expected 32 bytes");
      return data;
    }
    if (lengthOp == Ics23Helper.LengthOp.REQUIRE_64_BYTES) {
      require(data.length == 64, "Expected 64 bytes");
      return data;
    }
    revert("Unsupported lengthop");
  }

  function encodeVarintProto(uint64 n) internal pure returns (bytes memory) {
    // Count the number of groups of 7 bits
    // We need this pre-processing step since Solidity doesn't allow dynamic memory resizing
    uint64 tmp = n;
    uint64 num_bytes = 1;
    while (tmp > 0x7F) {
      tmp = tmp >> 7;
      num_bytes += 1;
    }

    bytes memory buf = new bytes(num_bytes);

    tmp = n;
    for (uint64 i = 0; i < num_bytes; i++) {
      // Set the first bit in the byte for each group of 7 bits
      buf[i] = bytes1(0x80 | uint8(tmp & 0x7F));
      tmp = tmp >> 7;
    }
    // Unset the first bit of the last byte
    buf[num_bytes - 1] &= 0x7F;

    return buf;
  }

  function verify(
    Ics23Helper.ExistenceProof memory proof,
    Ics23Helper.ProofSpec memory spec,
    bytes memory root,
    bytes memory key,
    bytes memory value
  ) public pure {
    checkAgainstSpec(proof, spec);
    require(Bytes.equals(proof.key, key), "Provided key doesn't match proof");
    require(
      Bytes.equals(proof.value, value),
      "Provided value doesn't match proof"
    );
    require(
      Bytes.equals(calculate(proof), root),
      "Calculcated root doesn't match provided root"
    );
  }

  function checkAgainstSpec(Ics23Helper.ExistenceProof memory proof, Ics23Helper.ProofSpec memory spec)
    private
    pure
  {
    checkAgainstSpecLeafOp(proof.leaf, spec);
    require(
      spec.minDepth == 0 || proof.path.length >= uint256(spec.minDepth),
      "InnerOps depth too short"
    );
    require(
      spec.maxDepth == 0 || proof.path.length >= uint256(spec.maxDepth),
      "InnerOps depth too short"
    );

    for (uint256 i = 0; i < proof.path.length; i++) {
      checkAgainstSpecInnerOp(proof.path[i], spec);
    }
  }

  // Calculate determines the root hash that matches the given proof.
  // You must validate the result is what you have in a header.
  // Returns error if the calculations cannot be performed.
  function calculate(Ics23Helper.ExistenceProof memory p)
    internal
    pure
    returns (bytes memory)
  {
    // leaf step takes the key and value as input
    bytes memory res = applyValueLeafOp(p.leaf, p.key, p.value);

    // the rest just take the output of the last step (reducing it)
    for (uint256 i = 0; i < p.path.length; i++) {
      res = applyValueInnerOp(p.path[i], res);
    }
    return res;
  }
}