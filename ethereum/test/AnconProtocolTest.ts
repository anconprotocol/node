// const {utils, ethers} = require("ethers")
// const {AncProt__Factory} = require("../types/ethers-contracts/factories/AnconProtocol__factory.ts");
import { utils, ethers } from "ethers";
import { ExistenceProofStruct } from "../types/ethers-contracts/AnconProtocol";
import { AnconProtocol__factory } from "../types/ethers-contracts/factories/AnconProtocol__factory";
import { ics23 } from "@confio/ics23";
import { arrayify, base64 } from "ethers/lib/utils";

const RPC_HOST = "http://localhost:8545/";
const CONTRACT_ADDRESS = "0x9cb049eB339C2cBFdf182455aa6DE7566a1C5C4D";
const BLOCK_NUMBER = 13730326;
const ACCOUNT_ADDRESS = "0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6";

let proofCombined = [
  {
    exist: {
      key: "YW5jb25iYWZ5cmVpY2Q0Z3ZwYmpwNHdlY3RiejQyNzRhZ3Zmc3gybjN3cnhxZHp1cDQzZXVrM3R1ZDI3Y3htZQ==",
      value:
        "qmNkaWRgZGtpbmRobWV0YWRhdGFkbmFtZWp0ZW5kZXJtaW50ZWltYWdldWh0dHA6Ly9sb2NhbGhvc3Q6MTMxN2VsaW5rc4HYKlglAAFxEiAe9b351lHuEQAZl2qWbdxrvaLk6O7sh5TzqMDBMwQ3EmVvd25lcngzZGlkOmV0aHI6MHhlZUM1OEU4OTk5NjQ5NjY0MGM4YjU4OThBN2UwMjE4RTliNkU5MGNCZnBhcmVudGBnc291cmNlc4F4LlFtUVJXR3hvRGhnOEpNb0RaYU4xdDFLaVlDYzVkTTlOUFhZRVk3VzhrSjhKd2prZGVzY3JpcHRpb25qdGVuZGVybWludHV2ZXJpZmllZENyZWRlbnRpYWxSZWZg",
      leaf: {
        hash: "SHA256",
        prehash_key: "NO_HASH",
        prehash_value: "SHA256",
        len: "VAR_PROTO",
        prefix: "AAKSAQ==",
        valid: true,
      },
      path: [
        {
          hash: "SHA256",
          prefix: "AgSSASA=",
          suffix: "INdSgAXFKAv2D5wqrrbM+uSs2ynW0VuytR2UdOMuNagz",
        },
      ],
    },
  },
];

let exProof: ExistenceProofStruct;
exProof = [] as any;

async function main() {
  const provider = new ethers.providers.JsonRpcProvider(RPC_HOST);
  const contract = AnconProtocol__factory.connect(CONTRACT_ADDRESS, provider);

  exProof.key = proofCombined[0].exist.key;
  exProof.leaf = proofCombined[0].exist.leaf;
  exProof.path = [
    {
      hash: proofCombined[0].exist.path[0].hash,
      prefix: proofCombined[0].exist.path[0].prefix,
      suffix: proofCombined[0].exist.path[0].suffix,
      valid: true,
    },
  ]; //proofCombined[0].exist.path

  exProof.valid = true;
  exProof.value = proofCombined[0].exist.value;
  // await contract.updateProtocolHeader()
  const resRootCalc = await contract.queryRootCalculation(
    base64.decode(exProof.leaf.prefix as any),
    base64.decode(exProof.path[0].prefix as any),
    base64.decode(exProof.path[0].suffix as any),
    base64.decode(exProof.key),
    base64.decode(exProof.value),
  )
  console.log(resRootCalc)
 // await contract.updateProtocolHeader(resRootCalc)
  const resVerifyProof = await contract.verifyProof(
    base64.decode(exProof.key),
    base64.decode(exProof.value),
    base64.decode(exProof.leaf.prefix as any),
    base64.decode(exProof.path[0].prefix as any),
    base64.decode(exProof.path[0].suffix as any),
    arrayify(resRootCalc)
  );

  const contractName = "ANCON PROTOCOL";
  console.log(`Our ${contractName} root calculation is: ${resRootCalc}`);
  console.log(`Our ${contractName} proof is: ${resVerifyProof}`);

  console.log(`Listing Transfer events for block ${BLOCK_NUMBER}`);

  // const eventsFilter = contract.filters.Transfer();
  // const events = await contract.queryFilter(
  //   eventsFilter,
  //   BLOCK_NUMBER,
  //   BLOCK_NUMBER
  // );

  // for (const event of events) {
  //   console.log(
  //     `${event.args.src} -> ${event.args.dst} | ${utils.formatEther(
  //       event.args.wad
  //     )} ${contractName}`
  //   );
  // }
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
