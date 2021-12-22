// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

contract OnchainMetadata {
    event AddOnchainMetadata(
        string name,
        string description,
        string image,
        string owner,
        string parent,
        bytes sources
    );

    event EncodeDagJson(string path, string hexdata);

    event EncodeDagCbor(string path, string hexdata);

    constructor() {}

// register L2 ACCOUNT
    function setOnchainMetadata(
        string memory proofKey,
        string memory proofValue,
        string memory metadataUri
    ) public {}

    function registerL2Account(
        string memory didAddress,
        bytes memory key,
        bytes memory value,
        bytes memory _prefix,
        uint256[] memory _leafOpUint,
        bytes memory _innerOpPrefix,
        bytes memory _innerOpSuffix,
        uint256 existenceProofInnerOpHash
    ) public payable returns (bool) {}


    function sum(uint256 x, uint256 y) public pure returns (uint256) {
        return 0;
    }

    function encodeDagjsonBlock(string memory path, string memory hexdata)
        public
        returns (bool)
    {
        emit EncodeDagJson(path, hexdata);

        return true;
    }
    //emit AddOnchainMetadata()
}
