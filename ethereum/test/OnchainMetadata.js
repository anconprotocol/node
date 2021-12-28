global["fetch"] = import("node-fetch");
// const ethers = require("ethers");
const Bluebird = require("bluebird");
const { ethers } = require("ethers");
require("dotenv").config({ path: "../.env" });
const Web3 = require("web3");
const AnconProtocol = artifacts.require("AnconProtocol");
// const ONMETA = artifacts.require("OnchainMetadata");

const { ANCON, PKEY } = process.env;
require("ethers");
const {
  AnconProtocol__factory,
} = require("../types/lib/factories/AnconProtocol__factory");

// const signer = new ethers.Wallet(PKEY)
// const SENDER_ADDRESS = signer.address

const RPC_HOST = "http://localhost:8545/";
const CONTRACT_ADDRESS = "0x77C51E844495899727dB63221af46220b0b13B37";
const BLOCK_NUMBER = 13730326;
const ACCOUNT_ADDRESS = "0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6";

let proofCombined = [
  {
    exist: {
      key: "L2FuY29ucHJvdG9jb2wvZTA1NDZjMDZlNDZlYWIzNDcyMmVhMTNjNTAyNGNiMDBmYjEzNmVmZDg3OGY0NThiNTViMDQ3YzhkOGU4Y2JiNi91c2VyL2JhZ3VxZWVyYWVocGRhN3pwcmJ1bXhoZzVncWlmbHkzYm1kbmFlb2NxbzRmbHZub2ZzaXB0ZmVoa3F5d3E=",
      leaf: {
        hash: 1,
        length: 1,
        prefix: "AAIi",
        prehash_value: 1,
      },
      path: [
        {
          hash: 1,
          prefix: "AgQiIA==",
          suffix: "IEga3tvbzXYCsg0xSYPX36OPlz1bwOPxMp239NZ9XwG+",
        },
        {
          hash: 1,
          prefix: "BAgiIKLI0YDL04HqPwDpHF8ht5pGCiq+Uw/0DTzSMDp7pE6mIA==",
          suffix: "",
        },
        {
          hash: 1,
          prefix: "BgwiIADQ5R72VKesArREiAa1kDeAFiOgt9tZ+32NLHPOPjP9IA==",
          suffix: "",
        },
        {
          hash: 1,
          prefix: "CBIiIA==",
          suffix: "IPhUy0NkQUD/fk42UPtKzT0QWd1NDgJshHnLRSm3j7XQ",
        },
        {
          hash: 1,
          prefix: "CiAiIA==",
          suffix: "IMo1qcvqM7Duwq5Ac3wtuoTitJjOTDVrB92+pkPPAlzW",
        },
      ],
      value: "ZGlkOndlYjppcGZzOnVzZXI6dGVzdA==",
    },
  },
];

function toABIproofs(proofCombined) {
  let innerOps = proofCombined[0].exist.path.map((p) => [
    web3.utils.toHex(ethers.utils.base64.decode(p.prefix)),
    web3.utils.toHex(ethers.utils.base64.decode(p.suffix)),
  ]);

  let key = proofCombined[0].exist.key;

  return {
    key: ethers.utils.base64.decode(key),
    value: ethers.utils.base64.decode(proofCombined[0].exist.value),
    prefix: ethers.utils.base64.decode(proofCombined[0].exist.leaf.prefix),
    innerOps: innerOps,
  };
}

contract("Onchain Metadata", (accounts) => {
  describe("when requesting to add metadata onchain", () => {
    it("should return true and emit event", async () => {
      try {
        // const provider = new ethers.providers.Web3Provider(web3.givenProvider);
        // const contract = AnconProtocol__factory.connect(
        //   CONTRACT_ADDRESS,
        //   provider
        // );
        const contract = await AnconProtocol.deployed();

        const abiProof = toABIproofs(proofCombined);
        const resRootCalc = await contract.queryRootCalculation(
          web3.utils.toHex(abiProof.prefix),
          web3.utils.toHex(abiProof.key),
          web3.utils.toHex(abiProof.value),
          abiProof.innerOps
        );

        const restUpdtHeader = await contract.updateProtocolHeader(
          resRootCalc,
          { from: accounts[0] }
        );
        console.log(restUpdtHeader);
        
      } catch (e) {
        console.log(e);
      }
    });
  });
});
