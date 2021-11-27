const { BigNumber } = require('bignumber.js')
const { assert } = require('chai')
const Web3 = require('web3')
const ethers = require('ethers')
const Bluebird = require('bluebird')
//const TestUSDC = artifacts.require("USDC");
const AnconVerifier = artifacts.require('AnconVerifier')
//const XDVController = artifacts.require("XDVController");
const ics23 = require('@confio/ics23')
const { base64 } = require('ethers/lib/utils')
const Faucet = artifacts.require('Faucet')
const Dai = artifacts.require('DAI')

contract('Ancon - ICS23 Javascript to ABI', (accounts) => {
  // let erc20Contract;
  // let controllerContract;
  // let documentMinterAddress;
  let anconVerifierContract
  let faucetContract
  let daiContract

  let toABI = ({ exist }) => {
    const innerOp = []
    innerOp.push([
      ethers.utils.base64.decode(exist.path[0].prefix),
      ethers.utils.base64.decode(exist.path[0].suffix),
    ])
    const leafOp = [
      ics23.ics23.HashOp[exist.leaf.hash],
      ics23.ics23.HashOp[exist.leaf.prehash_key],
      ics23.ics23.HashOp[exist.leaf.prehash_value],
      ics23.ics23.LengthOp[exist.leaf.length],
    ]

    return {
      leafOp,
      prefix: ethers.utils.base64.decode(exist.leaf.prefix),
      innerOp,
      innerOpHash: ics23.ics23.HashOp[exist.path[0].hash],
      key: exist.key,
      value: exist.value,
    }
  }

  let proofs = [
    {
      exist: {
        key:
          'YW5jb25iYWZ5cmVpY2Q0Z3ZwYmpwNHdlY3RiejQyNzRhZ3Zmc3gybjN3cnhxZHp1cDQzZXVrM3R1ZDI3Y3htZQ==',
        value:
          'qmNkaWRgZGtpbmRobWV0YWRhdGFkbmFtZWp0ZW5kZXJtaW50ZWltYWdldWh0dHA6Ly9sb2NhbGhvc3Q6MTMxN2VsaW5rc4HYKlglAAFxEiAe9b351lHuEQAZl2qWbdxrvaLk6O7sh5TzqMDBMwQ3EmVvd25lcngzZGlkOmV0aHI6MHhlZUM1OEU4OTk5NjQ5NjY0MGM4YjU4OThBN2UwMjE4RTliNkU5MGNCZnBhcmVudGBnc291cmNlc4F4LlFtUVJXR3hvRGhnOEpNb0RaYU4xdDFLaVlDYzVkTTlOUFhZRVk3VzhrSjhKd2prZGVzY3JpcHRpb25qdGVuZGVybWludHV2ZXJpZmllZENyZWRlbnRpYWxSZWZg',
        leaf: {
          hash: 'SHA256',
          prehash_key: 'NO_HASH',
          prehash_value: 'SHA256',
          length: 'VAR_PROTO',
          prefix: 'AALEAg==',
        },
        path: [
          {
            hash: 'SHA256',
            prefix: 'AgTEAiA=',
            suffix: 'INIPC/ialnKao6gMD70kc1cjFB4KTNRMaxU5n+KxvbAV',
          },
        ],
      },
    },
  ]
  let proofCombined = [
    {
      exist: {
        key:
          'YW5jb25iYWZ5cmVpY2Q0Z3ZwYmpwNHdlY3RiejQyNzRhZ3Zmc3gybjN3cnhxZHp1cDQzZXVrM3R1ZDI3Y3htZQ==',
        value:
          'qmNkaWRgZGtpbmRobWV0YWRhdGFkbmFtZWp0ZW5kZXJtaW50ZWltYWdldWh0dHA6Ly9sb2NhbGhvc3Q6MTMxN2VsaW5rc4HYKlglAAFxEiAe9b351lHuEQAZl2qWbdxrvaLk6O7sh5TzqMDBMwQ3EmVvd25lcngzZGlkOmV0aHI6MHhlZUM1OEU4OTk5NjQ5NjY0MGM4YjU4OThBN2UwMjE4RTliNkU5MGNCZnBhcmVudGBnc291cmNlc4F4LlFtUVJXR3hvRGhnOEpNb0RaYU4xdDFLaVlDYzVkTTlOUFhZRVk3VzhrSjhKd2prZGVzY3JpcHRpb25qdGVuZGVybWludHV2ZXJpZmllZENyZWRlbnRpYWxSZWZg',
        leaf: {
          hash: 'SHA256',
          prehash_key: 'NO_HASH',
          prehash_value: 'SHA256',
          length: 'VAR_PROTO',
          prefix: 'AAKSAQ==',
        },
        path: [
          {
            hash: 'SHA256',
            prefix: 'AgSSASA=',
            suffix: 'INdSgAXFKAv2D5wqrrbM+uSs2ynW0VuytR2UdOMuNagz',
          },
        ],
      },
    },
  ]

  // Initialize the contracts and make sure they exist
  before(async () => {
    ; ({ anconVerifierContract, faucetContract, daiContract } = await Bluebird.props({
      anconVerifierContract: AnconVerifier.deployed(),
      //erc20Contract: TestUSDC.deployed(),
      //con trollerContract: XDVController.deployed(),
      faucetContract: Faucet.deployed(),
      daiContract: Dai.deployed(),
    }))
  })

  describe('when requesting an ownership update to a NFT using an ics23 vector commitment proof', () => {
    it('should verify ownership ok', async () => {
      const abiProof = toABI(proofCombined[0])

      // const r = ics23.calculateExistenceRoot(
      // proofs[0].exist,
      // )
      // console.log(ethers.utils.hexlify(r))

      const root = await anconVerifierContract.requestRoot(
        abiProof.leafOp,
        abiProof.prefix,
        abiProof.innerOp,
        abiProof.innerOpHash,
        ethers.utils.base64.decode(abiProof.key),
        ethers.utils.base64.decode(abiProof.value),
        {
          from: accounts[0],
        },
      )
      console.log(root)
      const res = await anconVerifierContract.changeOwnerWithProof(
        abiProof.leafOp,
        abiProof.prefix,
        abiProof.innerOp,
        abiProof.innerOpHash,
        ethers.utils.base64.decode(abiProof.key),
        ethers.utils.base64.decode(abiProof.value),
        root,
        ethers.utils.base64.decode(abiProof.key),
        ethers.utils.base64.decode(abiProof.value),
      )
      //   function changeOwnerWithProof(
      //     bytes[] memory existenceProofLeafOp,
      //     bytes[][] memory existenceProofInnerOp,
      //     bytes memory existenceProofKey,
      //     bytes memory existenceProofValue,
      //     bytes memory rootBz,
      //     bytes memory pathBz,
      //     bytes memory value
      // ) public returns (bool) {

      assert.equal(res, true)
    })
  })

  describe('Requesting funds from faucet', () => {
    it('should recieve tokens from faucet', async () => {

      await daiContract.mint(faucetContract.address, BigInt(1e18));

      await faucetContract.getTokenFromFaucet(accounts[1]);
      const res = await daiContract.balanceOf(accounts[1]);
      assert.equal(res, BigInt(1e18))
    })
  })
})
