// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.7;
pragma experimental ABIEncoderV2;

import "./MerkleTree.sol";
import "./ics23.sol";

contract ICS23Verifier {

    function hashFromByteSlices(bytes[] memory data)
        public
        pure
        returns (bytes32)
    {
        return MerkleTree.hashFromByteSlices(data);
    }

    // function verifyMembership(
    //     bytes memory proofBz,
        
    //     bytes memory rootBz,
    //     bytes memory pathBz,
    //     bytes memory value,
    //     ExistenceProof memory proof
    // ) public view returns (bool) {
        
    //     return verifyMembership(
    //         iavlSpec,
    //         rootBz,
    //         proof,
    //         pathBz,
    //         value
    //     );
    // }
}
