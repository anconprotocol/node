global["fetch"] = import("node-fetch");
require("dotenv").config({ path: "../.env" });
const Web3 = require("web3");

var assert = require("assert");
const OnChainMetaJSON = require("../build/contracts/OnchainMetadata.json");

// const { ANCON, PKEY } = process.env;
let OnchainContractInstance;

describe("Onchain Metadata Contract", (deployer) => {
  before(async () => {
    //Ganache local
    const web3 = new Web3(
      new Web3.providers.HttpProvider("http://localhost:8545")
    );
    console.log("Metadata", OnChainMetaJSON);
    const OnchainMetaContract = new web3.eth.Contract(OnChainMetaJSON.abi);

    return OnchainMetaContract.deploy({
      data: OnChainMetaJSON.bytecode,
      arguments: [],
    })
      .send({
        from: "0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6",
        gas: 1500000,
        gasPrice: "30000000000",
      })
      .then(function (newContractInstance) {
        OnchainContractInstance = newContractInstance;
        console.log(
          "\n Contract ADDRESS\n",
          newContractInstance.options.address,
        ); // instance with the new contract address
      });
    // console.log("\n\n\n\n TEEEEEESTT",OnchainContractInstance)
  });

  describe("when requesting to add metadata onchain", () => {
    it("should return true and emit event", async () => {
      const toBytes= Web3.utils.fromUtf8
      let sources = {
        "mylink" : "QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D"
      }
      try {
        const res = await OnchainContractInstance.methods.setOnchainMetadata(
          "",
          "",
          "",
          "",
          "",
          toBytes(JSON.stringify(sources))
        ).send({
          from: "0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6",
          gas: 1500000,
          gasPrice: "30000000000",
        });
        console.log("\n\n\n\n Result",res);
      } catch (e) {
        console.log(e);
      }
    });
  });
});
