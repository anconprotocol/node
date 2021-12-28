// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";

contract AnconProtocol is ICS23 {
    address public owner;
    bytes public relayNetworkHash;

    mapping(bytes => bytes) public accountProofs;
    mapping(address => bytes) public accountByAddrProofs;
    mapping(bytes => bool) public proofs;

    event ProofPacketSubmitted(bytes key, bytes packet);

    constructor(address _onlyOwner) public {
        owner = _onlyOwner;
    }

    function enrollL2Account(
        bytes memory key, // did cid
        bytes memory did, // did id
        bytes memory _prefix,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix
    ) public payable returns (bool) {
        require(verifyProof(key, did, _prefix, _innerOpPrefix, _innerOpSuffix));
        accountProofs[(did)] = key;
        accountByAddrProofs[msg.sender] = key;
        return true;
    }

    function updateProtocolHeader(bytes memory rootHash) public returns (bool) {
        require(msg.sender == owner);
        relayNetworkHash = rootHash;
        return true;
    }

    function submitPacketWithProof(
        bytes memory key,
        bytes memory packet,
        bytes memory _prefix,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix
    ) public payable returns (bool) {
        // 1. Verify
        require(verifyProof(key, packet, _prefix, _innerOpPrefix, _innerOpSuffix));

        proofs[key] = true;

        // 2. Submit event
        emit ProofPacketSubmitted(key, packet);

        return true;
    }

    function convertProof(
        bytes memory key,
        bytes memory value,
        bytes memory _prefix,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix
    ) public pure returns (ExistenceProof memory) {
        LeafOp memory leafOp = LeafOp(
            true,
            HashOp.SHA256,
            HashOp.NO_HASH,
            HashOp.SHA256,
            LengthOp.VAR_PROTO,
            _prefix
        );

        // // innerOpArr
        InnerOp[] memory innerOpArr = new InnerOp[](1);

        innerOpArr[0] = InnerOp({
            valid: true,
            hash: HashOp.SHA256,
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
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix
    )
        public
        view
        returns (
            bool
        )
    {
        ExistenceProof memory exProof = convertProof(
            key,
            value,
            _prefix,
            _innerOpPrefix,
            _innerOpSuffix
        );

        // Verify membership
        verify(exProof, getIavlSpec(), relayNetworkHash, key, exProof.value);

        return true;
    }

    function queryRootCalculation(
        bytes memory prefix,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix,
        bytes memory existenceProofKey,
        bytes memory existenceProofValue
    ) public view returns (bytes memory) {
        ExistenceProof memory proof = convertProof(
            existenceProofKey,
            existenceProofValue,
            prefix,
            _innerOpPrefix,
            _innerOpSuffix
        );
        return bytes(calculate(proof));
    }
}
