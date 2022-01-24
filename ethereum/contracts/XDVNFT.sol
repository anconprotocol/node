// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "./ancon/IAnconProtocol.sol";
import "./ancon/TrustedOffchainHelper.sol";

import "./ics23/Ics23Helper.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";

//  a NFT secure document
contract XDVNFT is
    ERC721Burnable,
    ERC721URIStorage,
    IERC721Receiver,
    Ownable,
    TrustedOffchainHelper
{
    using Counters for Counters.Counter;

    Counters.Counter private _tokenIds;
    IERC20 public stablecoin;
    IAnconProtocol public anconprotocol;

    bytes32 moniker = keccak256("anconprotocol");
    uint256 chainId = 0;

    event Withdrawn(address indexed paymentAddress, uint256 amount);
    event ServiceFeePaid(
        address indexed from,
        uint256 paidToContract,
        uint256 paidToPaymentAddress
    );

    /**
     * XDVNFT Data Token
     */
    constructor(
        string memory name,
        string memory symbol,
        address tokenERC20,
        address ancon,
        uint256 chain
    ) ERC721(name, symbol) {
        stablecoin = IERC20(tokenERC20);
        anconprotocol = IAnconProtocol(ancon);
        chainId = chain;
    }

    function mint(address toAddress, uint256 tokenId)
        external
        returns (bytes32)
    {
        revert UsageInformation(
            "Requires anconprotocol proof to execute minting. See https://github.com/anconprotocol for more info"
        );
    }

    function onERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes calldata data
    ) external returns (bytes4) {
        return this.onERC721Received.selector;
    }

    /**
     * @dev Mints a XDV Data Token
     */
    function mintWithProof(
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) public returns (uint256) {
        _tokenIds.increment();
        uint256 newItemId = _tokenIds.current();
        require(
            anconprotocol.submitPacketWithProof(
                moniker,
                msg.sender,
                userProof,
                key,
                packet,
                proof
            ),
            "invalid packet proof"
        );
        (address user, string memory uri) = abi.decode(
            packet,
            (address, string)
        );
        require(
            hash == keccak256(abi.encodePacked(user, uri)),
            "Invalid packet"
        );

        _safeMint(user, newItemId);
        _setTokenURI(newItemId, uri);
        //Newly minted NFTs are not locked
        return newItemId;
    }

    /**
     * @dev Just overrides the superclass' function. Fixes inheritance
     * source: https://forum.openzeppelin.com/t/how-do-inherit-from-erc721-erc721enumerable-and-erc721uristorage-in-v4-of-openzeppelin-contracts/6656/4
     */
    function _burn(uint256 tokenId)
        internal
        override(ERC721, ERC721URIStorage)
    {
        super._burn(tokenId);
    }

    /**
     * More info at https://github.com/renproject/ren/wiki#cross-chain-transactions
     * @dev Locks a XDV Data Token
     */
    function lockWithProof(
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) public payable returns (uint256) {
        require(
            anconprotocol.submitPacketWithProof(
                moniker,
                msg.sender,
                userProof,
                key,
                packet,
                proof
            ),
            "invalid packet proof"
        );
        (uint256 id, bytes32 contractIdentifier) = abi.decode(
            packet,
            (uint256, bytes32)
        );
        require(
            hash == keccak256(abi.encodePacked(id, contractIdentifier)),
            "invalid packet"
        );

        require(msg.sender == this.ownerOf(id), "invalid owner");
        approve(address(this), id);
        safeTransferFrom(msg.sender, address(this), id);
        return id;
    }

    /**
     * More info at https://github.com/renproject/ren/wiki#cross-chain-transactions
     * @dev Releases a XDV Data Token
     */
    function releaseWithProof(
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) public payable returns (uint256) {
        require(
            anconprotocol.submitPacketWithProof(
                moniker,
                msg.sender,
                userProof,
                key,
                packet,
                proof
            ),
            "invalid packet proof"
        );
        (
            uint256 id,
            string memory metadataUri,
            address newOwner,
            //            bytes32 lockTransactionHash,
            bytes32 contractIdentifier
        ) = abi.decode(packet, (uint256, string, address, bytes32));

        require(
            hash ==
                keccak256(
                    abi.encodePacked(
                        id,
                        metadataUri,
                        newOwner,
                        //                        lockTransactionHash,
                        contractIdentifier
                    )
                ),
            "invalid packet"
        );
        require(
            msg.sender == newOwner && msg.sender != this.ownerOf(id),
            "invalid owner"
        );
        // approve(address(this), id);
        safeTransferFrom(address(this), newOwner, id);
        _setTokenURI(id, metadataUri);
        return id;
    }

    /**
     * @dev Just overrides the superclass' function. Fixes inheritance
     * source: https://forum.openzeppelin.com/t/how-do-inherit-from-erc721-erc721enumerable-and-erc721uristorage-in-v4-of-openzeppelin-contracts/6656/4
     */
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    function withdrawBalance(address payable payee) public onlyOwner {
        uint256 balance = stablecoin.balanceOf(address(this));

        require(stablecoin.transfer(payee, balance), "XDV: Transfer failed");

        emit Withdrawn(payee, balance);
    }

    function withdraw(address payable payee) public onlyOwner {
        uint256 b = address(this).balance;
        (bool sent, bytes memory data) = payee.call{value: b}("");
        require(sent, "Failed to send Ether");

        emit Withdrawn(payee, b);
    }
}
