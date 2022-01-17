const fs = require('fs')
const ContractImportBuilder = require('../contract-import-builder')
const XDVNFT = artifacts.require('XDVNFT')

const { ethers } = require('ethers')
const Bytes = artifacts.require('Bytes')
const Memory = artifacts.require('Memory')
const ICS23 = artifacts.require('ICS23')
const Ics23Helper = artifacts.require('Ics23Helper')

const AnconProtocol = artifacts.require('AnconProtocol')
const {
  AnconProtocol__factory,
} = require('../types/lib/factories/AnconProtocol__factory')
const { base64, hexlify } = require('ethers/lib/utils')

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

module.exports = async (deployer, network, accounts) => {
  const builder = new ContractImportBuilder()
  const path = `${__dirname}/../abi-export/ics23.js`

  builder.setOutput(path)
  builder.onWrite = (output) => {
    fs.writeFileSync(path, output)
  }

  if (network === 'kovan') {
    await deployer.deploy(
      XDVNFT,
      'XDVNFT',
      'XDVNFT',
      '0x4f96fe3b7a6cf9725f59d353f723c1bdb64ca6aa',
      '0xD5BA6f5eCfc80F13784446511B2be73f9d5e2707',
    )
  } else {
    await deployer.deploy(
      XDVNFT,
      'XDVNFT',
      'XDVNFT',
      '0xec5dcb5dbf4b114c9d0f65bccab49ec54f6a0867',
      '0xDB37c3D3316455d16D7788d47805e1458c8cdbBa',
    )
  }
  const c = await XDVNFT.deployed()

  await c.setServiceFeeForContract('5')
  await c.setServiceFeeForPaymentAddress('5')

  builder.addContract('XDVNFT', c, c.address, network)
}
