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
    // bytes public relayNetworkHash;

    IERC20 public stablecoin;
    uint256 public protocolFee = 0;
    uint256 public accountRegistrationFee = 0;

    mapping(bytes => bytes) public accountProofs; //did user-assigned proof key
    mapping(address => bytes) public accountByAddrProofs; //proof key-assigned eth address
    mapping(bytes => bool) public proofs; //if proof key was submitted to the blockchain

    //index is network rootkey, value= signature (chainAnumericIndex + chainBnumericIndex)
    //chainId is not numeric, it must be translated
    //must do easy recover first & then set it
    mapping(bytes32 => address) public whitelistedDagGraph;

    mapping(bytes32 => bytes) public relayerHashTable;

    event Withdrawn(address indexed paymentAddress, uint256 amount);

    event ServiceFeePaid(address indexed from, uint256 fee);

    event HeaderUpdated(bytes hash);
    event ProofPacketSubmitted(bytes key, bytes packet, bytes32 moniker);
    event AccountRegistered(bool enrolledStatus, bytes key, bytes value);

    constructor(address tokenAddress) public {
        owner = msg.sender;
        stablecoin = IERC20(tokenAddress);
    }

    //Must make payable
    function setWhitelistedDagGraph(
        bytes32 moniker,
        address dagAddress,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public {
        require(owner == msg.sender);
        address result = ecrecover(moniker, v, r, s);
        require(dagAddress == result);
        whitelistedDagGraph[moniker] = result;
    }

    //Must make payable
    function updateRelayerHeader(bytes32 moniker, bytes memory rootHash)
        public
    {
        require(whitelistedDagGraph[moniker] == msg.sender);
        // require(msg.sender == relayer);
        relayerHashTable[moniker] = rootHash;
        // emit HeaderUpdated(rootHash);
    }

    function setPaymentToken(address tokenAddress) public {
        require(owner == msg.sender);
        stablecoin = IERC20(tokenAddress);
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
    {
        require(
            stablecoin.balanceOf(address(msg.sender)) > 0,
            "no enough balance"
        );
        if ((paymentType) == ENROLL_PAYMENT) {
            require(
                stablecoin.transferFrom(
                    tokenHolder,
                    address(this),
                    accountRegistrationFee
                ),
                "transfer failed for recipient"
            );
            emit ServiceFeePaid(tokenHolder, accountRegistrationFee);
        }
        if ((paymentType) == SUBMIT_PAYMENT) {
            require(
                stablecoin.transferFrom(
                    tokenHolder,
                    address(this),
                    protocolFee
                ),
                "transfer failed for recipient"
            );
            emit ServiceFeePaid(tokenHolder, protocolFee);
        } // Transfer tokens to pay service fee
    }

    function setProtocolFee(uint256 _fee) public {
        require(owner == msg.sender);
        protocolFee = _fee;
    }

    function setAccountRegistrationFee(uint256 _fee) public {
        require(owner == msg.sender);
        accountRegistrationFee = _fee;
    }

    function getProtocolHeader(bytes32 moniker)
        public
        view
        returns (bytes memory)
    {
        return relayerHashTable[moniker];
    }

    function getProof(bytes memory did) public view returns (bytes memory) {
        return accountProofs[did];
    }

    function hasProof(bytes memory key) public view returns (bool) {
        return proofs[key];
    }

    function enrollL2Account(
        bytes32 moniker,
        bytes memory key,
        bytes memory did,
        Ics23Helper.ExistenceProof memory proof
    ) public payable returns (bool) {
        // require(keccak256(proof.key) == keccak256(key), "invalid key");

        require(verifyProof(moniker, proof), "invalid proof");

        require(
            keccak256(key) != keccak256(accountProofs[did]),
            "user already registered"
        );

        protocolPayment(ENROLL_PAYMENT, msg.sender);

        accountProofs[(did)] = key;
        accountByAddrProofs[msg.sender] = key;

        emit AccountRegistered(true, key, did);
        return true;
    }

    // function updateProtocolHeader(bytes memory rootHash) public {
    //     require(owner == msg.sender);
    //     // require(msg.sender == relayer);
    //     relayNetworkHash = rootHash;
    //     emit HeaderUpdated(rootHash);
    // }

    function submitPacketWithProof(
        bytes32 moniker,
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
        require(verifyProof(moniker, userProof), "invalid user proof");
        require(verifyProof(moniker, proof));

        proofs[key] = true;

        protocolPayment(SUBMIT_PAYMENT, sender);

        // 2. Submit event
        emit ProofPacketSubmitted(key, packet, moniker);

        return true;
    }

    function verifyProof(
        bytes32 moniker,
        Ics23Helper.ExistenceProof memory exProof
    ) internal view returns (bool) {
        // Verify membership
        verify(
            exProof,
            getIavlSpec(),
            relayerHashTable[moniker],
            exProof.key,
            exProof.value
        );

        return true;
    }

    function verifyProofWithKV(
        bytes32 moniker,
        bytes memory key,
        bytes memory value,
        Ics23Helper.ExistenceProof memory exProof
    ) external view returns (bool) {
        // Verify membership
        verify(exProof, getIavlSpec(), relayerHashTable[moniker], key, value);

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
