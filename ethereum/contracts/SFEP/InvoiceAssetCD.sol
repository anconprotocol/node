pragma solidity ^0.8.7;
import {ILendingPool, IProtocolDataProvider, IStableDebtToken} from "../aave-v2/Interfaces.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
// InvoiceAssetRequest contains the request for tokenization
contract InvoiceAssetCD {
    using SafeERC20 for IERC20;

    uint256 public requestCount;
    uint256 public fee;
    IERC20 public token;
    ILendingPool public lendingPool;
    IProtocolDataProvider public dataProvider;
    uint256 chainId = 0;
    address public owner;

    constructor(
        address _owner,
        address tokenERC20,
        address _lendingPool,
        address _dataProvider,
        uint256 chain
    ) public {
        token = IERC20(tokenERC20);
        lendingPool = ILendingPool(_lendingPool);
        dataProvider = IProtocolDataProvider(_dataProvider);
        chainId = chain;        
        owner = _owner;
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
    ) public  {
        require(owner == msg.sender, "only assigned owner can delegate");
        (, address stableDebtTokenAddress, ) = dataProvider
            .getReserveTokensAddresses(asset);
        IStableDebtToken(stableDebtTokenAddress).approveDelegation(
            borrower,
            amount
        );
    }

    function borrow(
        address assetToBorrow,
        uint256 amountToBorrowInWei,
        uint256 interestRateMode,
        uint16 referralCode,
        address delegatorAddress
    ) public {
        lendingPool.borrow(
            assetToBorrow,
            amountToBorrowInWei,
            interestRateMode,
            referralCode,
            delegatorAddress
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
    function withdrawCollateral(address asset) public  {
        require(owner == msg.sender, "only assigned owner can delegate");
        (address aTokenAddress, , ) = dataProvider.getReserveTokensAddresses(
            asset
        );
        uint256 assetBalance = IERC20(aTokenAddress).balanceOf(address(this));
        lendingPool.withdraw(asset, assetBalance, msg.sender);
    }
}
