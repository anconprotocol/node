const fs = require('fs')
const ContractImportBuilder = require('../contract-import-builder')
const InstantRelayer = artifacts.require('InstantRelayer')

module.exports = async (deployer, network, accounts) => {
  const builder = new ContractImportBuilder()
  const path = `${__dirname}/../abi-export/instantrelayer.js`

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
  await deployer.deploy(
    InstantRelayer,
    token,
    '0x3A942779cBc73D5DA159DDcCe3FE9c1A16E5Fcba',
    chainId,
  )
  const instantrelayer = await InstantRelayer.deployed()

  //  await wxdv.enrollNFT(nft.address)
  // builder.addContract('XDVNFT', nft, nft.address, network)
  // builder.addContract('WXDV', wxdv, wxdv.address, network)
  builder.addContract('InstantRelayer', instantrelayer, instantrelayer.address, network)
}
