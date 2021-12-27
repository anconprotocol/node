// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "./ics23.sol";

/// @title A title that should describe the contract/interface
/// @author The name of the author
/// @notice Explain to an end user what this does
/// @dev Explain to a developer any extra details
contract AnconVerifier is ICS23 {
    address public owner;
    bytes relayNetworkHash;

    constructor(address onlyOwner) public {
        owner = onlyOwner;
    }

    function setRootHash(bytes memory rootHash)
        public
        returns (bool)
    {
        require(msg.sender == owner);
        relayNetworkHash = rootHash;
        return true;
    }
    function convertProof(
        bytes memory key,
        bytes memory value,
        bytes memory _prefix,
        uint256[] memory _leafOpUint,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix,
        uint256 existenceProofInnerOpHash
    ) public pure returns (ExistenceProof memory) {
        LeafOp memory leafOp = LeafOp(
            true,
            HashOp((_leafOpUint[0])),
            HashOp((_leafOpUint[1])),
            HashOp((_leafOpUint[2])),
            LengthOp((_leafOpUint[3])),
            _prefix
        );

        // // innerOpArr
        InnerOp[] memory innerOpArr = new InnerOp[](1);

        innerOpArr[0] = InnerOp({
            valid: true,
            hash: HashOp(existenceProofInnerOpHash),
            prefix: _innerOpPrefix,
            suffix: _innerOpSuffix
        });
        ExistenceProof memory proof = ExistenceProof({
            valid: true,
            key: key,
            value: value,
            leaf: leafOp,
            path: innerOpArr
        });

        return proof;
    }

    function verifyProof(ExistenceProof memory proof, bytes memory key)
        public
        view
        returns (bool)
    {
        // Verify membership
        verify(proof, getIavlSpec(), relayNetworkHash, key, proof.value);

        return true;
    }
}
