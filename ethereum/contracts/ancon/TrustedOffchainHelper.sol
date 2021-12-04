//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.4;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/Address.sol";

abstract contract TrustedOffchainHelper is Ownable {
    using ECDSA for bytes32;
    using Address for address;
    string public url;
    address private _signer;
    mapping(bytes32 => bool) executed;
    error OffchainLookup(string url, bytes prefix);

    event ProofAccepted(address sender, bytes32 signatureHash);

    constructor() {}

    function setUrl(string memory url_) external onlyOwner {
        url = url_;
    }

    function setSigner(address signer_) external onlyOwner {
        _signer = signer_;
    }

    function getSigner() external view returns (address) {
        return _signer;
    }

    function isValidProof(bytes32 digest, bytes memory signature)
        internal
        returns (bool           )
    {
        bytes32 r;
        bytes32 s;
        uint8 v;
        // ecrecover takes the signature parameters, and the only way to get them
        // currently is to use assembly.
        assembly {
            r := mload(add(signature, 0x20))
            s := mload(add(signature, 0x40))
            v := byte(0, mload(add(signature, 0x60)))
        }

        return digest.recover((v + 27), r, s) == _signer;
    }
}
