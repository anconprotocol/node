pragma solidity ^0.8.7;
import {ILendingPool, IProtocolDataProvider, IStableDebtToken} from "../aave-v2/Interfaces.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "../ancon/IAnconProtocol.sol";
import "../MFNFT/XDVContainerNFT.sol";
import "../MFNFT/IMFNFT.sol";
import "../MFNFT/MFNFT.sol";

// InvoiceAssetRequest contains the request for tokenization
contract InvoiceAssetCD is Ownable {
    using SafeERC20 for IERC20;
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
    ILendingPool public lendingPool;
    IProtocolDataProvider public dataProvider;
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
        lendingPool = ILendingPool(_lendingPool);
        dataProvider = IProtocolDataProvider(_dataProvider);
        chainId = chain;
    }

    /**
     * Deposits collateral into the Aave, to enable credit delegation
     * This would be called by the delegator.
     * @param asset The asset to be deposited as collateral
     * @param amount The amount to be deposited as collateral
     * @param isPull Whether to pull the funds from the caller, or use funds sent to this contract
     *  User must have approved this contract to pull funds if `isPull` = true
     *
     */
    function depositCollateral(
        address asset,
        uint256 amount,
        bool isPull
    ) public {
        if (isPull) {
            IERC20(asset).safeTransferFrom(msg.sender, address(this), amount);
        }
        IERC20(asset).safeApprove(address(lendingPool), amount);
        lendingPool.deposit(asset, amount, address(this), 0);
    }

    /**
     * Approves the borrower to take an uncollaterised loan
     * @param borrower The borrower of the funds (i.e. delgatee)
     * @param amount The amount the borrower is allowed to borrow (i.e. their line of credit)
     * @param asset The asset they are allowed to borrow
     *
     * Add permissions to this call, e.g. only the owner should be able to approve borrowers!
     */
    function approveBorrower(
        address borrower,
        uint256 amount,
        address asset
    ) public {
        (, address stableDebtTokenAddress, ) = dataProvider
            .getReserveTokensAddresses(asset);
        IStableDebtToken(stableDebtTokenAddress).approveDelegation(
            borrower,
            amount
        );
    }

    /**
     * Repay an uncollaterised loan
     * @param amount The amount to repay
     * @param asset The asset to be repaid
     *
     * User calling this function must have approved this contract with an allowance to transfer the tokens
     *
     * You should keep internal accounting of borrowers, if your contract will have multiple borrowers
     */
    function repayBorrower(uint256 amount, address asset) public {
        IERC20(asset).safeTransferFrom(msg.sender, address(this), amount);
        IERC20(asset).safeApprove(address(lendingPool), amount);
        lendingPool.repay(asset, amount, 1, address(this));
    }

    /**
     * Withdraw all of a collateral as the underlying asset, if no outstanding loans delegated
     * @param asset The underlying asset to withdraw
     *
     * Add permissions to this call, e.g. only the owner should be able to withdraw the collateral!
     */
    function withdrawCollateral(address asset) public onlyOwner {
        (address aTokenAddress, , ) = dataProvider.getReserveTokensAddresses(
            asset
        );
        uint256 assetBalance = IERC20(aTokenAddress).balanceOf(address(this));
        lendingPool.withdraw(asset, assetBalance, msg.sender);
    }

}
