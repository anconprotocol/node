// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract AnconProtocol is ICS23 {
    bytes32 public ENROLL_PAYMENT = keccak256("ENROLL_PAYMENT");
    bytes32 public SUBMIT_PAYMENT = keccak256("SUBMIT_PAYMENT");

    address public owner;
    address public relayer;
    bytes public relayNetworkHash;

    IERC20 public stablecoin;
    uint256 public protocolFee = 0;
    uint256 public accountRegistrationFee = 0;

    mapping(bytes => bytes) public accountProofs; //did user-assigned proof key
    mapping(address => bytes) public accountByAddrProofs; //proof key-assigned eth address
    mapping(bytes => bool) public proofs; //if proof key was submitted to the blockchain

    event Withdrawn(address indexed paymentAddress, uint256 amount);

    event ServiceFeePaid(address indexed from, uint256 fee);

    event HeaderUpdated(bytes hash);
    event ProofPacketSubmitted(
        bytes key,
        bytes packet,
        Ics23Helper.ExistenceProof proof
    );
    event AccountRegistered(
        bool enrolledStatus,
        bytes key,
        bytes value,
        Ics23Helper.ExistenceProof proof
    );

    constructor() public {
        owner = msg.sender;
    }

    function setPaymentToken(IERC20 tokenAddress) public {
        require(owner == msg.sender);
        stablecoin = tokenAddress;
    }

    function withdraw(address payable payee) public {
        require(owner == msg.sender);
        uint256 b = address(this).balance;
        (bool sent, bytes memory data) = payee.call{value: b}("");
        require(sent, "Failed to send Ether");

        emit Withdrawn(payee, b);
    }

    function withdrawToken(address payable payee, address erc20token) public {
        require(owner == msg.sender);
        uint256 balance = IERC20(erc20token).balanceOf(address(this));

        // Transfer tokens to pay service fee
        require(IERC20(erc20token).transfer(payee, balance), "transfer failed");

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
            stablecoin.transferFrom(tokenHolder, address(this), fee),
            "transfer failed for recipient"
        );

        emit ServiceFeePaid(tokenHolder, fee);
    }

    function setProtocolFee(uint256 _fee) public {
        require(owner == msg.sender);
        protocolFee = _fee;
    }

    function setAccountRegistrationFee(uint256 _fee) public {
        require(owner == msg.sender);
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
        bytes memory key,
        bytes memory did,
        Ics23Helper.ExistenceProof memory proof
    ) public payable returns (bool) {
        require(keccak256(proof.key) == keccak256(key), "invalid key");

        require(verifyProof(proof), "invalid proof");

        require(
            keccak256(key) != keccak256(accountProofs[did]),
            "user already registered"
        );

        protocolPayment(ENROLL_PAYMENT, msg.sender);

        accountProofs[(did)] = key;
        accountByAddrProofs[msg.sender] = key;

        emit AccountRegistered(true, key, did, proof);
        return true;
    }

    function updateProtocolHeader(bytes memory rootHash) public {
        require(owner == msg.sender);
        // require(msg.sender == relayer);
        relayNetworkHash = rootHash;
        emit HeaderUpdated(rootHash);
    }

    function submitPacketWithProof(
        address sender,
        Ics23Helper.ExistenceProof memory userProof,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory proof
    ) external payable returns (bool) {
        // 1. Verify
        require(proofs[key] == false, "proof has been submitted (found key)");
        require(keccak256(proof.key) == keccak256(key), "invalid key");
        require(
            keccak256(accountByAddrProofs[sender]) == keccak256(userProof.key),
            "invalid user key"
        );
        require(verifyProof(userProof), "invalid user proof");
        require(verifyProof(proof));

        proofs[key] = true;

        protocolPayment(SUBMIT_PAYMENT, sender);

        // 2. Submit event
        emit ProofPacketSubmitted(key, packet, proof);

        return true;
    }

    function verifyProof(Ics23Helper.ExistenceProof memory exProof)
        internal
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

    function verifyProofWithKV(
        bytes memory key,
        bytes memory value,
        Ics23Helper.ExistenceProof memory exProof
    ) external view returns (bool) {
        // Verify membership
        verify(exProof, getIavlSpec(), relayNetworkHash, key, value);

        return true;
    }

    function queryRootCalculation(Ics23Helper.ExistenceProof memory proof)
        internal
        pure
        returns (bytes memory)
    {
        return bytes(calculate(proof));
    }
}
