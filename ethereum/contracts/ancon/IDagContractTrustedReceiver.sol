// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

/**
 * @title DAG Contract receiver interface
 * @dev Interface for any contract that wants to support DAG offchain.
 */
interface IDagContractTrustedReceiver {
    /**
     * @dev Whenever an {IERC721} `tokenId` token is transferred to this contract via {IERC721-safeTransferFrom}
     * by `operator` from `from`, this function is called.
     *
     * It must return its Solidity selector to confirm the token transfer.
     * If any other value is returned or the interface is not implemented by the recipient, the transfer will be reverted.
     *
     * The selector can be obtained in Solidity with `IERC721.onERC721Received.selector`.
     */
    function onDagContractResponseReceived(
        address operator,
        address from,
        string memory parentCid,
        string memory newCid,
        bytes calldata data
    ) external returns (bytes4);
}
