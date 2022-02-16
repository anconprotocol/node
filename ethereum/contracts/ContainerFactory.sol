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
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";

//  WIP: ContainerFactory creates NFT Containers for UI no-code users
contract ContainerFactory is
    ERC721Burnable,
    ERC721Pausable,
    IERC721Receiver,
    ERC721URIStorage,
    Ownable,
    TrustedOffchainHelper
{
    using Counters for Counters.Counter;
    bytes32 public TOKEN_LOCKED = keccak256("TOKEN_LOCKED");
    bytes32 public TOKEN_BURNED = keccak256("TOKEN_BURNED");
    bytes32 public TOKEN_AVAILABLE = keccak256("TOKEN_AVAILABLE");
    bytes32 public ENROLL_NFT = keccak256("ENROLL_NFT");

    Counters.Counter private _tokenIds;
    IERC20 public stablecoin;
    IAnconProtocol public anconprotocol;
    address public dagContractOperator;
    uint256 public NFTRegistrationFee = 0;
    uint256 public serviceFeeForPaymentAddress = 0;
    uint256 public serviceFeeForContract = 0;
    mapping(address => mapping(uint256 => bytes32)) public tokenLockStorage;
    uint256 chainId = 0;
    mapping(address => bool) nftRegistry;

    event Withdrawn(address indexed paymentAddress, uint256 amount);
    event Locked(address nftContractAddress, uint256 indexed id);
    event Released(address sender, uint256 indexed id);
    event ServiceFeePaid(address indexed from, uint256 fee);

    event NFTEnrolled(bool enrolledStatus, address NFTaddress);

    /**
     * WXDV
     */
    constructor(
        string memory name,
        string memory symbol,
        address tokenERC20,
        address anconprotocolAddr,
        uint256 chain
    ) ERC721(name, symbol) {
        stablecoin = IERC20(tokenERC20);
        anconprotocol = IAnconProtocol(anconprotocolAddr);
        chainId = chain;
    }

    function onERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes calldata data
    ) external returns (bytes4) {
        return this.onERC721Received.selector;
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
    function submitMintWithProof(
        address sender,
        uint256 newItemId,
        bytes32 moniker,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) external payable returns (bool) {
        require(
            anconprotocol.submitPacketWithProof(
                moniker,
                sender,
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
        //Newly minted NFTs are not locked
        tokenLockStorage[sender][newItemId] = TOKEN_AVAILABLE;
        return true;
    }

    // protocolPayment handles contract payment protocol fee types
    function protocolPayment(bytes32 paymentType, address tokenHolder)
        internal
    {
        require(
            stablecoin.balanceOf(address(msg.sender)) > 0,
            "no enough balance"
        );
        if ((paymentType) == ENROLL_NFT) {
            require(
                stablecoin.transferFrom(
                    tokenHolder,
                    address(this),
                    NFTRegistrationFee
                ),
                "transfer failed for recipient"
            );
            emit ServiceFeePaid(tokenHolder, NFTRegistrationFee);
        }
    }

    function enrollNFT(address NFTaddress) public payable returns (bool) {
        require(nftRegistry[NFTaddress] == false, "NFT is already in registry");

        protocolPayment(ENROLL_NFT, msg.sender);

        nftRegistry[NFTaddress] = true;

        emit NFTEnrolled(true, NFTaddress);
        return true;
    }

    function deactivateNFT(address NFTaddress) public onlyOwner returns (bool) {
        require(nftRegistry[NFTaddress] == true, "missing nft address");

        nftRegistry[NFTaddress] = false;

        emit NFTEnrolled(false, NFTaddress);
        return true;
    }

    /**
     * @dev Burns a XDV Data Token
     */
    // function burnWithProof(
    //     bytes memory key,
    //     bytes memory packet,
    //     Ics23Helper.ExistenceProof memory userProof,
    //     Ics23Helper.ExistenceProof memory proof,
    //     bytes32 hash
    // ) public returns (uint256) {
    //     require(
    //         anconprotocol.submitPacketWithProof(
    //             moniker,
    //             msg.sender,
    //             userProof,
    //             key,
    //             packet,
    //             proof
    //         ),
    //         "invalid packet proof"
    //     );
    //     uint256 id = abi.decode(packet, (uint256));
    //     require(hash == keccak256(abi.encodePacked(id)), "Invalid packet");
    //     _burn(id);
    //     tokenLockStorage[id] = TOKEN_BURNED;
    //     return id;
    // }

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
        address sender,
        bytes32 moniker,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) external payable returns (uint256) {
        require(nftRegistry[msg.sender] == true, "nft must be registered");
        require(
            anconprotocol.submitPacketWithProof(
                moniker,
                sender,
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

        ERC721 nftContractCaller = ERC721(msg.sender);
        require(sender == nftContractCaller.ownerOf(id), "invalid owner");

        if (
            tokenLockStorage[sender][id] == TOKEN_AVAILABLE //nftContractCaller.ownerOf(id) == sender
        ) {
            // require(
            //     nftContractCaller.getApproved(id) == address(this),
            //     "WXDV needs to be approved for lock operation"
            // );
            // nftContractCaller.safeTransferFrom(sender, address(this), id);
            // // Set as locked
            lock(sender, id);
            emit Locked(sender, id);
        } else {
            // todo: must be sender
            require(ownerOf(id) == address(this), "Is not a wrapped token");
            _burn(id); // todo: must give an incentive to holder
        }
        return id;
    }

    function lock(address holder, uint256 tokenId) internal {
        require(
            tokenLockStorage[holder][tokenId] == TOKEN_AVAILABLE,
            "Token is already locked"
        );
        tokenLockStorage[holder][tokenId] = TOKEN_LOCKED;
    }

    /**
     * More info at https://github.com/renproject/ren/wiki#cross-chain-transactions
     * @dev Releases a XDV Data Token
     */
    function releaseWithProof(
        address sender,
        bytes32 moniker,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) external payable returns (uint256) {
        require(nftRegistry[msg.sender] == true, "nft must be registered");
        require(
            anconprotocol.submitPacketWithProof(
                moniker,
                sender,
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


        ERC721 nftContractCaller = ERC721(msg.sender);
        require(sender == nftContractCaller.ownerOf(id), "invalid owner");

        if (
            tokenLockStorage[sender][id] == TOKEN_LOCKED //nftContractCaller.ownerOf(id) == sender
        ) {
            // require(
            //     nftContractCaller.getApproved(id) == address(this),
            //     "WXDV needs to be approved for lock operation"
            // );
            // nftContractCaller.safeTransferFrom(sender, address(this), id);
            // // Set as locked
            unlock(sender, id);
            emit Released(sender, id);
            return 2;
        } else {
            // todo: must be sender

            // if token doesnt exist I need to create a wrapped xdv token
            _tokenIds.increment();
            uint256 newItemId = _tokenIds.current();
            _safeMint(newOwner, newItemId);
            _setTokenURI(id, metadataUri);

            return  1;
        }
        return 0;
    }

    function unlock(address sender, uint256 tokenId) internal {
        require(
            tokenLockStorage[sender][tokenId] == TOKEN_LOCKED,
            "Token is already unlocked"
        );
        tokenLockStorage[sender][tokenId] = TOKEN_AVAILABLE;
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

        emit ServiceFeePaid(tokenHolder, serviceFeeForContract);
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

    // add two function modifiers
    // a modifier to check the owner of the contract
    // a modifier to determine if the locked flag is true of false
    // modifier onlyOwner() {
    //     require(msg.sender == owner, "Not Owner");
    //     _;
    // }

    // Add a setter to change the locked flag
    // only the owner of the contract can call because a modifier is specified
    function islocked(uint256 tokenId, address sender) public returns (bool) {
        return tokenLockStorage[sender][tokenId] == TOKEN_LOCKED;
    }

    // add the function modifier to the transfer function
    // if the locked==false then one can not trade
    // function transfer(address _to, uint256 _value)
    //     public
    //     returns (bool success)
    // {
    //     require(islocked() == false, "Token is locked");
    //     if (_value > 0 && _value <= balanceOf(msg.sender)) {
    //         __balanceOf[msg.sender] -= _value;
    //         __balanceOf[_to] += _value;
    //         emit Transfer(msg.sender, _to, _value);
    //         return true;
    //     }
    //     return false;
    // }
}
