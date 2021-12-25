import http from 'http'
import dotenv from 'dotenv'
import { arrayify, hexlify, hexValue } from '@ethersproject/bytes'
import fs from 'fs'
import { promisify } from 'util'
import { base64 } from 'ethers/lib/utils'
import { ethers } from 'ethers'
import { Command } from 'commander'
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
    const wallet = new ethers.Wallet(this.privateKey)
    const from = await wallet.getAddress()
    let data = {
      from: process.env.CID,
      code,
      // signature: '',
    } as any
    const bz = ethers.utils.toUtf8Bytes(code+data.from)
    const hash = ethers.utils.keccak256(bz)
    const sig = await Secp256k1.createSignature(
      arrayify(hash),
      Buffer.from(process.env.PRIVATE_KEY || '', 'hex'),
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
