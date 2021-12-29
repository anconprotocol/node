module.exports = {"VERSION":"1.0.0","AnconProtocol":{"raw":{"abi":[{"inputs":[{"internalType":"address","name":"_onlyOwner","type":"address"}],"stateMutability":"nonpayable","type":"constructor","signature":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bool","name":"enrolledStatus","type":"bool"},{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"value","type":"bytes"}],"name":"EnrollL2Account","type":"event","signature":"0x13280ffee0312bd0feaf5a148e6153d9636f3657872c53a8a058371fbd1e4914"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"packet","type":"bytes"}],"name":"ProofPacketSubmitted","type":"event","signature":"0x10a499eb855a3bf46db4fa7a4aa05f939a5d06c8a3a96ad7f4d840ee9817924e"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"accountByAddrProofs","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x73d15717"},{"inputs":[{"internalType":"bytes","name":"","type":"bytes"}],"name":"accountProofs","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x01644028"},{"inputs":[],"name":"getIavlSpec","outputs":[{"components":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum ICS23.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct ICS23.LeafOp","name":"leafSpec","type":"tuple"},{"components":[{"internalType":"uint256[]","name":"childOrder","type":"uint256[]"},{"internalType":"uint256","name":"childSize","type":"uint256"},{"internalType":"uint256","name":"minPrefixLength","type":"uint256"},{"internalType":"uint256","name":"maxPrefixLength","type":"uint256"},{"internalType":"bytes","name":"emptyChild","type":"bytes"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"}],"internalType":"struct ICS23.InnerSpec","name":"innerSpec","type":"tuple"},{"internalType":"uint256","name":"maxDepth","type":"uint256"},{"internalType":"uint256","name":"minDepth","type":"uint256"}],"internalType":"struct ICS23.ProofSpec","name":"","type":"tuple"}],"stateMutability":"pure","type":"function","constant":true,"signature":"0x27dcd78c"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x8da5cb5b"},{"inputs":[{"internalType":"bytes","name":"","type":"bytes"}],"name":"proofs","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xd56a07e3"},{"inputs":[],"name":"relayNetworkHash","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x410d0d84"},{"inputs":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum ICS23.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct ICS23.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct ICS23.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct ICS23.ExistenceProof","name":"proof","type":"tuple"},{"components":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum ICS23.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct ICS23.LeafOp","name":"leafSpec","type":"tuple"},{"components":[{"internalType":"uint256[]","name":"childOrder","type":"uint256[]"},{"internalType":"uint256","name":"childSize","type":"uint256"},{"internalType":"uint256","name":"minPrefixLength","type":"uint256"},{"internalType":"uint256","name":"maxPrefixLength","type":"uint256"},{"internalType":"bytes","name":"emptyChild","type":"bytes"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"}],"internalType":"struct ICS23.InnerSpec","name":"innerSpec","type":"tuple"},{"internalType":"uint256","name":"maxDepth","type":"uint256"},{"internalType":"uint256","name":"minDepth","type":"uint256"}],"internalType":"struct ICS23.ProofSpec","name":"spec","type":"tuple"},{"internalType":"bytes","name":"root","type":"bytes"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"}],"name":"verify","outputs":[],"stateMutability":"pure","type":"function","constant":true,"signature":"0xb0d264e7"},{"inputs":[{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"did","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum ICS23.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct ICS23.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct ICS23.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct ICS23.ExistenceProof","name":"proof","type":"tuple"}],"name":"enrollL2Account","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"payable","type":"function","payable":true,"signature":"0x32208491"},{"inputs":[{"internalType":"bytes","name":"rootHash","type":"bytes"}],"name":"updateProtocolHeader","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function","signature":"0xc935256b"},{"inputs":[{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"packet","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum ICS23.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct ICS23.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct ICS23.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct ICS23.ExistenceProof","name":"proof","type":"tuple"}],"name":"submitPacketWithProof","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"payable","type":"function","payable":true,"signature":"0x354763d9"},{"inputs":[{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"internalType":"bytes","name":"_prefix","type":"bytes"},{"internalType":"bytes[][]","name":"_innerOp","type":"bytes[][]"}],"name":"convertProof","outputs":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum ICS23.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct ICS23.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct ICS23.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct ICS23.ExistenceProof","name":"","type":"tuple"}],"stateMutability":"pure","type":"function","constant":true,"signature":"0x13f85fee"},{"inputs":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum ICS23.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct ICS23.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct ICS23.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct ICS23.ExistenceProof","name":"exProof","type":"tuple"}],"name":"verifyProof","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x391506e5"},{"inputs":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum ICS23.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum ICS23.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct ICS23.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum ICS23.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct ICS23.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct ICS23.ExistenceProof","name":"proof","type":"tuple"}],"name":"queryRootCalculation","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"pure","type":"function","constant":true,"signature":"0xd4f3237c"}]},"address":{"development":"0x2AE945d577685aBa3A8B564D255454a77AFfDf6c"}}}