/**
 * This example will use its app key as a signing key, and sign anything it is asked to.
 */

const ethers = require('ethers')

/*
 * The `wallet` API is a superset of the standard provider,
 * and can be used to initialize an ethers.js provider like this:
 */
const provider = new ethers.providers.Web3Provider(wallet)

wallet.registerRpcMessageHandler(async (_originString, requestObject) => {
  console.log('received request', requestObject)
  const privKey = await wallet.request({
    method: 'snap_getAppKey',
  })
  console.log(`privKey is ${privKey}`)
  const ethWallet = new ethers.Wallet(privKey, provider)
  console.dir(ethWallet)
  switch (requestObject.method) {
    case 'address':
      return ethWallet.address

    case 'signMessage': {
      const message = requestObject.params[0]

      const result = await wallet.request({
        method: 'snap_confirm',
        params: [
          {
            prompt: 'Would you like to take the action?',
            description: 'The action is...',
            textAreaContent: 'Very detailed information about the action...',
          },
        ],
      })

      if (result) {
        return ethWallet.signMessage(message)
      }
      return
    }

    case 'sign': {
      const message = requestObject.params[0]

      const result = await wallet.request({
        method: 'snap_confirm',
        params: [
          {
            prompt: 'Please sign message - Ancon Protocol',
            description: 'Signs message',
            textAreaContent: 'Name: Test playground \r\nDescription: NFT\r\nUser DID: did:ethr:bnbt:0x00\r\nChain: Binance Smart Chain',
          },
        ],
      })

      if (result) {
        return ethWallet.signMessage(message)
      }
      return
    }

    default:
      throw new Error('Method not found.')
  }
})
