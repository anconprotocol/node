// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract AnconProtocol is ICS23 {
    struct SubscriptionTier {
        address token;
        uint256 amount;
        uint256 amountStaked;
        uint256 includedBlocks;
        bytes32 id;
        uint256 incentiveBlocksMonthly;
        uint256 incentivePercentageMonthly;
        uint256 includedBlocksStarted;
        uint256 setupFee;
    }
    address public owner;
    address public relayer;

    IERC20 public stablecoin;
    uint256 chainId = 0;

    mapping(bytes => bytes) public accountProofs; //did user-assigned proof key
    mapping(address => bytes) public accountByAddrProofs; //proof key-assigned eth address
    mapping(bytes => bool) public proofs; //if proof key was submitted to the blockchain

    mapping(bytes32 => address) public whitelistedDagGraph;
    mapping(bytes32 => SubscriptionTier) public tiers;
    mapping(address => SubscriptionTier) public dagGraphSubscriptions;
    mapping(address => uint256) public totalHeaderUpdatesByDagGraph;
    mapping(address => mapping(address => uint256))
        public totalSubmittedByDagGraphUser;

    uint256 public seq;
    mapping(address => uint256) public nonce;
    mapping(bytes32 => bytes) public latestRootHashTable;
    mapping(bytes32 => mapping(uint256 => bytes)) public relayerHashTable;
    uint256 public INCLUDED_BLOCKS_EPOCH = 200000; // 200 000 chain blocks

    event Withdrawn(address indexed paymentAddress, uint256 amount);

    event ServiceFeePaid(
        address indexed from,
        bytes32 indexed tier,
        bytes32 indexed moniker,
        address token,
        uint256 fee
    );

    event HeaderUpdated(bytes32 indexed moniker);

    event ProofPacketSubmitted(
        bytes indexed key,
        bytes packet,
        bytes32 moniker
    );

    event TierAdded(bytes32 indexed id);

    event TierUpdated(
        bytes32 indexed id,
        address token,
        uint256 fee,
        uint256 staked,
        uint256 includedBlocks
    );

    event AccountRegistered(
        bool enrolledStatus,
        bytes key,
        bytes value,
        bytes32 moniker
    );

    constructor(
        address tokenAddress,
        uint256 network,
        uint256 starterFee,
        uint256 startupFee
    ) public {
        owner = msg.sender;
        stablecoin = IERC20(tokenAddress);
        chainId = network;

        // add tiers
        // crear un solo tier `default` Ancon token, starterFee 0.50, blocks, fee y staked en 0
        addTier(keccak256("starter"), tokenAddress, starterFee, 0, 100, 0);
        addTier(keccak256("startup"), tokenAddress, startupFee, 0, 500, 0);
        addTier(keccak256("pro"), tokenAddress, 0, 0, 1000, 150);
        /* setTierSettings(
            keccak256("pro"),
            tokenAddress,
            500000000,
            1000 ether,
            1000
        ); */
        addTier(keccak256("defi"), tokenAddress, 0, 0, 10000, 1500);
        addTier(keccak256("luxury"), tokenAddress, 0, 0, 100000, 9000);
    }

    // getContractIdentifier is used to identify a contract protocol deployed in a specific chain
    function getContractIdentifier() public view returns (bytes32) {
        return keccak256(abi.encodePacked(chainId, address(this)));
    }

    // verifyContractIdentifier verifies a nonce is from  a specific chain
    function verifyContractIdentifier(
        uint256 usernonce,
        address sender,
        bytes32 hash
    ) public view returns (bool) {
        return
            keccak256(abi.encodePacked(chainId, address(this))) == hash &&
            nonce[sender] == usernonce;
    }

    function getNonce() public view returns (uint256) {
        return nonce[msg.sender];
    }

    // registerDagGraphTier
    function registerDagGraphTier(
        bytes32 moniker,
        address dagAddress,
        bytes32 tier
    ) public payable {
        require(whitelistedDagGraph[moniker] == address(0), "moniker exists");
        require(tier == tiers[tier].id, "missing tier");

        if(tiers[tier].setupFee > 0){
            IERC20 token = IERC20(tiers[tier].token);
            require(token.balanceOf(address(msg.sender)) > tiers[tier].setupFee, "no enough balance");
            require(
                token.transferFrom(
                    msg.sender,
                    address(this),
                    tiers[tier].setupFee
                ),
                "transfer failed for recipient"
            );
        }
        whitelistedDagGraph[moniker] = dagAddress;
        dagGraphSubscriptions[dagAddress] = tiers[tier];
    }

    // updateRelayerHeader updates offchain dag graphs signed by dag graph key pair
    function updateRelayerHeader(
        bytes32 moniker,
        bytes memory rootHash,
        uint256 height
    ) public payable {
        require(msg.sender == whitelistedDagGraph[moniker], "invalid user");

        SubscriptionTier memory t = dagGraphSubscriptions[msg.sender];
        IERC20 token = IERC20(tiers[t.id].token);
        require(token.balanceOf(address(msg.sender)) > 0, "no enough balance");

        if (t.includedBlocks > 0) {
            t.includedBlocks = t.includedBlocks - 1;
        } else {
            // tier has no more free blocks for this epoch, charge protocol fee
            require(
                token.transferFrom(
                    msg.sender,
                    address(this),
                    tiers[t.id].amount
                ),
                "transfer failed for recipient"
            );
        }
        // reset tier includede blocks every elapsed epoch
        uint256 elapsed = block.number - t.includedBlocksStarted;
        if (elapsed > INCLUDED_BLOCKS_EPOCH) {
            // must always read from latest tier settings
            t.includedBlocks = tiers[t.id].includedBlocks;
            t.includedBlocksStarted = block.number;
        }
        // set hash
        relayerHashTable[moniker][height] = rootHash;
        latestRootHashTable[moniker] = rootHash;
        emit ServiceFeePaid(
            msg.sender,
            moniker,
            t.id,
            tiers[t.id].token,
            tiers[t.id].amount
        );

        seq = seq + 1;
        totalHeaderUpdatesByDagGraph[msg.sender] =
            totalHeaderUpdatesByDagGraph[msg.sender] +
            1;
        emit HeaderUpdated(moniker);
    }

    // setPaymentToken sets token used for protocol fees
    function setPaymentToken(address tokenAddress) public {
        require(owner == msg.sender);
        stablecoin = IERC20(tokenAddress);
    }

    // addTier
    function addTier(
        bytes32 id,
        address tokenAddress,
        uint256 amount,
        uint256 amountStaked,
        uint256 includedBlocks,
        uint256 setupFee
    ) public {
        require(owner == msg.sender, "invalid owner");
        require(tiers[id].id != id, "tier already in use");
        tiers[id] = SubscriptionTier({
            token: tokenAddress,
            amount: amount,
            amountStaked: amountStaked,
            includedBlocks: includedBlocks,
            id: id,
            incentiveBlocksMonthly: 0,
            incentivePercentageMonthly: 0,
            includedBlocksStarted: block.number,
            setupFee: setupFee
        });
        emit TierAdded(id);
    }

    // setTierSettings
    function setTierSettings(
        bytes32 id,
        address tokenAddress,
        uint256 amount,
        uint256 amountStaked,
        uint256 includedBlocks,
        uint256 setupFee
    ) public {
        require(owner == msg.sender, "invalid owner");
        require(tiers[id].id == id, "missing tier");
        tiers[id].token = tokenAddress;
        tiers[id].amount = amount;
        tiers[id].amountStaked = amountStaked;
        tiers[id].includedBlocks = includedBlocks;
        tiers[id].setupFee = setupFee;
        // incentiveBlocksMonthly: 0,
        // incentivePercentageMonthly: 0
        emit TierUpdated(
            id,
            tokenAddress,
            amount,
            amountStaked,
            includedBlocks
        );
    }

    // withdraws gas token, must be admin
    function withdraw(address payable payee) public {
        require(owner == msg.sender);
        uint256 b = address(this).balance;
        (bool sent, bytes memory data) = payee.call{value: b}("");
    }

    // withdraws protocol fee token, must be admin
    function withdrawToken(address payable payee, address erc20token) public {
        require(owner == msg.sender);
        uint256 balance = IERC20(erc20token).balanceOf(address(this));

        // Transfer tokens to pay service fee
        require(IERC20(erc20token).transfer(payee, balance), "transfer failed");

        emit Withdrawn(payee, balance);
    }

    function getProtocolHeader(bytes32 moniker)
        public
        view
        returns (bytes memory)
    {
        return latestRootHashTable[moniker];
        require(sent, "Failed to send Ether");

        emit Withdrawn(payee, b);
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
    ) public returns (bool) {
        require(keccak256(proof.key) == keccak256(key), "invalid key");

        require(verifyProof(moniker, proof), "invalid proof");

        require(
            keccak256(key) != keccak256(accountProofs[did]),
            "user already registered"
        );

        totalSubmittedByDagGraphUser[whitelistedDagGraph[moniker]][msg.sender] =
            totalSubmittedByDagGraphUser[whitelistedDagGraph[moniker]][
                msg.sender
            ] +
            1;
        accountProofs[(did)] = key;
        accountByAddrProofs[msg.sender] = key;

        emit AccountRegistered(true, key, did, moniker);
        return true;
    }

    // submitPacketWithProof registers packet onchain using ICS23 proofs, multi tenant using dag graph moniker
    function submitPacketWithProof(
        bytes32 moniker,
        address sender,
        Ics23Helper.ExistenceProof memory userProof,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory proof
    ) external returns (bool) {
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
        totalSubmittedByDagGraphUser[whitelistedDagGraph[moniker]][sender] =
            totalSubmittedByDagGraphUser[whitelistedDagGraph[moniker]][
                sender
            ] +
            1;

        nonce[sender] = nonce[sender] + 1;
        // 2. Submit event
        emit ProofPacketSubmitted(key, packet, moniker);

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
            latestRootHashTable[moniker],
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
            latestRootHashTable[moniker],
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
