
const BigNumber = require('bignumber.js')
const fs = require('fs')
const Web3 = require("web3");
const ContractImportBuilder = require('../contract-import-builder');

const XDVNFT = artifacts.require('XDVNFT')
const DAI = artifacts.require('DAI')
const TrustedOffchainHelper = artifacts.require("TrustedOffchainHelper");

const ECDSA = artifacts.require("ECDSA");
const Address = artifacts.require("Address");

module.exports = async (deployer, network, accounts) => {
  const builder = new ContractImportBuilder();
  const path = `${__dirname}/../abi-export/xdv.js`;

  builder.setOutput(path);
  builder.onWrite = (output) => {
    fs.writeFileSync(path, output);
  };

  let stableCoinAddress = process.env.STABLE_COIN_ADDRESS;


  // DAG Contract
  await deployer.deploy(Address);
  await deployer.link(Address, XDVNFT);
  await deployer.deploy(ECDSA);
  await deployer.link(ECDSA, XDVNFT);

  // ERC20 - stablecoin
  await deployer.deploy(DAI);
  const dai = await DAI.deployed();
  stableCoinAddress = stableCoinAddress || dai.address;
 

  // NFT - ERC721
  await deployer.deploy(XDVNFT, "XDVNFT", "XDVNFT", stableCoinAddress );

  const xdvnft = await XDVNFT.deployed();
  await xdvnft.setServiceFeeForContract(new BigNumber(1 * 1e18));
  const fee_bn = new BigNumber(5 * 1e18);

  // Mint
  await dai.mint(accounts[0], fee_bn);
  await dai.approve(xdvnft.address, fee_bn);
  await xdvnft.mint(accounts[0], "bafyreicztwstn4ujtsnabjabn3hj7mvbhsgrvefbh37ddnx4w2pvghvsfm")

  // Configure DAG contract
  // Set Durin Gateway endpoint
  await xdvnft.setUrl("http://localhost:7788/rpc");

  // Set Trusted Signer Address
  await xdvnft.setSigner("0x2a3D91a8D48C2892b391122b6c3665f52bCace23");

  builder.addContract(
    'XDVNFT',
    xdvnft,
    xdvnft.address,
    network
  );

  builder.addContract(
    'DAI',
    dai,
    dai.address,
    network
  );
};
