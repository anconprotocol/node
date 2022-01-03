// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

interface IAnconProtocol {
    function submitPacketWithProof(
        bytes memory key,
        bytes memory packet,
        ExistenceProof memory proof
    ) external payable returns (bool);

    function verifyProof(ExistenceProof memory exProof)
        external
        view
        returns (bool);
}
