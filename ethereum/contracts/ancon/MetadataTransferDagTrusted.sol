//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.4;

import "./IDagContractTrustedReceiver.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/Address.sol";

contract MetadataTransferDagTrusted is Ownable {
    using ECDSA for bytes32;
    using Address for address;
    string public url;
    address private _signer;
    mapping(bytes32 => bool) executed;
    error OffchainLookup(string url, bytes prefix);
    struct MetadataTransferProofPacket {
        string metadataCid;
        string fromOwner;
        string toOwner;
        string resultCid;
        string toAddress;
        string tokenId;
        string prefix;
        bytes32 signature;
    }

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

    /**
     * @dev Requests a DAG contract offchain execution
     */
    function request(address toAddress, uint256 tokenId)
        external    
        returns (bytes32)
    {
        revert OffchainLookup(
            url,
            abi.encodeWithSignature(
                "requestWithProof(address toAddress, uint256 tokenId, MetadataTransferProofPacket memory proof)",
                toAddress,
                tokenId
            )
        );
    }

    /**
     * @dev Requests a DAG contract offchain execution with proof
     */
    function requestWithProof(
        string memory toAddress,
        string memory tokenId,
        bytes memory proof
    ) external returns (bool) {
        (
            bytes memory metadataCid,
            bytes memory fromOwner,
            bytes memory resultCid,
            bytes memory toOwner,
            ,
            ,
            bytes memory prefix,
            bytes memory signature
        ) = abi.decode(
                proof,
                (bytes, bytes, bytes, bytes, bytes, bytes, bytes, bytes)
            );

        if (executed[bytes32(signature)]) {
            revert("metadata dag transfer:  invalid proof");
            return false;
        } else {
            bytes32 digest = keccak256(
                abi.encodePacked(
                    "\x19Ethereum Signed Message:\n32",
                    keccak256(
                        abi.encodePacked(
                            metadataCid,
                            fromOwner,
                            resultCid,
                            toOwner,
                            toAddress,
                            tokenId,
                            prefix
                        )
                    )
                )
            );

            require(
                _signer == isValidProof(digest,  signature),
                "Signer is not the signer of the token"
            );
            {
                executed[bytes32(signature)] = true;
                _onDagContractResponseReceived(
                    address(bytes20(bytes(toAddress))),
                    address(this),
                    msg.sender,
                    string(metadataCid),
                    (string(resultCid)),
                    proof
                );
            }
            return true;
        }
    }

    function isValidProof(bytes32 digest,bytes  memory   signature) internal returns (address){

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

            return digest.recover((v + 27), r, s);

    }
    /**
     * @dev Receives the DAG contract execution result
     */
    // function onDagContractResponseReceived(
    //     address operator,
    //     address from,
    //     string memory parentCid,
    //     string memory newCid,
    //     bytes calldata data
    // ) external returns (bytes4) {
    //     return IDagContractTrustedReceiver.onDagContractResponseReceived.selector;
    // }

    function _onDagContractResponseReceived(
        address to,
        address operator,
        address from,
        string memory parentCid,
        string memory newCid,
        bytes memory _data
    ) private returns (bool) {
        if (to.isContract()) {
            try
                IDagContractTrustedReceiver(to).onDagContractResponseReceived(
                    operator,
                    from,
                    parentCid,
                    newCid,
                    _data
                )
            returns (bytes4 retval) {
                return
                    retval ==
                    IDagContractTrustedReceiver
                        .onDagContractResponseReceived
                        .selector;
            } catch (bytes memory reason) {
                if (reason.length == 0) {
                    revert(
                        "metadata dag transfer: invalid receiver implementer"
                    );
                } else {
                    assembly {
                        revert(add(32, reason), mload(reason))
                    }   
                }
            }
        } else {
            return true;
        }
    }
}
