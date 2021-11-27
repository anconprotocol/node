pragma solidity ^0.8.7;
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

//token1 token2, withdraw only to these tokens
contract Faucet {
    address public token1;
    address public token2;
    mapping (address=>uint256) waitingList;

    constructor(address _token1) {
        token1 = _token1;
    }

    event Withdrawn(address indexed paymentAddress, uint256 amount);

    function getTokenFromFaucet(address payable payee) public {
        require(waitingList[payee] > 1 days, "Must wait 1 day to withdraw from this address");
        
        IERC20 token = IERC20 (token1);

        //token.approve(payee, 1 ether);  
        
        require(
            token.allowance(payee, address(this)) >= 2 ether,
            "Invalid token allowance"
        );

        uint256 balance = token.balanceOf(address(this));

        require(balance > (1 ether), "Balance must be higher than 1 ether ");
        require(token.transfer(payee, 1 ether), "Ancon Faucet: Transfer failed");
        
        waitingList[payee] = block.timestamp;
        
        emit Withdrawn(payee, balance);
    }
}
