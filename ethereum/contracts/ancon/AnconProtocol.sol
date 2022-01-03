// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract AnconProtocol is ICS23, Ownable {
    bytes32 public ENROLL_PAYMENT = keccak256("ENROLL_PAYMENT");
    bytes32 public SUBMIT_PAYMENT = keccak256("SUBMIT_PAYMENT");

    address public relayer;
    bytes public relayNetworkHash;

    IERC20 public stablecoin;
    uint256 public protocolFee = 0;
    uint256 public accountRegistrationFee = 0;

    mapping(bytes => bytes) public accountProofs;
    mapping(address => bytes) public accountByAddrProofs;
    mapping(bytes => bool) public proofs;

    event Withdrawn(address indexed paymentAddress, uint256 amount);

    event ServiceFeePaid(address indexed from, uint256 fee);

    event HeaderUpdated(bytes hash);
    event ProofPacketSubmitted(bytes key, bytes packet);
    event AccountRegistered(bool enrolledStatus, bytes key, bytes value);

    constructor(address _relayer, address tokenERC20) public {
        relayer = _relayer;
        stablecoin = IERC20(tokenERC20);
    }

    function withdrawBalance(address payable payee) public onlyOwner {
        uint256 balance = stablecoin.balanceOf(address(this));

        require(stablecoin.transfer(payee, balance), "transfer failed");

        emit Withdrawn(payee, balance);
    }

    function protocolPayment(bytes32 paymentType, address tokenHolder)
        internal
        virtual
    {
        uint256 fee = 0;
        if ((paymentType) == ENROLL_PAYMENT) {
            fee = accountRegistrationFee;
        }
        if ((paymentType) == SUBMIT_PAYMENT) {
            fee = protocolFee;
        } // Transfer tokens to pay service fee
        require(
            stablecoin.transferFrom(tokenHolder, address(this), protocolFee),
            "transfer failed for recipient"
        );

        emit ServiceFeePaid(tokenHolder, protocolFee);
    }

    function setProtocolFee(uint256 _fee) public onlyOwner {
        protocolFee = _fee;
    }

    function setAccountRegistrationFee(uint256 _fee) public onlyOwner {
        accountRegistrationFee = _fee;
    }

    function getProtocolHeader() public view returns (bytes memory) {
        return relayNetworkHash;
    }

    function getProof(bytes memory did) public view returns (bytes memory) {
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
        require(keccak256(proof.key) == keccadk256(key), "invalid key");

        require(verifyProof(proof), "invalid proof");

        require(
            keccak256(key) == keccak256(accountProofs[did]),
            "user already registered"
        );

        protocolPayment(ENROLL_PAYMENT, msg.sender);

        accountProofs[(did)] = key;
        accountByAddrProofs[msg.sender] = key;

        emit AccountRegistered(true, key, did);
        return true;
    }

    function updateProtocolHeader(bytes memory rootHash) public {
        require(msg.sender == relayer);
        relayNetworkHash = rootHash;
        emit HeaderUpdated(rootHash);
    }

    function submitPacketWithProof(
        bytes memory key,
        bytes memory packet,
        ExistenceProof memory proof
    ) public payable returns (bool) {
        // 1. Verify
        require(keccak256(proof.key) == keccadk256(key), "invalid key");

        require(verifyProof(proof));

        proofs[key] = true;

        protocolPayment(SUBMIT_PAYMENT, msg.sender);

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
        internal
        pure
        returns (bytes memory)
    {
        return bytes(calculate(proof));
    }
}
