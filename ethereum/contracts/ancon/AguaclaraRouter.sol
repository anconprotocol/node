// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "./AnconVerifier.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Pausable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";

 contract AguaclaraRouter {
    address public owner;
    AnconVerifier public verifier;
    mapping(bytes32 => bytes) public headers;

    event MetadataOwnershipChanged(
        address receiver,
        string uri,
        uint256 tokenID
    );

    constructor(address _onlyOwner, address _verifier) public {
        owner = _onlyOwner;
        verifier = AnconVerifier(verifier);
    }

    function updateHeader(bytes32 key, bytes memory rootHash)
        public
        returns (bool)
    {
        headers[key] = rootHash;
        return true;
    }

    function verifyOwnershipReference(
           // -- existence proof
        uint256[] memory leafOpUint,
        bytes memory prefix,
        bytes[][] memory existenceProofInnerOp,
        uint256 existenceProofInnerOpHash,
        bytes memory existenceProofKey,
        bytes memory existenceProofValue,
        bytes memory key,
        bytes memory value,
        bytes memory packetMetadataUri
    ) public view returns (bool) {
        // 1. Verify
        bytes memory calculatedHash = verifier.queryRootCalculation(
            leafOpUint,
            prefix,
            existenceProofInnerOp,
            existenceProofInnerOpHash,
            existenceProofKey,
            existenceProofValue
        );
        // // https://github.com/smartcontractkit/solidity-cborutils
        // // Use https://github.com/chrisdotn/jsmnSol
        // require(
        //     keccak256(value) == keccak256(packetMetadataUri),
        //     "Invalid Proof for key"
        // );

        return true;
    }

    function sendMetadataOwnership(
        // -- existence proof
        uint256[] memory leafOpUint,
        bytes memory prefix,
        bytes[][] memory existenceProofInnerOp,
        uint256 existenceProofInnerOpHash,
        bytes memory existenceProofKey,
        bytes memory existenceProofValue,
        bytes memory key,
        bytes memory value,
        // -- ics23 packet
        // bytes memory packet, (abiDecoder de ethers)
        string memory metadata,
        uint256 tokenId,
        address tokenAddress,
        address to
    ) public returns (bool) {
        // 1. Verify
        bytes memory calculatedHash = verifier.queryRootCalculation(
            leafOpUint,
            prefix,
            existenceProofInnerOp,
            existenceProofInnerOpHash,
            existenceProofKey,
            existenceProofValue
        );
        bytes memory current = headers[keccak256(key)];
        require(
            keccak256(current) == keccak256(calculatedHash),
            "Invalid Proof for key"
        );

        // 2. Instantiate NFT and query
        ERC721 nft = ERC721(tokenAddress);
        require(nft.ownerOf(tokenId) == msg.sender, "Invalid token id");

        // if (isNew) {
        //     // todo: should mint
        // } else {
        nft.safeTransferFrom(
            msg.sender,
            to,
            tokenId,
            abi.encodePacked(metadata, tokenId)
        );
        //        }

// Ics23 
// packet - dag toMetadata => link
// evm verifies ics23
// browser  verifies packet is child of ics23 metadata
        emit MetadataOwnershipChanged(to, metadata, tokenId);

        return true;
    }
}
