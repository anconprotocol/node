// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetMinterPauser.sol";

contract DAI is ERC20PresetMinterPauser, Ownable {
    constructor() ERC20PresetMinterPauser("DAI", "DAI") {
        mint(msg.sender, 1000000 ether);
    }
}
