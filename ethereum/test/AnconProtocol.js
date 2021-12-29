global['fetch'] = import('node-fetch')
// const ethers = require("ethers");
const Bluebird = require('bluebird')
const { ethers } = require('ethers')
const { base64, hexlify } = require('ethers/lib/utils')
require('dotenv').config({ path: '../.env' })
const Web3 = require('web3')
const AnconProtocol = artifacts.require('AnconProtocol')
// const ONMETA = artifacts.require("OnchainMetadata");

const { ANCON, PKEY } = process.env
require('ethers')
const {
  AnconProtocol__factory,
} = require('../types/lib/factories/AnconProtocol__factory')

// const signer = new ethers.Wallet(PKEY)
// const SENDER_ADDRESS = signer.address

const RPC_HOST = 'http://localhost:8545/'
const CONTRACT_ADDRESS = '0x77C51E844495899727dB63221af46220b0b13B37'
const BLOCK_NUMBER = 13730326
const ACCOUNT_ADDRESS = '0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6'

let proofCombined = [
  {
    exist: {
      valid: true,
      key:
        'L2FuY29ucHJvdG9jb2wvZTA1NDZjMDZlNDZlYWIzNDcyMmVhMTNjNTAyNGNiMDBmYjEzNmVmZDg3OGY0NThiNTViMDQ3YzhkOGU4Y2JiNi91c2VyL2JhZ3VxZWVyYWVocGRhN3pwcmJ1bXhoZzVncWlmbHkzYm1kbmFlb2NxbzRmbHZub2ZzaXB0ZmVoa3F5d3E=',
      value: 'ZGlkOndlYjppcGZzOnVzZXI6dGVzdA==',
      leaf: {
        valid: true,
        hash: 1,
        length: 1,
        prefix: 'AAIi',
        prehash_value: 1,
      },
      path: [
        {
          valid: true,
          hash: 1,
          prefix: 'AgQiIA==',
          suffix: 'IEga3tvbzXYCsg0xSYPX36OPlz1bwOPxMp239NZ9XwG+',
        },
        {
          valid: true,
          hash: 1,
          prefix: 'BAgiIKLI0YDL04HqPwDpHF8ht5pGCiq+Uw/0DTzSMDp7pE6mIA==',
        },
        {
          valid: true,
          hash: 1,
          prefix: 'BgwiIADQ5R72VKesArREiAa1kDeAFiOgt9tZ+32NLHPOPjP9IA==',
        },
        {
          valid: true,
          hash: 1,
          prefix: 'CBIiIA==',
          suffix: 'IPhUy0NkQUD/fk42UPtKzT0QWd1NDgJshHnLRSm3j7XQ',
        },
        {
          valid: true,
          hash: 1,
          prefix: 'CiAiIA==',
          suffix: 'IMo1qcvqM7Duwq5Ac3wtuoTitJjOTDVrB92+pkPPAlzW',
        },
      ],
    },
  },
]

function toABIproofs() {
  let z = { ...proofCombined[0].exist }
        z.key = hexlify(base64.decode(z.key))
        z.value = hexlify(base64.decode(z.value))
        z.leaf.prefix = hexlify(base64.decode(z.leaf.prefix))
        z.leaf.hash = 1
        z.path = z.path.map((x) => {
          let suffix
          if (!!x.suffix) {
            suffix = hexlify(base64.decode(x.suffix))
            return {
              valid: true,
              prefix: hexlify(base64.decode(x.prefix)),
              suffix: suffix,
              hash: 1,
            }
          } else {
            return {
              valid: true,
              prefix: hexlify(base64.decode(x.prefix)),
              hash: 1,
              suffix: '0x',
            }
          }
        })
        z.leaf.prehash_key = 0
        z.leaf.len = z.leaf.length

  return z
}

contract('QueryRootCalculation', (accounts) => {
  describe('when requesting to add metadata onchain', () => {
    it('should return true and emit event', async () => {
      try {
        const provider = new ethers.providers.Web3Provider(web3.currentProvider)
        const contract = await AnconProtocol.deployed()
        const contract2 = AnconProtocol__factory.connect(
          contract.address,
          provider,
        )
        
        z = toABIproofs();
        console.log(z)
        const resRootCalc = await contract2.callStatic.queryRootCalculation({
          ...z,
        })

        const restUpdtHeader = await contract.updateProtocolHeader(
          resRootCalc,
          { from: accounts[0] },
        )
        console.log(resRootCalc)
        const k = await contract2.verifyProof({ ...z })
        const restL2 = await contract2.enrollL2Account(z.key, z.value, z)
        assert.equal(k, true)
      } catch (e) {
        console.log(e)
      }
    })
  })
})