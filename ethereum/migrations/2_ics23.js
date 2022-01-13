const fs = require('fs')
const ContractImportBuilder = require('../contract-import-builder')
const XDVNFT = artifacts.require('XDVNFT')

const { ethers } = require('ethers')
const Bytes = artifacts.require('Bytes')
const Memory = artifacts.require('Memory')
const ICS23 = artifacts.require('ICS23')
const Ics23Helper = artifacts.require('Ics23Helper')

const AnconProtocol = artifacts.require('AnconProtocol')
const {
  AnconProtocol__factory,
} = require('../types/lib/factories/AnconProtocol__factory')
const { base64, hexlify } = require('ethers/lib/utils')

let proofCombined = [
  {
    exist: {
      valid: true,
      key:
        'L2FuY29ucHJvdG9jb2wvZTA1NDZjMDZlNDZlYWIzNDcyMmVhMTNjNTAyNGNiMDBmYjEzNmVmZDg3OGY0NThiNTViMDQ3YzhkOGU4Y2JiNi91c2VyL2JhZ3VxZWVyYWVocGRhN3pwcmJ1bXhoZzVncWlmbHkzYm1kbmFlb2NxbzRmbHZub2ZzaXB0ZmVoa3F5d3E=',
      value: 'ZGlkOndlYjppcGZzOnVzZXI6dGVzdA==',
      leaf: {
        valid: true,
        hash: 1,
        length: 1,
        prefix: 'AAIi',
        prehash_value: 1,
      },
      path: [
        {
          valid: true,
          hash: 1,
          prefix: 'AgQiIA==',
          suffix: 'IEga3tvbzXYCsg0xSYPX36OPlz1bwOPxMp239NZ9XwG+',
        },
        {
          valid: true,
          hash: 1,
          prefix: 'BAgiIKLI0YDL04HqPwDpHF8ht5pGCiq+Uw/0DTzSMDp7pE6mIA==',
        },
        {
          valid: true,
          hash: 1,
          prefix: 'BgwiIADQ5R72VKesArREiAa1kDeAFiOgt9tZ+32NLHPOPjP9IA==',
        },
        {
          valid: true,
          hash: 1,
          prefix: 'CBIiIA==',
          suffix: 'IPhUy0NkQUD/fk42UPtKzT0QWd1NDgJshHnLRSm3j7XQ',
        },
        {
          valid: true,
          hash: 1,
          prefix: 'CiAiIA==',
          suffix: 'IMo1qcvqM7Duwq5Ac3wtuoTitJjOTDVrB92+pkPPAlzW',
        },
      ],
    },
  },
]

module.exports = async (deployer, network, accounts) => {
  const builder = new ContractImportBuilder()
  const path = `${__dirname}/../abi-export/ics23.js`

  builder.setOutput(path)
  builder.onWrite = (output) => {
    fs.writeFileSync(path, output)
  }

  await deployer.deploy(Memory)
  await deployer.link(Memory, Bytes)

  await deployer.deploy(Bytes)
  await deployer.link(Bytes, ICS23, Ics23Helper, AnconProtocol)

  await deployer.deploy(AnconProtocol,  '0xec5dcb5dbf4b114c9d0f65bccab49ec54f6a0867')
  const verifier = await AnconProtocol.deployed('0xec5dcb5dbf4b114c9d0f65bccab49ec54f6a0867')
  await deployer.deploy(
    XDVNFT,
    'XDVNFT',
    'XDVNFT',
    '0xec5dcb5dbf4b114c9d0f65bccab49ec54f6a0867',
    verifier.address,
  )
  const c = await XDVNFT.deployed()

  await verifier.setPaymentToken('0xec5dcb5dbf4b114c9d0f65bccab49ec54f6a0867')
  await verifier.updateProtocolHeader('0x')
  // await verifier.setProtocolFee(new BigNumber(500000000))
  // await verifier.setAccountRegistrationFee(new BigNumber(500000000))

  builder.addContract('XDVNFT', c, c.address, network)

  //  const provider = new ethers.providers.Web3Provider(web3.currentProvider);
  //  const contract = await AnconProtocol.deployed();
  // const contract2 = AnconProtocol__factory.connect(contract.address, provider);

  // // z = toABIproofs();
  // // console.log(z);
  // // const resRootCalc = await contract2.callStatic.queryRootCalculation({
  // //   ...z,
  // // });

  // // const restUpdtHeader = await contract.updateProtocolHeader(resRootCalc, {
  // //   from: accounts[0],
  // // });
  // // console.log(resRootCalc);
  // // console.log(restUpdtHeader);

  builder.addContract('AnconProtocol', verifier, verifier.address, network)
}

function toABIproofs() {
  let z = { ...proofCombined[0].exist }
  z.key = hexlify(base64.decode(z.key))
  z.value = hexlify(base64.decode(z.value))
  z.leaf.prefix = hexlify(base64.decode(z.leaf.prefix))
  z.leaf.hash = 1
  z.path = z.path.map((x) => {
    let suffix
    if (!!x.suffix) {
      suffix = hexlify(base64.decode(x.suffix))
      return {
        valid: true,
        prefix: hexlify(base64.decode(x.prefix)),
        suffix: suffix,
        hash: 1,
      }
    } else {
      return {
        valid: true,
        prefix: hexlify(base64.decode(x.prefix)),
        hash: 1,
        suffix: '0x',
      }
    }
  })
  z.leaf.prehash_key = 0
  z.leaf.len = z.leaf.length

  return z
}
