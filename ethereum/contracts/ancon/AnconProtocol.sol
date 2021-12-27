// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";

contract AnconProtocol is ICS23 {
    address public owner;
    bytes public relayNetworkHash;
    
    mapping(string => ExistenceProof) public accountProofs;
    mapping(address => ExistenceProof) public accountByAddrProofs;
    mapping(bytes => ExistenceProof) public proofs;

    event ProofPacketSubmitted(bytes key, bytes packet);

    constructor(address _onlyOwner) public {
        owner = _onlyOwner;
    }

    function enrollL2Account(string memory did, ExistenceProof memory proof)
        public
        payable
        returns (bool)
    {
        accountProofs[did] = proof;
        accountByAddrProofs[msg.sender] = proof;
        return true;
    }

    function updateProtocolHeader(bytes memory rootHash)
        public
        returns (bool)
    {
        require(msg.sender == owner);
        relayNetworkHash = rootHash;
        return true;
    }

    function submitPacketWithProof(
        ExistenceProof memory proof,
        bytes memory packet,
        bytes memory key
    ) public payable returns (bool) {
        // 1. Verify
        require(
            keccak256(proof.value) == keccak256(packet),
            "bad packet: packet hash is different from ics23 value"
        );
        verify(proof, getIavlSpec(), relayNetworkHash, key, proof.value);

        proofs[key] = proof;

        // 2. Submit event
        emit ProofPacketSubmitted(key, packet);

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
