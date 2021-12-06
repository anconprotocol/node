const OnChainMet = artifacts.require('OnchainMetadata')

module.exports = async function (deployer) {
  await deployer.deploy(OnChainMet)
}
