// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";

contract AnconProtocol is ICS23 {
    address public owner;
    bytes public relayNetworkHash;

    mapping(string => bytes) public accountProofs;
    mapping(address => bytes) public accountByAddrProofs;
    mapping(bytes => bool) public proofs;

    event ProofPacketSubmitted(bytes key, bytes packet);

    constructor(address _onlyOwner) public {
        owner = _onlyOwner;
    }

    function enrollL2Account(string memory did, ExistenceProof memory proof)
        public
        payable
        returns (bool)
    {
        // require(verifyProof(proof, proof.key));
        accountProofs[did] = proof.key;
        accountByAddrProofs[msg.sender] = proof.key;
        return true;
    }

    function updateProtocolHeader(bytes memory rootHash) public returns (bool) {
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
                // require(verifyProof(proof, proof.key));

        proofs[key] = true;

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

    function verifyProof(
        bytes memory key,
        bytes memory value,
        bytes memory _prefix,
        uint256[] memory _leafOpUint,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix,
        uint256 existenceProofInnerOpHash
    )
        public
        view
        returns (bool)
    {

        ExistenceProof memory exProof = convertProof(
            key, 
            value, 
            _prefix, 
            _leafOpUint, 
            _innerOpPrefix,
            _innerOpSuffix, 
            existenceProofInnerOpHash
        );

        // Verify membership
        verify(exProof, getIavlSpec(), relayNetworkHash, key, exProof.value);

        return true;
    }
    function queryRootCalculation(
        uint256[] memory leafOpUint,
        bytes memory prefix,
        // bytes[][] memory existenceProofInnerOp,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix,
        uint256 existenceProofInnerOpHash,
        bytes memory existenceProofKey,
        bytes memory existenceProofValue
    ) public view returns (bytes memory) {
        ExistenceProof memory proof = convertProof(
            existenceProofKey,
            existenceProofValue,
            prefix,
            leafOpUint,
            _innerOpPrefix,
            _innerOpSuffix,
            existenceProofInnerOpHash
        );
        return bytes(calculate(proof));
    }

}
