//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.4;

import "./IDagContractTrustedReceiver.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract DagContractTrusted is Ownable, IDagContractTrustedReceiver {
    using ECDSA for bytes32;

    string public url;
    address private _signer;
    mapping(bytes32 => bool) executed;
    error OffchainLookup(string url, bytes prefix);
    struct DagContractRequestProof {
        string schemaCid;
        string dataSourceCid;
        string variables;
        string contractMutation;
        string result;
        address addr;
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
    function request(address addr) public returns (bytes32) {
        revert OffchainLookup(
            url,
            abi.encodeWithSignature(
                "requestWithProof(address addr, DagContractRequestProof memory proof)",
                contractMutation,
                addr
            )
        );
    }

    /**
     * @dev Requests a DAG contract offchain execution with proof
     */
    function requestWithProof(
        address addr,
        DagContractRequestProof memory proof
    ) external returns (bool) {
        if (executed[proof.signature]) {
            return false;
        } else {
            bytes32 digest = keccak256(
                abi.encodePacked(
                    "\x19Ethereum Signed Message:\n32",
                    keccak256(
                        abi.encodePacked(
                            proof.schemaCid,
                            proof.dataSourceCid,
                            proof.variables,
                            proof.contractMutation,
                            proof.result,
                            addr
                        )
                    )
                )
            );
            address recovered = digest.recover(proof.signature);

            require(
                _signer == recovered,
                "Signer is not the signer of the token"
            );
            executed(proof.signature);
            bytes memory none = bytes("");
            onDagContractResponseReceived(
                this.address,
                msg.sender,
                proof.dataSourceCid, 
                proof.result,
                none
            );
            return true;
        }
    }



    /**
     * @dev Receives the DAG contract execution result
     */
    function onDagContractResponseReceived(
        address operator,
        address from,
        string memory parentCid,
        string memory newCid,
        bytes calldata data
    ) external returns (bytes4) {
        return DagContractTrusted.onDagContractResponseReceived.selector;  
    }
}
