global['fetch'] = import('node-fetch')
// const ethers = require("ethers");
const Bluebird = require('bluebird')
require('dotenv').config({ path: '../.env' })
const Web3 = require('web3')

// const ONMETA = artifacts.require("OnchainMetadata");

const { ANCON, PKEY } = process.env
// const { fetchJson } = require("@ethersproject/web");

// const provider = new ethers.providers.JsonRpcProvider(ANCON);
const metadata = require('../build/contracts/OnchainMetadata.json')
const { formatBytes32String } = require('@ethersproject/strings')
// const signer = new ethers.Wallet(PKEY)
// const SENDER_ADDRESS = signer.address
console.log({
  ANCON,
  PKEY,
})

contract('Onchain Metadata', (accounts) => {
  let contract
  before(async () => {
    console.log('Metadata', metadata)
    const ONCHAINMET = new web3.eth.Contract(metadata.abi)
    let onChainMet
    ;({ onChainMet } = await Bluebird.props({
      onChainMet: ONCHAINMET.deploy({
        data: metadata.bytecode,
        arguments: [],
      }),
    }))

    contract = await onChainMet.send({
      from: '0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6',
      gas: 1500000,
      gasPrice: '30000000000',
    })
    console.log(
      'Contract deployed at address: ',
      contract.options.address,
    )
  })

  describe('when requesting to add a dagjson block', () => {
    it('should return true and emit event', async () => {

      const name = "metadata sample"
      const owner = "0x2a3D91a8D48C2892b391122b6c3665f52bCace23p"
      const description = "testing sample"
      const image = "baguqeeraouhd5jr7ktftgkbo5ufmihs3e2yqonajjejhomvbtsrwynlqdxba"

      const toBytes = Web3.default.utils.fromUtf8
      try {
        const res = await contract.encodeDagjsonBlock(
"/","=="
        )
        console.log(res)
      } catch (e) {
        console.log(e)
      }
    })
  })
})
