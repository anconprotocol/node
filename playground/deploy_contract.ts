import fetch from 'node-fetch'
import { hexlify } from '@ethersproject/bytes'
import fs from 'fs'
import { promisify } from 'util'
import { base64 } from 'ethers/lib/utils'
import { ethers } from 'ethers'
import { Command } from 'commander'
require('dotenv').config()
global['fetch'] = require('node-fetch')

export class DeployContractTool {
  constructor(
    private privateKey: string,
    private apiUrl: string = 'http://localhost:7788',
  ) {}

  async deploy(wasm: string): Promise<any> {
    const file = await promisify(fs.readFile)(wasm)
    const code = base64.encode(file)
    const wallet = new ethers.Wallet(this.privateKey)
    const from = await wallet.getAddress()
    let data = {
      from,
      code,
      signature: '',
    }
    const hash = ethers.utils.keccak256(JSON.stringify(data))
    const sig = await wallet.signMessage(hash)
    data.signature = sig

    const res = await (
      await fetch(`${this.apiUrl}/v0/code`, {
        body: JSON.stringify(data),
        headers: {
          'Content-Type': 'application/json',
        },
      })
    ).json()

    return res
  }
}

;(async () => {
  const deployer = new DeployContractTool(
    process.env.PRIVATE_KEY,
    process.env.URL,
  )

  const result = await deployer.deploy(process.env.WASM_PATH)

  console.log(result)
})()
