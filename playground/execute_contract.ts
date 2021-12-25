import dotenv from 'dotenv'
import { arrayify, hexlify } from '@ethersproject/bytes'
import fs from 'fs'
import { promisify } from 'util'
import { ethers } from 'ethers'
import web3 from 'web3'
import { Secp256k1 } from '@cosmjs/crypto'
import { fetchJson } from '@ethersproject/web'
import { toUtf8Bytes, toUtf8String } from 'ethers/lib/utils'
dotenv.config()

export class ExecuteContractExample {
  constructor(
    private privateKey: string,
    private apiUrl: string = 'http://localhost:7788',
  ) {}

  async execute({ to, from, query }: any): Promise<any> {
    let data = {
      to,
      from,
      data: query,
      // signature: '',
    } as any
    const bz = ethers.utils.toUtf8Bytes(data.to + data.from + data.data)
    const hash = ethers.utils.keccak256(bz)
    const sig = await Secp256k1.createSignature(
      arrayify(hash),
      Buffer.from(this.privateKey, 'hex'),
    )
    data.signature = hexlify(sig.toFixedLength())

    const body = {
      jsonrpc: '2.0',
      method: 'ancon_call',
      params: [data.to, data.from, data.signature, data.data],
      id: 1,
    }
    
    const result = await fetchJson(
      {
        url: `${this.apiUrl}/rpc`,
      },
      JSON.stringify(body),
    )

    return result
  }
}

;(async () => {
  const contract = new ExecuteContractExample(
    process.env.PRIVATE_KEY || '',
    process.env.URL,
  )

  //@ts-ignore
  let metadata = 'baguqeerak5hfvtgsvaaxtm5cbce2khnh7y4ijfnbrgeygwu4wixrz4z52vja'
  const result = await contract.execute({
    to: process.env.CONTRACT,
    from: process.env.DID_USER,
    query: `query   { metadata(cid: \"${metadata}\", path: \"/\"){name} }`
  })

  console.log(result)
  const j = (web3.utils.hexToString(result.result))
  console.log(JSON.parse(j))
  
})()
