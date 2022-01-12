// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Pausable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "./ancon/IAnconProtocol.sol";
import "./ancon/TrustedOffchainHelper.sol";
import "./ics23/Ics23Helper.sol";

//  a NFT secure document
contract XDVNFT is
    ERC721Burnable,
    ERC721Pausable,
    ERC721URIStorage,
    Ownable,
    TrustedOffchainHelper
{
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;
    IERC20 public stablecoin;
    IAnconProtocol public anconprotocol;
    address public dagContractOperator;
    uint256 public serviceFeeForPaymentAddress = 0;
    uint256 public serviceFeeForContract = 0;

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
        address anconprotocolAddr
    ) ERC721(name, symbol) {
        stablecoin = IERC20(tokenERC20);
        anconprotocol = IAnconProtocol(anconprotocolAddr);
    }

    function setServiceFeeForPaymentAddress(uint256 _fee) public onlyOwner {
        serviceFeeForPaymentAddress = _fee;
    }

    function setServiceFeeForContract(uint256 _fee) public onlyOwner {
        serviceFeeForContract = _fee;
    }

    function mint(address toAddress, uint256 tokenId)
        external
        returns (bytes32)
    {
        revert UsageInformation(
            "Requires anconprotocol proof to execute minting. See https://github.com/anconprotocol for more info"
        );
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
        require(
            anconprotocol.submitPacketWithProof(
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
        _tokenIds.increment();
        uint256 newItemId = _tokenIds.current();
        _safeMint(user, newItemId);
        _setTokenURI(newItemId, uri);

        return newItemId;
    }

    /**
     * @dev Burns a XDV Data Token
     */
    function burnWithProof(
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) public returns (uint256) {
        require(
            anconprotocol.submitPacketWithProof(
                msg.sender,
                userProof,
                key,
                packet,
                proof
            ),
            "invalid packet proof"
        );
        uint256 id = abi.decode(packet, (uint256));
        require(hash == keccak256(abi.encodePacked(id)), "Invalid packet");
        _burn(id);
        return id;
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
     * @dev Locks a XDV Data Token
     */
    function lockWithProof(
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) public returns (uint256) {
        // require(
        //     anconprotocol.submitPacketWithProof(
        //         msg.sender,
        //         userProof,
        //         key,
        //         packet,
        //         proof
        //     ),
        //     "invalid packet proof"
        // );
        // uint256 id = abi.decode(packet, (uint256));
        // require(hash == keccak256(abi.encodePacked(id)), "Invalid packet");
        // _lock(id);
        // return id;
    }

    function _lock(uint256 tokenId)
        internal
        override(ERC721, ERC721URIStorage)
    {
        // super._burn(tokenId);
    }

    /**
     * @dev Locks a XDV Data Token
     */
    function unlockWithProof(
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) public returns (uint256) {
        // require(
        //     anconprotocol.submitPacketWithProof(
        //         msg.sender,
        //         userProof,
        //         key,
        //         packet,
        //         proof
        //     ),
        //     "invalid packet proof"
        // );
        // uint256 id = abi.decode(packet, (uint256));
        // require(hash == keccak256(abi.encodePacked(id)), "Invalid packet");
        // _unlock(id);
        // return id;
    }

    function _unlock(uint256 tokenId)
        internal
        override(ERC721, ERC721URIStorage)
    {
        // super._burn(tokenId);
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

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId
    ) internal virtual override(ERC721, ERC721Pausable) {
        require(!paused(), "XDV: Token execution is paused");

        if (from == address(0)) {
            paymentBeforeMint(msg.sender);
        }

        super._beforeTokenTransfer(from, to, tokenId);
    }

    /**
     * @dev tries to execute the payment when the token is minted.
     * Reverts if the payment procedure could not be completed.
     */
    function paymentBeforeMint(address tokenHolder) internal virtual {
        // Transfer tokens to pay service fee
        require(
            stablecoin.transferFrom(
                tokenHolder,
                address(this),
                serviceFeeForContract
            ),
            "XDV: Transfer failed for recipient"
        );

        emit ServiceFeePaid(
            tokenHolder,
            serviceFeeForContract,
            serviceFeeForPaymentAddress
        );
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
