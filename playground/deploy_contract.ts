import dotenv from 'dotenv'
import { arrayify, hexlify } from '@ethersproject/bytes'
import fs from 'fs'
import { promisify } from 'util'
import { ethers } from 'ethers'
import axios from 'axios'
import { Secp256k1 } from '@cosmjs/crypto'
dotenv.config()

export class DeployContractTool {
  constructor(
    private privateKey: string,
    private apiUrl: string = 'http://localhost:7788',
  ) {}

  async deploy(wasm: string): Promise<any> {
    const file = await promisify(fs.readFile)(wasm)
    const code = ethers.utils.hexlify(file)
    let data = {
      from: process.env.CID,
      code,
      // signature: '',
    } as any
    const bz = ethers.utils.toUtf8Bytes(code+data.from)
    const hash = ethers.utils.keccak256(bz)
    const sig = await Secp256k1.createSignature(
      arrayify(hash),
      Buffer.from(this.privateKey, 'hex'),
    )
    data.signature = hexlify((sig.toFixedLength()))
    const res = await axios.post(`${this.apiUrl}/v0/code`, {
      body: data,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    return res.data
  }
}

;(async () => {
  const deployer = new DeployContractTool(
    process.env.PRIVATE_KEY || '',
    process.env.URL,
  )

  //@ts-ignore
  const result = await deployer.deploy(process.env.WASM_PATH)

  console.log(result)
})()
