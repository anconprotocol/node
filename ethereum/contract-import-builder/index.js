const fs = require('fs')
const networks = {
    ropsten: 3,
    rinkeby: 4,
    kovan: 42,
    mainnet: 1,
}

module.exports = class ContractImportBuilder {
    constructor() {
        this.importOutput = ''
        this.contractImport = {
            VERSION: '1.0.0'
        }
    }

    setOutput(output) {
        this.importOutput = output
    }

    addContract(name, abi, address, network) {
        let addr = {}

        if (this.importOutput.length < 5) {
            throw new Error(
                'A contract import must be set'
            )
        }

        fs.exists(this.importOutput, exists => {
            if (exists) {
                const contractImportExisting = require(this.importOutput)

                // get any existing addresses for the current selected network
                if (contractImportExisting && contractImportExisting[name]) {
                    addr = contractImportExisting[name].address
                }
            }

            addr[network] = address

            this.contractImport[name] = {
                // raw: require(`./build/contracts/${abi.contractName}.json`),
                raw: {
                    abi: abi.abi
                },
                address: {
                    ...addr
                }
            }

            const output = JSON.stringify(this.contractImport);
            const jsOutput = `module.exports = ${output}`;
            if (this.onWrite) {
                this.onWrite(jsOutput)
            }
        })
    }
}
