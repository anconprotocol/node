
const BigNumber = require('bignumber.js')
const fs = require('fs')
const XDVNFT = artifacts.require('XDVNFT')
const DAI = artifacts.require('DAI')

const AguaclaraRouter = artifacts.require("AguaclaraRouter");

const Web3 = require("web3");
const ContractImportBuilder = require('../contract-import-builder');

const Proofs = artifacts.require("IBCExistenceProof");
const Bytes = artifacts.require("Bytes");
const Memory = artifacts.require("Memory");
const ICS23 = artifacts.require("ICS23");
const AnconVerifier = artifacts.require("AnconVerifier");
const Faucet = artifacts.require("Faucet");

module.exports = async (deployer, network, accounts) => {
  const builder = new ContractImportBuilder();
  const path = `${__dirname}/../abi-export/xdv.js`;

  builder.setOutput(path);
  builder.onWrite = (output) => {
    fs.writeFileSync(path, output);
  };
  let xdvnft;
  let stableCoinAddress = "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d";
  let dai;

  if (network === "bsc") {
    stableCoinAddress = "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d";
  }
  else {
    await deployer.deploy(DAI);
    dai = await DAI.deployed();
    stableCoinAddress = dai.address;
  }


  await deployer.deploy(Memory);
  await deployer.link(Memory, Bytes);
  await deployer.deploy(Bytes);
  await deployer.link(Bytes, AnconVerifier);
  await deployer.deploy(AnconVerifier, accounts[0]);
  await deployer.deploy(Faucet, dai.address);
  await Faucet.deployed()
  const verifier = await AnconVerifier.deployed();

  await deployer.deploy(AguaclaraRouter,  accounts[0], verifier.address);

  const r = await AguaclaraRouter.deployed();
  

  // Deploy NFT
  await deployer.deploy(XDVNFT, "XDVNFT", "XDVNFT", stableCoinAddress, verifier.address);

  xdvnft = await XDVNFT.deployed();
  await xdvnft.setServiceFeeForContract(new BigNumber(1 * 1e18));
  const fee_bn = new BigNumber(5 * 1e18);

  // Mint
  await dai.mint(accounts[0], fee_bn);
  await dai.approve(xdvnft.address, fee_bn);
  await xdvnft.mint(accounts[0], "bafyreicztwstn4ujtsnabjabn3hj7mvbhsgrvefbh37ddnx4w2pvghvsfm")

  // Verify
  const xp = '0x0abc020a40616e636f6e6261667972656961706b6a65333237657270347537333662376a757532376a6433697375697a75706d6137373563707736703479766737706c796912ea01a86364696460646b696e64686d65746164617461646e616d656a74656e6465726d696e7465696d61676575687474703a2f2f6c6f63616c686f73743a31333137656c696e6b7381d82a582500017112201ef5bdf9d651ee110019976a966ddc6bbda2e4e8eeec8794f3a8c0c133043712656f776e657278296469643a6b65793a7a386d57614a48586965415678784c6167427064614e574645424b56576d4d694567736f757263657381d82a582500017112201ef5bdf9d651ee110019976a966ddc6bbda2e4e8eeec8794f3a8c0c1330437126b6465736372697074696f6e6a74656e6465726d696e741a0b0801180120012a0300020c'
  const root = '0x6d813d0dab056200a230d3791bc3b5f9ba49930e351e56fcc7ea21ab4b1da52b'
  const keyPath = Web3.utils.toHex('anconbafyreiapkje327erp4u736b7juu27jd3isuizupma775cpw6p4yvg7plyi')
  const value = '0xa86364696460646b696e64686d65746164617461646e616d656a74656e6465726d696e7465696d61676575687474703a2f2f6c6f63616c686f73743a31333137656c696e6b7381d82a582500017112201ef5bdf9d651ee110019976a966ddc6bbda2e4e8eeec8794f3a8c0c133043712656f776e657278296469643a6b65793a7a386d57614a48586965415678784c6167427064614e574645424b56576d4d694567736f757263657381d82a582500017112201ef5bdf9d651ee110019976a966ddc6bbda2e4e8eeec8794f3a8c0c1330437126b6465736372697074696f6e6a74656e6465726d696e74'

  // await verifier.changeOwnerWithProof(
  //   xp,
  //   root,
  //   keyPath,
  //   value
  // );
  builder.addContract(
    'XDVNFT',
    xdvnft,
    xdvnft.address,
    network
  );

  builder.addContract(
    'AnconMetadataOwnableBridge',
    verifier,
    verifier.address,
    network
  );
};
