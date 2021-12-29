// const {utils, ethers} = require("ethers")
// const {AncProt__Factory} = require("../types/ethers-contracts/factories/AnconProtocol__factory.ts");
import { utils, ethers } from 'ethers'
import { ExistenceProofStruct } from '../types/ethers-contracts/AnconProtocol'
import { AnconProtocol__factory } from '../types/ethers-contracts/factories/AnconProtocol__factory'
import { arrayify, base64 } from 'ethers/lib/utils'

const RPC_HOST = 'http://localhost:8545/'
const CONTRACT_ADDRESS = '0x77C51E844495899727dB63221af46220b0b13B37'
const BLOCK_NUMBER = 13730326
const ACCOUNT_ADDRESS = '0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6'

let proofCombined = [
  {
    exist: {
      key:
        'L2FuY29ucHJvdG9jb2wvZTA1NDZjMDZlNDZlYWIzNDcyMmVhMTNjNTAyNGNiMDBmYjEzNmVmZDg3OGY0NThiNTViMDQ3YzhkOGU4Y2JiNi91c2VyL2JhZ3VxZWVyYWVocGRhN3pwcmJ1bXhoZzVncWlmbHkzYm1kbmFlb2NxbzRmbHZub2ZzaXB0ZmVoa3F5d3E=',
      leaf: {
        hash: 1,
        length: 1,
        prefix: 'AAIi',
        prehash_value: 1,
      },
      path: [
        {
          hash: 1,
          prefix: 'AgQiIA==',
          suffix: 'IEga3tvbzXYCsg0xSYPX36OPlz1bwOPxMp239NZ9XwG+',
        },
        {
          hash: 1,
          prefix: 'BAgiIKLI0YDL04HqPwDpHF8ht5pGCiq+Uw/0DTzSMDp7pE6mIA==',
        },
        {
          hash: 1,
          prefix: 'BgwiIADQ5R72VKesArREiAa1kDeAFiOgt9tZ+32NLHPOPjP9IA==',
        },
        {
          hash: 1,
          prefix: 'CBIiIA==',
          suffix: 'IPhUy0NkQUD/fk42UPtKzT0QWd1NDgJshHnLRSm3j7XQ',
        },
        {
          hash: 1,
          prefix: 'CiAiIA==',
          suffix: 'IMo1qcvqM7Duwq5Ac3wtuoTitJjOTDVrB92+pkPPAlzW',
        },
      ],
      value: 'ZGlkOndlYjppcGZzOnVzZXI6dGVzdA==',
    },
  },
]

let exProof: ExistenceProofStruct
exProof = [] as any

function toABIproofs(proofCombined: any) {
  let innerOps = proofCombined[0].exist.path.map((p) => [p.prefix, p.suffix])

  let key = proofCombined[0].exist.key

  return {
    key,
    value: proofCombined[0].exist.value,
    prefix: proofCombined[0].exist.leaf.prefix,
    innerOps: innerOps,
  }
}

async function main() {
  const provider = new ethers.providers.JsonRpcProvider(RPC_HOST)
  const contract = AnconProtocol__factory.connect(CONTRACT_ADDRESS, provider)

  exProof.valid = true
  exProof.value = proofCombined[0].exist.value
  // await contract.updateProtocolHeader()
  const resRootCalc = await contract.queryRootCalculation({
    valid: false,
    key: '',
    value: '',
    leaf: {
      valid: false,
      hash: '',
      prehash_key: '',
      prehash_value: '',
      len: '',
      prefix: '',
    },
    path: [],
  })
  console.log(resRootCalc)
  // await contract.updateProtocolHeader(resRootCalc)
  const resVerifyProof = await contract.verifyProof(
    base64.decode(exProof.key),
    base64.decode(exProof.value),
    base64.decode(exProof.leaf.prefix as any),
    base64.decode(exProof.path[0].prefix as any),
    base64.decode(exProof.path[0].suffix as any),
    arrayify(resRootCalc),
  )

  const contractName = 'ANCON PROTOCOL'
  console.log(`Our ${contractName} root calculation is: ${resRootCalc}`)
  console.log(`Our ${contractName} proof is: ${resVerifyProof}`)

  console.log(`Listing Transfer events for block ${BLOCK_NUMBER}`)

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
  console.error(e)
  process.exit(1)
})
