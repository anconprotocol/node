// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "./ics23.sol";

/// @title A title that should describe the contract/interface
/// @author The name of the author
/// @notice Explain to an end user what this does
/// @dev Explain to a developer any extra details
contract AnconVerifier is ICS23 {
    address public owner;

    constructor(address onlyOwner) public {
        owner = onlyOwner;
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

    function verifyProof(
        uint256[] memory leafOpUint 
    ) public view returns (bool) {
        // (
        //     bytes memory prefix,
        //     bytes memory existenceProofInnerOpPrefix,
        //     bytes memory existenceProofInnerOpSuffix,
        //     uint256 existenceProofInnerOpHash
        // ) = abi.decode(
        //         existenceProof,
        //         (bytes, bytes, bytes, uint256)
        //     );
        // // todo: verify not empty
        // ExistenceProof memory proof = convertProof(
        //     key,
        //     value,
        //     prefix,
        //     leafOpUint,
        //     existenceProofInnerOpPrefix,
        //     existenceProofInnerOpSuffix,
        //     existenceProofInnerOpHash
        // );

        // // Verify membership
        // verify(proof, getIavlSpec(), root, key, value);

        return true;
    }
}
