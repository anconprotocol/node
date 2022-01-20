// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
contract AnconProtocol is ICS23 {
    bytes32 public ENROLL_PAYMENT = keccak256("ENROLL_PAYMENT");
    bytes32 public ENROLL_DAG = keccak256("ENROLL_DAG");
    bytes32 public SUBMIT_PAYMENT = keccak256("SUBMIT_PAYMENT");

    address public owner;
    address public relayer;
    // bytes public relayNetworkHash;
    struct Header {
        bytes roothash;
        uint256 height;
        uint8 v;
        bytes32 r;
        bytes32 s;
    }

    IERC20 public stablecoin;
    uint256 public protocolFee = 0;
    uint256 public accountRegistrationFee = 0;
    uint256 public dagRegistrationFee = 0;
    uint256 chainId = 0;

    mapping(bytes => bytes) public accountProofs; //did user-assigned proof key
    mapping(address => bytes) public accountByAddrProofs; //proof key-assigned eth address
    mapping(bytes => bool) public proofs; //if proof key was submitted to the blockchain

    //index is network rootkey, value= signature (chainAnumericIndex + chainBnumericIndex)
    //chainId is not numeric, it must be translated
    //must do easy recover first & then set it
    mapping(bytes32 => address) public whitelistedDagGraph;

    mapping(bytes32 => Header) public relayerHashTable;

    event Withdrawn(address indexed paymentAddress, uint256 amount);

    event ServiceFeePaid(address indexed from, uint256 fee);

    event HeaderUpdated(bytes32 indexed moniker, Header header);

    event ProofPacketSubmitted(
        bytes key,
        bytes packet,
        bytes32 moniker,
        Header header
    );

    event AccountRegistered(
        bool enrolledStatus,
        bytes key,
        bytes value,
        bytes32 moniker,
        Header header
    );

    constructor(address tokenAddress, uint256 network) public {
        owner = msg.sender;
        stablecoin = IERC20(tokenAddress);
        chainId = network;
    }

    // getContractIdentifier is used to identify an offchain proof in any chain  
    function getContractIdentifier() public returns (bytes32) {
        return keccak256(abi.encodePacked(
                chainId,
                address(this)
            ));
    }

    // setWhitelistedDagGraph registers offchain graphs by protocol admin
    function setWhitelistedDagGraph(
        bytes32 moniker,
        address dagAddress,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public payable {
        require(owner == msg.sender, "invalid owner");
        address result = ecrecover(moniker, v, r, s);
        require(dagAddress == result, "invalid address");

        protocolPayment(ENROLL_DAG, msg.sender);

        whitelistedDagGraph[moniker] = result;
    }

    // updateRelayerHeader updates offchain dag graphs signed by dag graph key pair
    function updateRelayerHeader(
        bytes32 moniker,
        bytes memory rootHash,
        uint256 height,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public {
        address result = ecrecover(moniker, v, r, s);
        require(
            whitelistedDagGraph[moniker] == result,
            "invalid network moniker"
        );
        // TODO:  Check to  see if  signer has n amount of token staked
        relayerHashTable[moniker] = Header(rootHash, height, v, r, s);
        emit HeaderUpdated(moniker, relayerHashTable[moniker]);
    }

    // setPaymentToken sets token used for protocol fees
    function setPaymentToken(address tokenAddress) public {
        require(owner == msg.sender);
        stablecoin = IERC20(tokenAddress);
    }

    // withdraws gas token, must be admin 
    function withdraw(address payable payee) public {
        require(owner == msg.sender);
        uint256 b = address(this).balance;
        (bool sent, bytes memory data) = payee.call{value: b}("");
        require(sent, "Failed to send Ether");

        emit Withdrawn(payee, b);
    }

    // withdraws protocol fee token, must be admin 
    function withdrawToken(address payable payee, address erc20token) public {
        require(owner == msg.sender);
        uint256 balance = IERC20(erc20token).balanceOf(address(this));

        // Transfer tokens to pay service fee
        require(IERC20(erc20token).transfer(payee, balance), "transfer failed");

        emit Withdrawn(payee, balance);
    }

    // protocolPayment handles contract payment protocol fee types 
    function protocolPayment(bytes32 paymentType, address tokenHolder)
        internal
    {
        require(
            stablecoin.balanceOf(address(msg.sender)) > 0,
            "no enough balance"
        );
        if ((paymentType) == ENROLL_DAG) {
            require(
                stablecoin.transferFrom(
                    tokenHolder,
                    address(this),
                    dagRegistrationFee
                ),
                "transfer failed for recipient"
            );
            emit ServiceFeePaid(tokenHolder, dagRegistrationFee);
        }
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

    function setDagGraphFee(uint256 _fee) public {
        require(owner == msg.sender);
        dagRegistrationFee = _fee;
    }

    function getProtocolHeader(bytes32 moniker)
        public
        view
        returns (Header memory)
    {
        return relayerHashTable[moniker];
    }

    function getProof(bytes memory did) public view returns (bytes memory) {
        return accountProofs[did];
    }

    function hasProof(bytes memory key) public view returns (bool) {
        return proofs[key];
    }

    // enrollL2Account registers offchain did user onchain using ICS23 proofs, multi tenant using dag graph moniker 
    function enrollL2Account(
        bytes32 moniker,
        bytes memory key,
        bytes memory did,
        Ics23Helper.ExistenceProof memory proof
    ) public payable returns (bool) {
        require(keccak256(proof.key) == keccak256(key), "invalid key");

        require(verifyProof(moniker, proof), "invalid proof");

        require(
            keccak256(key) != keccak256(accountProofs[did]),
            "user already registered"
        );

        protocolPayment(ENROLL_PAYMENT, msg.sender);

        accountProofs[(did)] = key;
        accountByAddrProofs[msg.sender] = key;

        emit AccountRegistered(
            true,
            key,
            did,
            moniker,
            relayerHashTable[moniker]
        );
        return true;
    }


    // enrollL2Account registers packet onchain using ICS23 proofs, multi tenant using dag graph moniker
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
        emit ProofPacketSubmitted(
            key,
            packet,
            moniker,
            relayerHashTable[moniker]
        );

        return true;
    }


    // verifies ICS23 proofs, multi tenant using dag graph moniker
    function verifyProof(
        bytes32 moniker,
        Ics23Helper.ExistenceProof memory exProof
    ) internal view returns (bool) {
        // Verify membership
        verify(
            exProof,
            getIavlSpec(),
            relayerHashTable[moniker].roothash,
            exProof.key,
            exProof.value
        );

        return true;
    }

    // verifies ICS23 proofs with key and value, multi tenant using dag graph moniker
    function verifyProofWithKV(
        bytes32 moniker,
        bytes memory key,
        bytes memory value,
        Ics23Helper.ExistenceProof memory exProof
    ) external view returns (bool) {
        // Verify membership
        verify(
            exProof,
            getIavlSpec(),
            relayerHashTable[moniker].roothash,
            key,
            value
        );

        return true;
    }

    // calculates root hash
    function queryRootCalculation(Ics23Helper.ExistenceProof memory proof)
        internal
        pure
        returns (bytes memory)
    {
        return bytes(calculate(proof));
    }
}
