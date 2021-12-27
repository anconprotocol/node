const fs = require('fs')
const ContractImportBuilder = require('../contract-import-builder');

const Bytes = artifacts.require("Bytes");
const Memory = artifacts.require("Memory");
const ICS23 = artifacts.require("ICS23");
const AnconProtocol = artifacts.require("AnconProtocol");

module.exports = async (deployer, network, accounts) => {
  const builder = new ContractImportBuilder();
  const path = `${__dirname}/../abi-export/ics23.js`;

  builder.setOutput(path);
  builder.onWrite = (output) => {
    fs.writeFileSync(path, output);
  };
  

  await deployer.deploy(Memory);
  await deployer.link(Memory, Bytes);
  
  await deployer.deploy(Bytes);
  await deployer.link(Bytes, ICS23, AnconProtocol);

  
  await deployer.deploy(AnconProtocol, accounts[0]);
  const verifier = await AnconProtocol.deployed();



  
  builder.addContract(
    'AnconProtocol',
    verifier,
    verifier.address,
    network
  );
};
