pragma solidity >=0.8.0;
import "./Bytes.sol";

library Ics23Helper {
  using Bytes for bytes;
  // Data structures and helper functions

  enum HashOp {
    NO_HASH,
    SHA256,
    SHA512,
    KECCAK,
    RIPEMD160,
    BITCOIN
  }
  enum LengthOp {
    NO_PREFIX,
    VAR_PROTO,
    VAR_RLP,
    FIXED32_BIG,
    FIXED32_LITTLE,
    FIXED64_BIG,
    FIXED64_LITTLE,
    REQUIRE_32_BYTES,
    REQUIRE_64_BYTES
  }

  struct ExistenceProof {
    bool valid;
    bytes key;
    bytes value;
    LeafOp leaf;
    InnerOp[] path;
  }

  struct NonExistenceProof {
    bool valid;
    bytes key;
    ExistenceProof left;
    ExistenceProof right;
  }

  struct LeafOp {
    bool valid;
    HashOp hash;
    HashOp prehash_key;
    HashOp prehash_value;
    LengthOp len;
    bytes prefix;
  }

  struct InnerOp {
    bool valid;
    HashOp hash;
    bytes prefix;
    bytes suffix;
  }

  struct ProofSpec {
    LeafOp leafSpec;
    InnerSpec innerSpec;
    uint256 maxDepth;
    uint256 minDepth;
  }

  struct InnerSpec {
    uint256[] childOrder;
    uint256 childSize;
    uint256 minPrefixLength;
    uint256 maxPrefixLength;
    bytes emptyChild;
    HashOp hash;
  }

  
}