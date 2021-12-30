// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";

contract AnconProtocol is ICS23 {
    address public owner;
    address public relayer;
    bytes public relayNetworkHash;

    mapping(bytes => bytes) public accountProofs;
    mapping(address => bytes) public accountByAddrProofs;
    mapping(bytes => bool) public proofs;

    event HeaderUpdated(bytes hash);
    event ProofPacketSubmitted(bytes key, bytes packet);
    event AccountRegistered(bool enrolledStatus, bytes key, bytes value);

    constructor(address _onlyOwner, address _relayer) public {
        owner = _onlyOwner;
        relayer = _relayer;
    }

    function getProtocolHeader() public view returns (bytes memory) {
        return relayNetworkHash;
    }

    function getProof(bytes memory  did) public view returns (bytes memory) {
        return accountProofs[did];
    }

    function hasProof(bytes memory key) public view returns (bool) {
        return proofs[key];
    }

    function enrollL2Account(
        bytes memory key, // proof key "/anconprotocol/root/user/diddocid"
        bytes memory did, // proof value did doc id
        ExistenceProof memory proof
    ) public payable returns (bool) {
        require(verifyProof(proof));
        accountProofs[(did)] = key;
        accountByAddrProofs[msg.sender] = key;
        emit AccountRegistered(true, key, did);
        return true;
    }

    function updateProtocolHeader(bytes memory rootHash) public returns (bool) {
        require(msg.sender == relayer);
        relayNetworkHash = rootHash;
        emit HeaderUpdated(rootHash);
        return true;
    }

    function submitPacketWithProof(
        bytes memory key,
        bytes memory packet,
        ExistenceProof memory proof
    ) public payable returns (bool) {
        // 1. Verify
        require(verifyProof(proof));

        proofs[key] = true;

        // 2. Submit event
        emit ProofPacketSubmitted(key, packet);

        return true;
    }

    function verifyProof(ExistenceProof memory exProof)
        public
        view
        returns (bool)
    {
        // Verify membership
        verify(
            exProof,
            getIavlSpec(),
            relayNetworkHash,
            exProof.key,
            exProof.value
        );

        return true;
    }

    function queryRootCalculation(ExistenceProof memory proof)
        public
        pure
        returns (bytes memory)
    {
        return bytes(calculate(proof));
    }
}
