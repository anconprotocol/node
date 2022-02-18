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
// InvoiceAssetRequest contains the request for tokenization
contract InvoiceAssetRequest is Ownable {

    struct Request {
        string cufeId;
        string cafeUri;
        address creator;
        uint256 kyxId;
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
        address tokenAddress,
        uint256 tokenId
    ); //Id must be a did
    event Withdrawn(address indexed payee, uint256 weiAmount);

    uint256 public requestCount;
    mapping(string => Request) public requests;
    uint256 public fee;
    IERC20 public token;
    IAnconProtocol public anconprotocol;
    uint256 chainId = 0;

    constructor(
        address tokenERC20,
        address ancon,
        address _lendingPool,
        address _dataProvider,
        uint256 chain
    ) public {
        token = IERC20(tokenERC20);
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
            uint256 kyxId,
            string memory diddoc
        ) = abi.decode(packet, (string, string, uint256, string));
        require(
            keccak256(abi.encodePacked(requests[cufeId].cufeId)) !=
                keccak256(abi.encodePacked(cufeId)),
            "request already exists"
        );
        // TODO: Verify the KYX
        requestCount = requestCount + 1;
        requests[cufeId] = Request({
            cufeId: cufeId,
            cafeUri: cafeUri,
            creator: msg.sender,
            kyxId: kyxId,
            diddoc: diddoc,
            minted: false
        });
        emit RequestAdded(cufeId, cafeUri, diddoc);
        return cufeId;
    }

    // Creates a new request
    function mintMFNFTwithProof(
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
        ); //Store recover valu of the rsa signature
        (
            string memory cufeId,
            string memory uri,
            bytes memory n,
            bytes memory e,
            bytes memory sig,
            address tokenAddress,
            uint256 tokenId,
            uint256 shares
        ) = //Pool address
            abi.decode(
                packet,
                (string, string, bytes, bytes, bytes, address, uint256, uint256)
            );
        require(
            keccak256(abi.encodePacked(requests[cufeId].cufeId)) ==
                keccak256(abi.encodePacked(cufeId)),
            "request must be created"
        );
        // TODO:KYX/RSA
        // WIP MINT...
        // WIP create InvoiceAssetCD
        //Example, if this is an invoice stored as an nft with 1000 shares
        //mfnt with the shares and then transfer to pool address
        requests[cufeId].minted = true;
        emit RequestMinted(cufeId, uri, tokenAddress, tokenId);
        return cufeId;
    }
}
