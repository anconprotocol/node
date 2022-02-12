// SPDX-License-Identifier: MIT

//** Standard ERC20 - Ancon Token */
//** Author IFESA : Ancon Protocol 2022 */

pragma solidity 0.8.5;

//** remove previous contract and create standard ERC20 contract */
import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20BurnableUpgradeable.sol";

event AddedStakeHolder {
    uint256 amount,
    address wallet,
    uint256 restrictedUntil
}

struct StakeHolderType {
    uint8 shType,
    uint256 restrictedDays,
}

contract SampleToken is ERC20Upgradeable, ERC20BurnableUpgradeable {
    
    /**
        distribution: wallet => restricted until timestamp date
        type 1. staff: 6 months
        type 2. advisors: 3 months
    */ 
    
    mapping (address => uint256) public distribution;

    StakeHolderType[] public types;
    uint256 releaseDate;

    function initialize() public initializer {
        __ERC20_init('Ancon Token', 'ANCON');
        types.push(StakeHolderType(1, 180 days));
        types.push(StakeHolderType(2, 90 days));
        releaseDate = block.timestamp;
    }

    function addStakeHolders(address[] wallets, uint256[] amounts, uint shTypeIndex) onlyOwner {
        // require amounts and wallets same length
        // require found type index on types array
        require(wallets.length == amounts.length, "Wallets and amounts length are different");
        require(types[shTypeIndex].shType, "StakeHolder not found");

        StakeHolderType sht = types[shTypeIndex];
        for(uint i = 0; i < wallets.length; i++) {
            address _wallet = wallets[i];
            uint256 _amount = amounts[i];
            uint256 restrictedDays = sht.restrictedDays;
            uint256 restrictedUntil = releaseDate.add(restrictedDays);

            distribution[_wallet] = restrictedUntil;
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

    function burnToken(uint256 amount)
        public
        payable
        onlyOwner
    {
        burn(address_);
    }
}
