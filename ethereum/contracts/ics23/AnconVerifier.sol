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
        bytes[][] memory _innerOp,
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
        InnerOp[] memory innerOpArr = new InnerOp[](_innerOp.length);

        for (uint256 i = 0; i < _innerOp.length; i++) {
            bytes[] memory temp = _innerOp[i];
            innerOpArr[i] = InnerOp({
                valid: true,
                hash: HashOp(existenceProofInnerOpHash),
                prefix: temp[0],
                suffix: temp[1]
            });
        }
        ExistenceProof memory proof = ExistenceProof({
            valid: true,
            key: key,
            value: value,
            leaf: leafOp,
            path: innerOpArr
        });

        return proof;
    }

    function queryRootCalculation(
        uint256[] memory leafOpUint,
        bytes memory prefix,
        bytes[][] memory existenceProofInnerOp,
        uint256 existenceProofInnerOpHash,
        bytes memory existenceProofKey,
        bytes memory existenceProofValue
    ) public view returns (bytes memory) {
        ExistenceProof memory proof = convertProof(
            existenceProofKey,
            existenceProofValue,
            prefix,
            leafOpUint,
            existenceProofInnerOp,
            existenceProofInnerOpHash
        );
        return bytes(calculate(proof));
    }

    function verifyProof(
        uint256[] memory leafOpUint,
        bytes memory prefix,
        bytes[][] memory existenceProofInnerOp,
        uint256 existenceProofInnerOpHash,
        bytes memory existenceProofKey,
        bytes memory existenceProofValue,
        bytes memory root,
        bytes memory key,
        bytes memory value
    ) public pure returns (bool) {
        // todo: verify not empty
        ExistenceProof memory proof = convertProof(
            existenceProofKey,
            existenceProofValue,
            prefix,
            leafOpUint,
            existenceProofInnerOp,
            existenceProofInnerOpHash
        );

        // Verify membership
        verify(proof, getIavlSpec(), root, key, value);

        return true;
    }
}