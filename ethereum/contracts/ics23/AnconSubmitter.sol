// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "./AnconVerifier.sol";

contract AnconSubmitter {
    address public owner;
    AnconVerifier public verifier;
    bytes public relayNetworkHash;
    mapping(bytes => bytes) public proofs;

    event ProofPacketSubmitted();

    constructor(address _onlyOwner, address _verifier) public {
        owner = _onlyOwner;
        verifier = AnconVerifier(verifier);
    }

    function updateRelayHeader(bytes memory rootHash)
        public
        returns (bool)
    {
        require(msg.sender == owner);
        relayNetworkHash = rootHash;
        return true;
    }

    function submitPacketWithProof(
        // -- existence proof payload
        uint256[] memory leafOpUint,
        bytes memory prefix,
        bytes[][] memory existenceProofInnerOp,
        uint256 existenceProofInnerOpHash,
        bytes memory key,
        bytes memory value,
        bytes memory packet
    ) public payable returns (bool) {
        // 1. Verify
        require(
            bytes32(value) == keccak256(packet),
            "bad packet: packet hash is different from ics23 value"
        );
        bytes memory calculatedHash = verifier.queryRootCalculation(
            leafOpUint,
            prefix,
            existenceProofInnerOp,
            existenceProofInnerOpHash,
            key,
            value
        );
        require(
            keccak256(relayNetworkHash) == keccak256(calculatedHash),
            "invalid proof for key"
        );
        proofs[key] = packet;
        // 2. Submit event
        emit ProofPacketSubmitted();

        return true;
    }
}
