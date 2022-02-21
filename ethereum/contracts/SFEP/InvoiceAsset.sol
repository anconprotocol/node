pragma solidity ^0.8.7;
import {ILendingPool, IProtocolDataProvider, IStableDebtToken} from "../aave-v2/Interfaces.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "../ancon/IAnconProtocol.sol";
import "../MFNFT/XDVContainerNFT.sol";
import "../MFNFT/IMFNFT.sol";
import "../MFNFT/MFNFT.sol";
import "./InvoiceAssetCD.sol";
import "../ancon/KYX.sol";
import "./ens-dnssec-oracle/RSAVerify.sol";

// InvoiceAsset contains the tokenization scheme
contract InvoiceAsset is Ownable {
    struct Request {
        string cufeId;
        string cafeUri;
        address creator;
        string kyxId;
        string diddoc;
        bool minted;
    }
    event RequestAdded(
        string indexed cufeId,
        string indexed cafeUri,
        string diddoc
    ); //Id must be a did
    event RequestMinted(
        string indexed cufeId,
        string indexed uri,
        address tokenAddress
    ); //Id must be a did
    event Withdrawn(address indexed payee, uint256 weiAmount);

    uint256 public requestCount;
    mapping(string => Request) public requests;
    mapping(string => InvoiceAssetCD) public invoiceAssetCDList;
    uint256 public fee;
    IERC20 public token;
    IAnconProtocol public anconprotocol;
    KYX public kyx;
    address public lendingPool;
    address public dataProvider;
    uint256 chainId = 0;

    constructor(
        address _tokenERC20,
        address _ancon,
        address _kyx,
        address _lendingPool,
        address _dataProvider,
        uint256 chain
    ) public {
        token = IERC20(_tokenERC20);
        anconprotocol = IAnconProtocol(_ancon);
        kyx = KYX(_kyx);
        chainId = chain;
        lendingPool = _lendingPool;
        dataProvider = _dataProvider;
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
    // Returns a request
    function getRequest(string memory cufeId) public returns (Request memory) {
        require(requests[cufeId].creator != address(0), "no request found");
        return requests[cufeId];
    }

    // Creates a new request
    function createRequestWithProof(
        bytes32 moniker,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory packetProof
    ) public returns (string memory) {
        require(
            keccak256(anconprotocol.getProtocolHeader(moniker)) !=
                keccak256(""),
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
        (
            string memory cufeId,
            string memory cafeUri,
            string memory kyxid,
            bytes32 category,
            string memory diddoc
        ) = abi.decode(packet, (string, string, string, bytes32, string));
        require(
            keccak256(abi.encodePacked(requests[cufeId].cufeId)) !=
                keccak256(abi.encodePacked(cufeId)),
            "request already exists"
        );
        require(
            kyx.getIssuer(category, kyxid).enabled == true &&
                kyx.getIssuer(category, kyxid).creator == msg.sender,
            "invalid KYC"
        );
        requestCount = requestCount + 1;
        requests[cufeId] = Request({
            cufeId: cufeId,
            cafeUri: cafeUri,
            creator: msg.sender,
            kyxId: kyxid,
            diddoc: diddoc,
            minted: false
        });
        emit RequestAdded(cufeId, cafeUri, diddoc);
        return cufeId;
    }

    // Creates a new request
    function mintCDwithProof(
        bytes32 moniker,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory packetProof
    ) public returns (string memory) {
        // Validate network
        require(
            keccak256(anconprotocol.getProtocolHeader(moniker)) !=
                keccak256(""),
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
        ); //Store recover value of the rsa signature
        (
            string memory cufeId,
            string memory kyxid,
            bytes32 category,
            string memory uri,
            bytes memory n,
            bytes memory e,
            bytes memory sig,
            address tokenAddress // uint256 tokenId, // uint256 shares //Pool address
        ) = abi.decode(
                packet,
                (string, string, bytes32, string, bytes, bytes, bytes, address)
            );

        // Verify packet
        require(
            keccak256(abi.encodePacked(requests[cufeId].cufeId)) ==
                keccak256(abi.encodePacked(cufeId)),
            "request must be created"
        );

        // Verify KYX is enabled and matches creator
        require(
            kyx.getIssuer(category, kyxid).enabled == true &&
                kyx.getIssuer(category, kyxid).creator == msg.sender,
            "invalid KYC"
        );
        // Verify is valid SFEP Invoice
        (bool ok, bytes memory x) = RSAVerify.rsarecover(n, e, sig);
        require(ok == true, "invalid invoice");

        InvoiceAssetCD  asset = new InvoiceAssetCD(
            msg.sender,
            tokenAddress,
            lendingPool,
            dataProvider,
            chainId
        );

        // invoice asset request minted ok
        requests[cufeId].minted = true;

        // invoice asset credit delegation
        invoiceAssetCDList[cufeId] = asset;

        emit RequestMinted(cufeId, uri, tokenAddress);
        return cufeId;
    }
}
