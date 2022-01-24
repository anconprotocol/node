// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;
import "./ics23/ics23.sol";
import "./ics23/Ics23Helper.sol";

interface IWXDV {
    function submitMintWithProof(
        address sender,
        uint256 newItemId,
        bytes32 moniker,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) external payable returns (bool);

    function lockWithProof(
        address sender,
        bytes32 moniker,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) external payable returns (uint256);

    function releaseWithProof(
        address sender,
        bytes32 moniker,
        bytes memory key,
        bytes memory packet,
        Ics23Helper.ExistenceProof memory userProof,
        Ics23Helper.ExistenceProof memory proof,
        bytes32 hash
    ) external payable returns (uint256);
}
