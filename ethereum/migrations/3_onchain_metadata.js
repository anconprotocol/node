const OnChainMet = artifacts.require('OnchainMetadata')

module.exports = async function (deployer) {
  try{
  await deployer.deploy(OnChainMet)
  } catch (e) {










    console.log(e)
  }
}
