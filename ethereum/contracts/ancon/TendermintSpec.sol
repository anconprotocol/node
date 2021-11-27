// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

contract TendermintSpec {


    event CrossMintCallbackReceived(
        address indexed newOwner,
        string indexed metadataHash,
        uint256 indexed tokenId
    );

    constructor(
    ) public {
    }

    // bytes data {
    //     // Onchain metadata
    //     string memory metadata,

    //     // From NFT contract address    
    //     address sourceNft,
    //     uint sourceChainId,
    //     address fromOwner,

    //     // To NFT contract address    
    //     address destinationNft,
    //     uint destinationChainId,
    //     address toOwner,

    //     // Source TokenID
    //     uint256 sourceTokenid,

    //     // ???
    //     bool isNew
    // }
}