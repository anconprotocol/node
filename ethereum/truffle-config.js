const HDWalletProvider = require('@truffle/hdwallet-provider')
require('dotenv').config()

module.exports = {
  compilers: {
    solc: {
      version: '^0.8.7',
      metadataLiteral: true,
      settings: {
        // See the solidity docs for advice about optimization and evmVersion
        optimizer: {
          enabled: true,
          runs: 1000,
        },
      },
    },
  },
  api_keys: {
    bscscan: process.env.ETHERSCAN_API,
  },
  plugins: ['truffle-contract-size', 'truffle-plugin-verify'],
  networks: {
    bsctestnet: {
      provider: () =>
        new HDWalletProvider(process.env.BSC_MNEMONIC, process.env.BSC_TESTNET),
      network_id: 97,
      confirmations: 10,
      timeoutBlocks: 200,
      skipDryRun: false,
    },
    bsc: {
      provider: () =>
        new HDWalletProvider(
          process.env.BSCMAINNET_MNEMONIC,
          `https://bsc-dataseed1.binance.org`,
        ),
      network_id: 56,
      confirmations: 10,
      timeoutBlocks: 200,
      ///gas: 5000000,
      // gasPrice: 135000000000,
      skipDryRun: false,
    },
    localhost: {
      from: '0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6',
      host: 'localhost',
      port: 8545,
      network_id: '*', // Match any network id
    },
    development: {
      from: '0x32A21c1bB6E7C20F547e930b53dAC57f42cd25F6',
      host: 'localhost',
      port: 8545,
      network_id: '*', // Match any network id
    },
    ancon: {
      provider: () =>
        new HDWalletProvider(process.env.MNEMONIC, process.env.ANCON),
      network_id: '9000',
      gas: 7000000,
      gasPrice: 30000000,
    },
    rinkeby: {
      provider: () =>
        new HDWalletProvider(process.env.MNEMONIC, process.env.RINKEBY),
      network_id: '*',
      gas: 7000000,
      gasPrice: 3000000000,
    },
    kovan: {
      provider: () =>
        new HDWalletProvider(process.env.MNEMONIC, process.env.KOVAN),
      network_id: 42,
      gas: 7000000,
      gasPrice: 30000000000,
    },
    ropsten: {
      provider: () =>
        new HDWalletProvider(process.env.MNEMONIC, process.env.ROPSTEN),
      network_id: 3,
      gas: 5000000,
      // timeoutBlocks: 3,
      gasPrice: 65000000000,
    },
  },
}
