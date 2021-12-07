const fs = require("fs");
var truffleContract = require("@truffle/contract");
// const Web3 = require('web3')

const metadata = require("../build/contracts/OnchainMetadata.json");
module.exports = async function (deployer) {
  try {
    console.log(metadata);
    let contract = new web3.eth.Contract(metadata.abi);

    contract = truffleContract(
      contract.deploy({
        data: metadata.bytecode,
        arguments: [],
      })
    );

    contractD = await deployer.deploy(contract);
    // const newContractInstance = await contract.send({
    //     from: '0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6',
    //     gas: 1500000,
    //     gasPrice: '30000000000'
    // })
    // console.log('Contract deployed at address: ', newContractInstance.options.address)
  } catch (e) {
    console.log(e);
  }
};
