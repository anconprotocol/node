// global['fetch'] = require('node-fetch')
import('node-fetch')
const ethers = require('ethers');
const Bluebird = require('bluebird');
const jayson = require('jayson');
require('dotenv').config({ path: '../.env' });
const BigNumber = require('bignumber.js')
const fs = require('fs')
const Web3 = require("web3");

const XDVNFT = artifacts.require('XDVNFT')
const DAI = artifacts.require('DAI')
const MetadataTransferDagTrusted = artifacts.require("MetadataTransferDagTrusted");
const {
    CHAIN_A,
    PKEY
} = process.env;
const metadataTransferAbi = require('../build/contracts/MetadataTransferDagTrusted');
const { ErrorFragment } = require('@ethersproject/abi');

const provider = new ethers.providers.JsonRpcProvider(CHAIN_A);
const signer = new ethers.Wallet(PKEY);

const erc20 = new ethers.Contract(XDVNFT.address, metadataTransferAbi.abi, provider);
const SENDER_ADDRESS = signer.address
console.log({
    CHAIN_A
})

const metadataTransferIface = new ethers.utils.Interface(metadataTransferAbi.abi)

function parseEthError(ethErr, iface) {
    try {
        const keys = Object.keys(ethErr.data)
        return (iface.decodeErrorResult('OffchainLookup', ethErr.data[keys[0]].return))
    } catch (e) {
        return null;
    }
}

contract('Ancon - ICS23 Javascript to ABI', (accounts) => {
    // let erc20Contract;
    // let controllerContract;
    // let documentMinterAddress;
    let metadataTransferDagTrusted
    let xdvnft
    let daiContract

    // Initialize the contracts and make sure they exist
    before(async () => {
        ; ({ metadataTransferDagTrusted, xdvnft, daiContract } = await Bluebird.props({
            metadataTransferDagTrusted: MetadataTransferDagTrusted.deployed(),
            xdvnft: XDVNFT.deployed(),
            //erc20Contract: TestUSDC.deployed(),
            //con trollerContract: XDVController.deployed(),
            daiContract: DAI.deployed(),
        }))
    })

    xdescribe('when requesting a transfer ownership', () => {
        it('should revert', async () => {

            try {
                await metadataTransferDagTrusted.request(xdvnft.address, 1);
            } catch (e) {
                if (e.message.match(/OffchainLookup/)) {
                    assert.equal(e.message.match(/OffchainLookup/), true)

                } else {
                    console.log({ e })
                    assert.notEqual(e.message.match(/OffchainLookup/), true)
                }
            }

        })
    })

    describe('when requesting a transfer ownership with proof', () => {
        it('should succeed', async () => {
            const toAddress = xdvnft.address;
            const tokenId = 1;
            const iface = new ethers.utils.Interface(metadataTransferAbi.abi)
            const requestWithProofIface = iface.requestWithProof
            try {
                await metadataTransferDagTrusted.request(xdvnft.address, 1);
            } catch (e) {
                const hasError = parseEthError(e, iface)
                if (hasError) {
                    const url = "http://localhost:7788/rpc"
                    const inputdata = iface.encodeFunctionData("requestWithProof", [toAddress, tokenId])
                    const body = {
                        jsonrpc: '2.0',
                        method: 'durin_call',
                        params: [{ to: TOKEN_ADDRESS, from: fromAddress, data: inputdata }],
                        id: 1,
                    }
                    const result = await (await node - fetch(url, {
                        method: 'post',
                        body: JSON.stringify(body),
                        headers: { 'Content-Type': 'application/json' },
                    })).json()
                    // const nonce = await provider.getTransactionCount(signer.address)
                    // const signedTransaction = await signer.signTransaction({
                    //     nonce,
                    //     gasLimit: 500000,
                    //     gasPrice: ethers.BigNumber.from("1000000000"),
                    //     to: TOKEN_ADDRESS,
                    //     data: result.result,
                    // });
                    // await provider.sendTransaction(signedTransaction)
                    console.log(result)
                } else {
                    console.log({ e })
                }
            }
        })
    })
})
