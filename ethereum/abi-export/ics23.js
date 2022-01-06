module.exports = {"VERSION":"1.0.0","XDVNFT":{"raw":{"abi":[{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"},{"internalType":"address","name":"tokenERC20","type":"address"},{"internalType":"address","name":"anconprotocolAddr","type":"address"}],"stateMutability":"nonpayable","type":"constructor","signature":"constructor"},{"inputs":[{"internalType":"string","name":"url","type":"string"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"name":"OffchainLookup","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"approved","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Approval","type":"event","signature":"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"bool","name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event","signature":"0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event","signature":"0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event","signature":"0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"bytes32","name":"signatureHash","type":"bytes32"}],"name":"ProofAccepted","type":"event","signature":"0xe37c3073f0ac6c2b04fc00b9673ac5cb0fbedc3e21d2b038714c08fb312df56c"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"uint256","name":"paidToContract","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"paidToPaymentAddress","type":"uint256"}],"name":"ServiceFeePaid","type":"event","signature":"0xff781da90c6b849a72bafc99aedde3dff0b039f697f31c21e96d587accee2760"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event","signature":"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event","signature":"0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"paymentAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"Withdrawn","type":"event","signature":"0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5"},{"inputs":[],"name":"anconprotocol","outputs":[{"internalType":"contract IAnconProtocol","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x3cd559a4"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"approve","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x095ea7b3"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x70a08231"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"burn","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x42966c68"},{"inputs":[],"name":"dagContractOperator","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x7927b898"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x081812fc"},{"inputs":[],"name":"getSigner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x7ac3c02f"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xe985e9c5"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x06fdde03"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x8da5cb5b"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x6352211e"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x5c975abb"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x715018a6"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x42842e0e"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0xb88d4fde"},{"inputs":[],"name":"serviceFeeForContract","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xe53ae102"},{"inputs":[],"name":"serviceFeeForPaymentAddress","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xc014b9da"},{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"bool","name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0xa22cb465"},{"inputs":[{"internalType":"address","name":"signer_","type":"address"}],"name":"setSigner","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x6c19e783"},{"inputs":[{"internalType":"string","name":"url_","type":"string"}],"name":"setUrl","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x252498a2"},{"inputs":[],"name":"stablecoin","outputs":[{"internalType":"contract IERC20","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xe9cbd822"},{"inputs":[{"internalType":"bytes4","name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x01ffc9a7"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x95d89b41"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x23b872dd"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0xf2fde38b"},{"inputs":[],"name":"url","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x5600f04f"},{"inputs":[{"internalType":"uint256","name":"_fee","type":"uint256"}],"name":"setServiceFeeForPaymentAddress","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0xb86be2ca"},{"inputs":[{"internalType":"uint256","name":"_fee","type":"uint256"}],"name":"setServiceFeeForContract","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x2c078e8d"},{"inputs":[{"internalType":"address","name":"toAddress","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"transferURI","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"nonpayable","type":"function","signature":"0xa0d9c721"},{"inputs":[{"internalType":"address","name":"toAddress","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"mint","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"nonpayable","type":"function","signature":"0x40c10f19"},{"inputs":[{"internalType":"string","name":"metadataUri","type":"string"},{"internalType":"address","name":"transferTo","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"transferMetadataOwnership","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"nonpayable","type":"function","signature":"0x2d4da937"},{"inputs":[{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"packet","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum Ics23Helper.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct Ics23Helper.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct Ics23Helper.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct Ics23Helper.ExistenceProof","name":"userProof","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum Ics23Helper.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct Ics23Helper.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct Ics23Helper.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct Ics23Helper.ExistenceProof","name":"proof","type":"tuple"}],"name":"mintWithProof","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function","signature":"0x1777110d"},{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"address","name":"from","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"onERC721Received","outputs":[{"internalType":"bytes4","name":"","type":"bytes4"}],"stateMutability":"nonpayable","type":"function","signature":"0x150b7a02"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xc87b56dd"},{"inputs":[{"internalType":"address payable","name":"payee","type":"address"}],"name":"withdrawBalance","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x756af45f"}]},"address":{"bsctestnet":"0x46d890d9e9BdB91bD7d31a5D8262586baD6A9399"}},"AnconProtocol":{"raw":{"abi":[{"inputs":[],"stateMutability":"nonpayable","type":"constructor","signature":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bool","name":"enrolledStatus","type":"bool"},{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"value","type":"bytes"}],"name":"AccountRegistered","type":"event","signature":"0xdf285b75c61e111633bd9fdd496ccf9b449555e34d68ea7bf9db18e8f1977c70"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"hash","type":"bytes"}],"name":"HeaderUpdated","type":"event","signature":"0xe0b001f59b54160030a2302b411d234315941c6c1d33a52bdb8f3c46a1dffeb8"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event","signature":"0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"packet","type":"bytes"}],"name":"ProofPacketSubmitted","type":"event","signature":"0x10a499eb855a3bf46db4fa7a4aa05f939a5d06c8a3a96ad7f4d840ee9817924e"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"}],"name":"ServiceFeePaid","type":"event","signature":"0xa70c9ef1994019c7c70e8134256a652460b545755ed8aad140daeaccc30446b3"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"paymentAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"Withdrawn","type":"event","signature":"0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5"},{"inputs":[],"name":"ENROLL_PAYMENT","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x40ae5ffa"},{"inputs":[],"name":"SUBMIT_PAYMENT","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xfc7ea420"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"accountByAddrProofs","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x73d15717"},{"inputs":[{"internalType":"bytes","name":"","type":"bytes"}],"name":"accountProofs","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x01644028"},{"inputs":[],"name":"accountRegistrationFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x5dad28de"},{"inputs":[],"name":"getIavlSpec","outputs":[{"components":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum Ics23Helper.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct Ics23Helper.LeafOp","name":"leafSpec","type":"tuple"},{"components":[{"internalType":"uint256[]","name":"childOrder","type":"uint256[]"},{"internalType":"uint256","name":"childSize","type":"uint256"},{"internalType":"uint256","name":"minPrefixLength","type":"uint256"},{"internalType":"uint256","name":"maxPrefixLength","type":"uint256"},{"internalType":"bytes","name":"emptyChild","type":"bytes"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"}],"internalType":"struct Ics23Helper.InnerSpec","name":"innerSpec","type":"tuple"},{"internalType":"uint256","name":"maxDepth","type":"uint256"},{"internalType":"uint256","name":"minDepth","type":"uint256"}],"internalType":"struct Ics23Helper.ProofSpec","name":"","type":"tuple"}],"stateMutability":"pure","type":"function","constant":true,"signature":"0x27dcd78c"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x8da5cb5b"},{"inputs":[{"internalType":"bytes","name":"","type":"bytes"}],"name":"proofs","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xd56a07e3"},{"inputs":[],"name":"protocolFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xb0e21e8a"},{"inputs":[],"name":"relayNetworkHash","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x410d0d84"},{"inputs":[],"name":"relayer","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x8406c079"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x715018a6"},{"inputs":[],"name":"stablecoin","outputs":[{"internalType":"contract IERC20","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xe9cbd822"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0xf2fde38b"},{"inputs":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum Ics23Helper.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct Ics23Helper.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct Ics23Helper.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct Ics23Helper.ExistenceProof","name":"proof","type":"tuple"},{"components":[{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum Ics23Helper.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct Ics23Helper.LeafOp","name":"leafSpec","type":"tuple"},{"components":[{"internalType":"uint256[]","name":"childOrder","type":"uint256[]"},{"internalType":"uint256","name":"childSize","type":"uint256"},{"internalType":"uint256","name":"minPrefixLength","type":"uint256"},{"internalType":"uint256","name":"maxPrefixLength","type":"uint256"},{"internalType":"bytes","name":"emptyChild","type":"bytes"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"}],"internalType":"struct Ics23Helper.InnerSpec","name":"innerSpec","type":"tuple"},{"internalType":"uint256","name":"maxDepth","type":"uint256"},{"internalType":"uint256","name":"minDepth","type":"uint256"}],"internalType":"struct Ics23Helper.ProofSpec","name":"spec","type":"tuple"},{"internalType":"bytes","name":"root","type":"bytes"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"}],"name":"verify","outputs":[],"stateMutability":"pure","type":"function","constant":true,"signature":"0xb0d264e7"},{"inputs":[{"internalType":"contract IERC20","name":"tokenAddress","type":"address"}],"name":"setPaymentToken","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x6a326ab1"},{"inputs":[{"internalType":"address payable","name":"payee","type":"address"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x51cff8d9"},{"inputs":[{"internalType":"address payable","name":"payee","type":"address"},{"internalType":"address","name":"erc20token","type":"address"}],"name":"withdrawToken","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x3aeac4e1"},{"inputs":[{"internalType":"uint256","name":"_fee","type":"uint256"}],"name":"setProtocolFee","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x787dce3d"},{"inputs":[{"internalType":"uint256","name":"_fee","type":"uint256"}],"name":"setAccountRegistrationFee","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0x8b178bec"},{"inputs":[],"name":"getProtocolHeader","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xf2a4147a"},{"inputs":[{"internalType":"bytes","name":"did","type":"bytes"}],"name":"getProof","outputs":[{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x693ac4fb"},{"inputs":[{"internalType":"bytes","name":"key","type":"bytes"}],"name":"hasProof","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true,"signature":"0xe9f49b53"},{"inputs":[{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"did","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum Ics23Helper.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct Ics23Helper.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct Ics23Helper.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct Ics23Helper.ExistenceProof","name":"proof","type":"tuple"}],"name":"enrollL2Account","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"payable","type":"function","payable":true,"signature":"0x32208491"},{"inputs":[{"internalType":"bytes","name":"rootHash","type":"bytes"}],"name":"updateProtocolHeader","outputs":[],"stateMutability":"nonpayable","type":"function","signature":"0xc935256b"},{"inputs":[{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"packet","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum Ics23Helper.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct Ics23Helper.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct Ics23Helper.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct Ics23Helper.ExistenceProof","name":"proof","type":"tuple"}],"name":"submitPacketWithProof","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"payable","type":"function","payable":true,"signature":"0x354763d9"},{"inputs":[{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"bytes","name":"key","type":"bytes"},{"internalType":"bytes","name":"value","type":"bytes"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_key","type":"uint8"},{"internalType":"enum Ics23Helper.HashOp","name":"prehash_value","type":"uint8"},{"internalType":"enum Ics23Helper.LengthOp","name":"len","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"}],"internalType":"struct Ics23Helper.LeafOp","name":"leaf","type":"tuple"},{"components":[{"internalType":"bool","name":"valid","type":"bool"},{"internalType":"enum Ics23Helper.HashOp","name":"hash","type":"uint8"},{"internalType":"bytes","name":"prefix","type":"bytes"},{"internalType":"bytes","name":"suffix","type":"bytes"}],"internalType":"struct Ics23Helper.InnerOp[]","name":"path","type":"tuple[]"}],"internalType":"struct Ics23Helper.ExistenceProof","name":"exProof","type":"tuple"}],"name":"verifyProofWithKV","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true,"signature":"0x01b3239e"}]},"address":{"development":"0x71E56696Eb1A1d0b0e96A01A03DA7481e0008F3F","rinkeby-fork":"0xb0c578D19f6E7dD455798b76CC92FfdDb61aD635","ropsten-fork":"0x067F04932a4808124541521E940fD2d21a44fEeB","bsctestnet":"0x846399a87FeaE3915F43a871075970DB40C25673"}}}