const fs = require('fs')
const ContractImportBuilder = require('../contract-import-builder')
const WXDV = artifacts.require('WXDV')
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

  await deployer.deploy(Memory)
  await deployer.link(Memory, Bytes)

  await deployer.deploy(Bytes)
  await deployer.link(Bytes, ICS23, Ics23Helper, AnconProtocol)

  let chainId = 97
  let token = '0xec5dcb5dbf4b114c9d0f65bccab49ec54f6a0867'

  if (network === 'bsctestnet') {
    // no op
  } else if (network === 'kovan') {
    chainId = 42
    token = '0x4f96fe3b7a6cf9725f59d353f723c1bdb64ca6aa'
  } else if (network === 'mumbai') {
    chainId = 80001
    token = '0x326c977e6efc84e512bb9c30f76e30c160ed06fb'
  } else if (network === 'gnosis') {
    chainId = 100
    token = '0xe91D153E0b41518A2Ce8Dd3D7944Fa863463a97d'
  } else if (network === 'auroratestnet') {
    chainId = 1313161555
    token = '0xc115851ca60aed2ccc6ee3d5343f590834e4a3ab'
  }
  await deployer.deploy(AnconProtocol, token, chainId)
  const verifier = await AnconProtocol.deployed()
  await deployer.deploy(WXDV, 'WXDV', 'WXDV', token, verifier.address, chainId)
  const wxdv = await WXDV.deployed()

  await verifier.setPaymentToken(token, { from: accounts[0] })
  await verifier.setWhitelistedDagGraph(
    web3.utils.keccak256('tensta'),
    '0x04cc4232356b66A112ED42E2c51b3B062b4c94C2',
    { from: accounts[0] },
  )
  await verifier.setWhitelistedDagGraph(
    web3.utils.keccak256('anconprotocol'),
    '0x28CB56Ef6C64B066E3FfD5a04E0214535732e57F',
    { from: accounts[0] },
  )

  await verifier.setAccountRegistrationFee('500000000', { from: accounts[0] })
  await verifier.setDagGraphFee('500000000', { from: accounts[0] })
  await verifier.setProtocolFee('500000000', { from: accounts[0] })

  await wxdv.setServiceFeeForContract('50000', { from: accounts[0] })
  await wxdv.setServiceFeeForPaymentAddress('50000', { from: accounts[0] })

  await deployer.deploy(
    XDVNFT,
    'XDVNFT',
    'XDVNFT',
    token,
    wxdv.address,
    chainId,
  )

  const nft = await XDVNFT.deployed()

  await wxdv.enrollNFT(nft.address)
  builder.addContract('XDVNFT', nft, nft.address, network)
  builder.addContract('WXDV', wxdv, wxdv.address, network)
  builder.addContract('AnconProtocol', verifier, verifier.address, network)
}
