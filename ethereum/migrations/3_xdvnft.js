const fs = require("fs");
const ContractImportBuilder = require("../contract-import-builder");

const { ethers } = require("ethers");

const XDVNFT = artifacts.require("XDVNFT");

module.exports = async (deployer, network, accounts) => {
  const builder = new ContractImportBuilder();
  const path = `${__dirname}/../abi-export/ics23.js`;

  builder.setOutput(path);
  builder.onWrite = (output) => {
    fs.writeFileSync(path, output);
  };

  // dai
  await deployer.deploy(XDVNFT, "XDVNFT", "XDVNFT", "0xec5dcb5dbf4b114c9d0f65bccab49ec54f6a0867","0xBE51ED643231c318c355f12495BD9B220209090A");
  const c = await XDVNFT.deployed();


  builder.addContract("XDVNFT", c, c.address, network);
};

