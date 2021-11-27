const path = require("path");
const {
  init,
  emulator,
  getContractAddress,
  executeScript,
  getAccountAddress,
  deployContractByName,
  mintFlow,
  shallResolve,
} = require("flow-js-testing");
const ics23 = require("@confio/ics23");
const ethers = require("ethers");
const { assert } = require("chai");

function toABI({ exist }) {
  const innerOp = [
    {
      valid: true,
      hash: ics23.ics23.HashOp[exist.path[0].hash],
      prefix: Array.from(ethers.utils.base64.decode(exist.path[0].prefix)),
      suffix: Array.from(ethers.utils.base64.decode(exist.path[0].suffix)),
    },
  ];

  const leafOp = {
    valid: true,
    hash: ics23.ics23.HashOp[exist.leaf.hash],
    prehash_key: ics23.ics23.HashOp[exist.leaf.prehash_key],
    prehash_value: ics23.ics23.HashOp[exist.leaf.prehash_value],
    len: ics23.ics23.LengthOp[exist.leaf.length],
    prefix: Array.from(ethers.utils.base64.decode(exist.leaf.prefix)),
  };

  return {
    leafOp,
    innerOp,
    prefix: Array.from(ethers.utils.base64.decode(exist.leaf.prefix)),
    innerOpHash: ics23.ics23.HashOp[exist.path[0].hash],
    key: Array.from(ethers.utils.base64.decode(exist.key)),
    value: Array.from(ethers.utils.base64.decode(exist.value)),
  };
}

const proofCombined = [
  {
    exist: {
      key: "YW5jb25iYWZ5cmVpY2Q0Z3ZwYmpwNHdlY3RiejQyNzRhZ3Zmc3gybjN3cnhxZHp1cDQzZXVrM3R1ZDI3Y3htZQ==",
      value:
        "qmNkaWRgZGtpbmRobWV0YWRhdGFkbmFtZWp0ZW5kZXJtaW50ZWltYWdldWh0dHA6Ly9sb2NhbGhvc3Q6MTMxN2VsaW5rc4HYKlglAAFxEiAe9b351lHuEQAZl2qWbdxrvaLk6O7sh5TzqMDBMwQ3EmVvd25lcngzZGlkOmV0aHI6MHhlZUM1OEU4OTk5NjQ5NjY0MGM4YjU4OThBN2UwMjE4RTliNkU5MGNCZnBhcmVudGBnc291cmNlc4F4LlFtUVJXR3hvRGhnOEpNb0RaYU4xdDFLaVlDYzVkTTlOUFhZRVk3VzhrSjhKd2prZGVzY3JpcHRpb25qdGVuZGVybWludHV2ZXJpZmllZENyZWRlbnRpYWxSZWZg",
      leaf: {
        hash: "SHA256",
        prehash_key: "NO_HASH",
        prehash_value: "SHA256",
        length: "VAR_PROTO",
        prefix: "AAKSAQ==",
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

const PRECALCULATED_MERKLE_ROOT =
  "0x16dbcad17f5eff1b8fc04ea7527023811794839765764a6e62e41bb79ebac2cc";

describe("AnconVerifier", () => {
  const contractName = `AnconVerifier`;
  let accountAddress;

  beforeAll(async () => {
    const basePath = path.resolve(__dirname, "..");
    const port = 8080;

    await init(basePath, { port });
    await emulator.start(port, false);

    // We need to have FLOW to be able to store the contract:
    // Source: https://docs.onflow.org/concepts/storage/#storage-capacity
    accountAddress = await getAccountAddress("emulator-account");
    await mintFlow(accountAddress, "1000.0");

    await deployContractByName({ to: accountAddress, name: contractName});
  });

  afterAll(async () => {
    return emulator.stop();
  });

  it("should have deployed correctly", async () => {
    const contractAddress = await getContractAddress(contractName);
    assert.equal(contractAddress, accountAddress);
  });

  it("should calculate manually ICS23 Merkle Proofs", async () => {
    const code = `
      import AnconVerifier from ${accountAddress}

      pub fun main(
        key: [UInt8],
        value: [UInt8],
        leafPrefix: [UInt8],
        pathPrefix: [UInt8],
        pathSuffix: [UInt8],
      ): [UInt8] {
        let leafStruct = AnconVerifier.LeafOp(
          valid: true,
          hash: AnconVerifier.HashOp.SHA256,
          prehash_key: AnconVerifier.HashOp.NO_HASH,
          prehash_value: AnconVerifier.HashOp.SHA256,
          len: AnconVerifier.LengthOp.VAR_PROTO,
          prefix: leafPrefix
        )

        let pathStruct = AnconVerifier.InnerOp(
          valid: true,
          hash: AnconVerifier.HashOp.SHA256,
          prefix: pathPrefix,
          suffix: pathSuffix,
        )

        let p = AnconVerifier.ExistenceProof(
          valid: true,
          key: key,
          value: value,
          leaf: leafStruct,
          path: [pathStruct]
        )
        return AnconVerifier.calculate(p:p)
      }
    `;

    const proof = toABI(proofCombined[0]);

    const args = [
      proof.key,
      proof.value,
      proof.leafOp.prefix,
      proof.innerOp[0].prefix,
      proof.innerOp[0].suffix,
    ];
    const result = await shallResolve(executeScript({ code, args }));
    assert.equal(ethers.utils.hexlify(result), PRECALCULATED_MERKLE_ROOT);
  });

  it("should correctly verify ownership via ICS23", async () => {
    const code = `
      import AnconVerifier from ${accountAddress}

      pub fun main(
        leafOp: [UInt8],
        prefix: [UInt8],
        existenceProofInnerOp: [[[UInt8]]],
        existenceProofInnerOpHash: UInt8,
        existenceProofKey: [UInt8],
        existenceProofValue: [UInt8],
      ): AnyStruct {
        let root = AnconVerifier.queryRootCalculation(
          leafOpUint: leafOp,
          prefix: prefix,
          existenceProofInnerOp: existenceProofInnerOp,
          existenceProofInnerOpHash: existenceProofInnerOpHash,
          existenceProofKey: existenceProofKey,
          existenceProofValue: existenceProofValue,
        );

        return AnconVerifier.changeOwnerWithProof(
          leafOpUint: leafOp,
          prefix: prefix,
          existenceProofInnerOp: existenceProofInnerOp,
          existenceProofInnerOpHash: existenceProofInnerOpHash,
          existenceProofKey: existenceProofKey,
          existenceProofValue: existenceProofValue,
          root: root,
          key: existenceProofKey,
          value: existenceProofValue,
        );
      }
    `;

    const proof = toABI(proofCombined[0]);
    const leafOps = [
      proof.leafOp.hash,
      proof.leafOp.prehash_key,
      proof.leafOp.prehash_value,
      proof.leafOp.len,
    ];
    const mappedInnerOp = proof.innerOp.map((io) => [io.prefix, io.suffix]);
    const args = [
      leafOps,
      proof.prefix,
      mappedInnerOp,
      proof.innerOpHash,
      proof.key,
      proof.value,
    ];

    const result = await shallResolve(executeScript({ code, args }));
    assert.isTrue(result, "Must have Validated");
  });
});
 