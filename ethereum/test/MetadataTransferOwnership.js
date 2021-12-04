global['fetch'] = import('node-fetch')
const ethers = require('ethers')
const Bluebird = require('bluebird')
const jayson = require('jayson')
require('dotenv').config({ path: '../.env' })
const BigNumber = require('bignumber.js')
const fs = require('fs')
const Web3 = require('web3')

const XDVNFT = artifacts.require('XDVNFT')
const DAI = artifacts.require('DAI')

const { CHAIN_A, PKEY } = process.env
const { fetchJson } = require('@ethersproject/web')

const provider = new ethers.providers.JsonRpcProvider(CHAIN_A)
const signer = new ethers.Wallet(PKEY)
const SENDER_ADDRESS = signer.address
console.log({
  CHAIN_A,
})

function parseEthError(ethErr, iface) {
  try {
    const keys = Object.keys(ethErr.data)
    return iface.decodeErrorResult(
      'OffchainLookup',
      ethErr.data[keys[0]].return,
    )
  } catch (e) {
    return null
  }
}


contract(
  'ancon protocol - dag contract with trusted gateway integration',
  (accounts) => {
    let metadataTransferDagTrusted
    let xdvnft
    let daiContract

    // Initialize the contracts and make sure they exist
    before(async () => {
      ;({
        xdvnft,
        daiContract,
      } = await Bluebird.props({
        xdvnft: XDVNFT.deployed(),
        daiContract: DAI.deployed(),
      }))
    })

    describe('when requesting a transfer ownership', () => {
      it('should revert', async () => {
        try {
          await xdvnft.transferURI(xdvnft.address, 1)
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
        const toAddress = xdvnft.address
        const tokenId = "1"
        const iface = new ethers.utils.Interface(xdvnft.abi)
        const metadataCid = "baguqeeramhau4x5j5zihi6ffksrl7wv7qnfd2urmut2e6c7oqdbu4jj7zljq"
        try 
        {
          await xdvnft.transferURI(xdvnft.address, 1)
        } catch (e) {
          const response = parseEthError(e, iface)
          if (response) {
            const url = response.url;

            const fromAddress = accounts[0]
            const inputdata ={
              toAddress,
              tokenId,
              metadataCid: metadataCid,
              fromOwner: accounts[1],
              toOwner: accounts[0],
              prefix: response.prefix,
            };

            const body = {
              jsonrpc: '2.0',
              method: 'durin_call',
              params: [
                xdvnft.address,
                fromAddress,
                inputdata,
                [iface.getFunction('transferURIWithProof')],
              ],
              id: 1,
            }

            const result = await fetchJson(
              {
                url,
              },
              JSON.stringify(body),
            )
            console.log(result)
            const resParse = JSON.parse(Web3.utils.hexToString(result.result))
            
            let txDataDecoded = ethers.utils.defaultAbiCoder.decode(['bytes', 'bytes', 'bytes', 'bytes', 'bytes', 'bytes', 'bytes'], ethers.utils.arrayify(resParse.txdata))
            
            console.log("RESULT CID DECODED", resParse.txdata )
            
            const resTransfer = await xdvnft.transferURIWithProof(toAddress, tokenId, ethers.utils.arrayify(resParse.txdata),{
              from: accounts[0]
            })
            console.log("TRANSFER RESULT", JSON.stringify(resTransfer))
            
            
            console.log(result)
          } else {
            console.log({ e })
          }
        }
      })
    })
  },
)
