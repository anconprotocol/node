//SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.7;

struct VerifiableCredential {
    address issuer;
    address subject;
    bytes data;
    uint256 validFrom;
    uint256 validTo;
}
