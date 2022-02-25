const fs = require('fs')
const ContractImportBuilder = require('../contract-import-builder')
const KYX = artifacts.require('KYX')
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
const { base64, hexlify, keccak256, toUtf8Bytes } = require('ethers/lib/utils')

module.exports = async (deployer, network, accounts) => {
  const builder = new ContractImportBuilder()
  const path = `${__dirname}/../abi-export/kyx.js`

  builder.setOutput(path)
  builder.onWrite = (output) => {
    fs.writeFileSync(path, output)
  }

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
  } else if (network == 'bsc') {
    chainId = 56
    token = '0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d' //USDC
  }

  await deployer.deploy(KYX, token, chainId)
  const verifier = await KYX.deployed()

  //  await wxdv.enrollNFT(nft.address)
  // builder.addContract('XDVNFT', nft, nft.address, network)
  // builder.addContract('WXDV', wxdv, wxdv.address, network)
  builder.addContract('AnconProtocol', verifier, verifier.address, network)
}
