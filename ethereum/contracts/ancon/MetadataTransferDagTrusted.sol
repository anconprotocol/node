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

    function getDigest(
        string memory toAddress,
        string memory tokenId,
        string memory metadataCid,
        string memory fromOwner,
        string memory toOwner,
        string memory resultCid
    ) public view returns (bytes32) {
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
                        tokenId
                    )
                )
            )
        );
        return digest;
    }

    /**
     * @dev Requests a DAG contract offchain execution
     */
    function request(address toAddress, uint256 tokenId)
        public
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
        address toAddress,
        uint256 tokenId,
        MetadataTransferProofPacket memory proof
    ) external returns (bool) {
        if (executed[proof.signature]) {
            return false;
        } else {
            bytes32 digest = keccak256(
                abi.encodePacked(
                    "\x19Ethereum Signed Message:\n32",
                    keccak256(
                        abi.encodePacked(
                            proof.metadataCid,
                            proof.fromOwner,
                            proof.resultCid,
                            proof.toOwner,
                            toAddress,
                            tokenId
                        )
                    )
                )
            );

            address recovered = digest.recover(digest, proof.signature);

            require(
                _signer == recovered,
                "Signer is not the signer of the token"
            );
            executed[proof.signature] = true;
            bytes memory data = abi.encodePacked(toAddress, tokenId);
            _onDagContractResponseReceived(
                toAddress,
                address(this),
                msg.sender,
                proof.metadataCid,
                proof.resultCid,
                data
            );
            return true;
        }
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
                    revert("DagContractTrusted: invalid receiver implementer");
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
