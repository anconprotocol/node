pragma solidity ^0.8.7;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./IAnconProtocol.sol";
import "./TrustedOffchainHelper.sol";
import "../ics23/Ics23Helper.sol";


// InstantRelayer is used by end users to pay for expedited block root hash commit
contract InstantRelayer is Ownable {
    struct Ticket {
        uint256 id;
        string destination;
        bool open;
    }
    event Withdrawn(address indexed payee, uint256 weiAmount);
    event InstantBlockPaid(
        bytes32 indexed moniker,
        address indexed from,
        uint256 id
    );
    event InstantBlockApplied(bytes32 indexed moniker, address indexed from,  string destination);
    mapping(bytes32 => mapping(address => Ticket)) public tickets;
    mapping(bytes32 => string[]) public relayers;
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

    // Returns a list of relayers uris
    function getRelayers(bytes32 moniker) public returns (string[] memory) {
        require(relayers[moniker].length > 0, "no relayer uris found");
        return relayers[moniker];
    }

    // Enrolls a relayer by moniker name, returns an id
    function enrollRelayer(bytes32 moniker, string memory uri)
        public
        returns (uint256)
    {
        require(
            relayers[moniker].length < 11,
            "maximum registered relayers is 10"
        );
        require(
            keccak256(anconprotocol.getProtocolHeader(moniker)) != keccak256(""),
            "Invalid moniker"
        );

        relayers[moniker].push(uri);
        return relayers[moniker].length;
    }

    // Pays for expedited block, must send a uuid as identifier. Returns ticket id.
    function payForExpediteBlock(bytes32 moniker, uint256 uuid, string memory destination)
        public
        payable
        returns (uint256)
    {
        require(relayers[moniker].length > 0, "moniker has no uri registered");

        require(
            stablecoin.balanceOf(address(msg.sender)) > fee,
            "no enough balance"
        );
        require(
            stablecoin.transferFrom(msg.sender, address(this), fee),
            "transfer failed for recipient"
        );
        tickets[moniker][msg.sender] = Ticket({id: uuid, open: true, destination: destination });
        emit InstantBlockPaid(
            moniker,
            msg.sender,
            tickets[moniker][msg.sender].id
        );
        return tickets[moniker][msg.sender].id;
    }

    // call by agent listening InstantBlockPaid, subscriber must match destination
    function applyBlock(
        bytes32 moniker,
        address sender,
        uint256 ticket
    ) public {
        require(
            tickets[moniker][sender].id == ticket &&
                tickets[moniker][sender].open == true,
            "invalid ticket or dag"
        );

        tickets[moniker][sender].open = false;
        // WIP: Add Rewards for mining
        emit InstantBlockApplied(moniker, sender,tickets[moniker][sender].destination);
    }
}
