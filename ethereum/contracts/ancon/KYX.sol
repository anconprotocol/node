pragma solidity ^0.8.7;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./IAnconProtocol.sol";
import "./TrustedOffchainHelper.sol";
import "../ics23/Ics23Helper.sol";

contract KYX is Ownable {
    struct Issuer {
        uint256 id;
        bytes32 category;
        string metadata;
        uint256 reputation;
        bool enabled;
        address creator;
    }
    event IssuerAdded(
        uint256 indexed id,
        bytes32 indexed category,
        string metadata
    );
    event Withdrawn(address indexed payee, uint256 weiAmount);

    mapping(bytes32 => uint256) public issuersCount;
    mapping(bytes32 => mapping(uint256 => Issuer)) public issuers;
    uint256 public fee;
    IERC20 public stablecoin;
    IAnconProtocol public anconprotocol;
    uint256 chainId = 0;

    constructor(
        address tokenERC20,
        address ancon,
        uint256 chain
    ) public {
        stablecoin = IERC20(tokenERC20);
        anconprotocol = IAnconProtocol(ancon);
        chainId = chain;
    }

    // withdraws gas token, must be admin
    function withdraw(address payable payee) public onlyOwner {
        uint256 b = address(this).balance;
        (bool sent, bytes memory data) = payee.call{value: b}("");
        require(sent, "Failed to send Ether");

        emit Withdrawn(payee, b);
    }

    // withdraws protocol fee token, must be admin
    function withdrawToken(address payable payee, address erc20token)
        public
        onlyOwner
    {
        uint256 balance = IERC20(erc20token).balanceOf(address(this));

        // Transfer tokens to pay service fee
        require(IERC20(erc20token).transfer(payee, balance), "transfer failed");

        emit Withdrawn(payee, balance);
    }

    function setFee(uint256 _fee) public onlyOwner {
        fee = _fee;
    }

    function getFee() public returns (uint256) {
        return fee;
    }

    // Implementation

    // Returns a count of issuers by category
    function getIssuerLength(bytes32 category) public returns (uint256) {
        require(issuersCount[category] > 0, "no issuers found");
        return issuersCount[category];
    }

    // Returns an issuer
    function getIssuer(bytes32 category, uint256 id)
        public
        returns (Issuer memory)
    {
        require(issuers[category][id].enabled == false, "no issuers found");
        return issuers[category][id];
    }

    // Enrolls a relayer by moniker name, returns an id
    function enrollIssuerWithProof(
        bytes32 moniker,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory packetProof
    ) public returns (uint256) {
        require(
            keccak256(anconprotocol.getProtocolHeader(moniker)) != keccak256(""),
            "Invalid moniker"
        );
        require(
            anconprotocol.submitPacketWithProof(
                moniker,
                msg.sender,
                userProof,
                packetProof.key,
                packet,
                packetProof
            ),
            "invalid packet proof"
        );
        (bytes32 category, uint256 id, string memory uri) = abi.decode(
            packet,
            (bytes32, uint256, string)
        );
        require(
            issuers[category][id].enabled == true,
            "issuer already exists and enabled"
        );

        issuersCount[category] = issuersCount[category] + 1;
        issuers[category][id] = Issuer({
            id: id,
            category: category,
            metadata: uri,
            enabled: true,
            creator: msg.sender,
            reputation: 0
        });
        emit IssuerAdded(id, category, uri);
        return id;
    }

    // Sets new offchain verifiable data reference for issuer and
    // if enabled or disabled
    function setIssuerWithProof(
        bytes32 category,
        uint256 issuerID,
        string memory metadataUri
    ) public {
        // require only creator can set issuer with proof
    }

     // Adds rating to an issuer, must post proof as evidence
    function setIssuerRatingWithProof(
        bytes32 category,
        uint256 issuerID,
        string memory metadataUri
    ) public {
        // add rating logic
        // add threshold that disables an issuer
    }

    
}
