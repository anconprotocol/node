const fs = require('fs')
const ContractImportBuilder = require('../contract-import-builder')
const AnconToken = artifacts.require('AnconToken')
const { deployProxy, upgradeProxy } = require('@openzeppelin/truffle-upgrades');

module.exports = async (deployer, network, accounts) => {
    const builder = new ContractImportBuilder()
    const path = `${__dirname}/../abi-export/token.js`

    builder.setOutput(path)
    builder.onWrite = (output) => {
        fs.writeFileSync(path, output)
    }

    await deployProxy(AnconToken, 
        [], { deployer });


    const token = await AnconToken.deployed();
    console.log('\tCreando stakeholders...')
    await token.addStakeHolders([accounts[1],accounts[2]], 
        ['100000000000000000000000',
            '300000000000000000000000'], 
        [1,2]);
    builder.addContract('AnconToken', token, token.address, network);
}