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
        bytes memory key, // proof key "/anconprotocol/root/user/diddocid"
        bytes memory did, // proof value did doc id
        bytes memory _prefix,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix,
        bytes[][] memory _innerOp
    ) public payable returns (bool) {
        require(
            verifyProof(
                key,
                did,
                _prefix,
                _innerOpPrefix,
                _innerOpSuffix,
                _innerOp
            )
        );
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
        bytes memory _innerOpSuffix,
        bytes[][] memory _innerOp
    ) public payable returns (bool) {
        // 1. Verify
        require(
            verifyProof(
                key,
                packet,
                _prefix,
                _innerOpPrefix,
                _innerOpSuffix,
                _innerOp
            )
        );

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
        bytes memory _innerOpSuffix,
        bytes[][] memory _innerOp
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
        InnerOp[] memory innerOpArr = new InnerOp[](_innerOp.length);

        for (uint256 i = 0; i < _innerOp.length; i++) {
            bytes[] memory temp = _innerOp[i];
            innerOpArr[i] = InnerOp({
                valid: true,
                hash: HashOp(1),
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

    function verifyProof(
        bytes memory key,
        bytes memory value,
        bytes memory _prefix,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix,
        bytes[][] memory _innerOp
    ) public view returns (bool) {
        ExistenceProof memory exProof = convertProof(
            key,
            value,
            _prefix,
            _innerOpPrefix,
            _innerOpSuffix,
            _innerOp
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
        bytes memory existenceProofValue,
        bytes[][] memory _innerOp
    ) public view returns (bytes memory) {
        ExistenceProof memory proof = convertProof(
            existenceProofKey,
            existenceProofValue,
            prefix,
            _innerOpPrefix,
            _innerOpSuffix,
            _innerOp
        );
        return bytes(calculate(proof));
    }
}
