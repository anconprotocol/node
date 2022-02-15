// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "../ics23/ics23.sol";
import "../ics23/Ics23Helper.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

interface IAnconProtocol {
    function verifyContractIdentifier(
        uint256 usernonce,
        address sender,
        bytes32 hash
    ) external view returns (bool);

    function getProtocolHeader(bytes32 moniker)
        external
        view
        returns (bytes memory);

    function getContractIdentifier() external view returns (bytes32);

    function submitPacketWithProof(
        bytes32 moniker,
        address sender,
        Ics23Helper.ExistenceProof memory userProof,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory proof
    ) external returns (bool);

    function verifyProof(
        bytes32 moniker,
        Ics23Helper.ExistenceProof memory exProof
    ) external view returns (bool);

    function verifyProofWithKV(
        bytes32 moniker,
        bytes memory key,
        bytes memory value,
        Ics23Helper.ExistenceProof memory exProof
    ) external view returns (bool);
}
