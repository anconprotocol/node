//SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.7;

import "./ClaimTypes.sol";

interface IClaimsVerifier {
    function verifyCredential(
        VerifiableCredential
         memory vc,
        uint8 v,
        bytes32 r,
        bytes32 s
    )
        external
        view
        returns (
            bool,
            bool,
            bool,
            bool,
            bool
        );

    function verifySigner(
        VerifiableCredential memory vc,
        bytes calldata _signature
    ) external view returns (bool);
}
