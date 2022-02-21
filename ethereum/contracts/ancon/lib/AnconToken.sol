// SPDX-License-Identifier: MIT

//** Standard ERC20 - Ancon Token */
//** Author IFESA : Ancon Protocol 2022 */

pragma solidity ^0.8.7;

//** remove previous contract and create standard ERC20 contract */
import "./OwnableUpgradeable.sol";
import "./ERC20Upgradeable.sol";

struct DistrubutionDetails {
    uint256 restrictedUntil;
    uint256 amount;
}

contract AnconToken is ERC20Upgradeable, OwnableUpgradeable {
    
    event AddedStakeHolder (
        uint256 amount,
        address wallet,
        uint256 restrictedUntil
    );

    /**
        distribution: wallet => DistrubutionDetails
        type 1. staff: 6 months
        type 2. advisors: 3 months
    */ 
    
    mapping (address => DistrubutionDetails) public distribution;

    mapping (uint8 => uint256) public types;
    uint256 releaseDate;

    function initialize() public initializer {
        __ERC20_init('Ancon Token', 'ANCON');
        __Ownable_init();
        types[1] = 180 days;
        types[2] = 90 days;
        releaseDate = block.timestamp;
        _mint(msg.sender, 2000000 ether);
    }

    function addStakeHolders(address[] memory wallets, uint256[] memory amounts, uint8[] memory shTypeIndexes) public onlyOwner {
        // require amounts and wallets same length
        require(wallets.length == amounts.length && 
            wallets.length == shTypeIndexes.length &&
            amounts.length == shTypeIndexes.length, "Wallets, amounts, indexes length are different");

        for(uint i = 0; i < wallets.length; i++) {
            address _wallet = wallets[i];
            uint256 _amount = amounts[i];
            uint8 _shTypeIndex = shTypeIndexes[i];
            
            // require found type index on types array
            require(types[_shTypeIndex] != 0, "StakeHolder not found");
            uint256 restrictedDays = types[_shTypeIndex];

            uint256 restrictedUntil = releaseDate + restrictedDays;
            
            DistrubutionDetails memory distributionDetail = DistrubutionDetails(restrictedUntil, _amount);

            distribution[_wallet] = distributionDetail;
            emit AddedStakeHolder(_amount, _wallet, restrictedUntil);
        }
    }

    // override the before tranfer
    /* function _beforeTokenTransfer(address sender, address recipient, uint256 amount) internal virtual override {
        if(distrbution[sender]){
            uint256 restrictedUntil = distrbution[sender];
            
        }
    } */
    // withdraw gas

    function mintToWallet(address address_, uint256 amount)
        public
        payable
        onlyOwner
    {
        _mint(address_, amount);
    }

    function burnToken(address address_, uint256 amount)
        public
        payable
        onlyOwner
    {
        _burn(address_, amount);
    }
}
