global["fetch"] = import("node-fetch");
// const ethers = require("ethers");
const Bluebird = require("bluebird");
require("dotenv").config({ path: "../.env" });
const Web3 = require("web3");

// const ONMETA = artifacts.require("OnchainMetadata");

const { ANCON, PKEY } = process.env;
// const { fetchJson } = require("@ethersproject/web");

// const provider = new ethers.providers.JsonRpcProvider(ANCON);
const metadata = require("../build/contracts/OnchainMetadata.json");
const { formatBytes32String } = require("@ethersproject/strings");
// const signer = new ethers.Wallet(PKEY)
// const SENDER_ADDRESS = signer.address
console.log({
  ANCON,
  PKEY,
});

contract("Onchain Metadata", (accounts) => {
  
  before(async () => {
    console.log("Metadata", metadata);
    const ONCHAINMET = new web3.eth.Contract(metadata.abi);
    let onChainMet;
    ({} = await Bluebird.props({
      onChainMet: ONCHAINMET.deploy({
        data: metadata.bytecode,
        arguments: [],
      }),
    }));

    const onChainIntance = await onChainMet.send({
      from: "0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6",
      gas: 1500000,
      gasPrice: "30000000000",
    });
    console.log(
      "Contract deployed at address: ",
      onChainIntance.options.address
    );
  });

  describe("when requesting to add metadata onchain", () => {
    it("should return true and emit event", async () => {
      try {
        const res = await onChainIntance.setOnchainMetadata(
          "",
          "",
          "",
          "",
          "",
          formatBytes32String(" ")
        );
        console.log(res);
      } catch (e) {
        console.log(e);
      }
    });
  });
});
