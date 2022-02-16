/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import {
  Signer,
  utils,
  Contract,
  ContractFactory,
  Overrides,
  BigNumberish,
} from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type { KYX, KYXInterface } from "../KYX";

const _abi = [
  {
    inputs: [
      {
        internalType: "address",
        name: "tokenERC20",
        type: "address",
      },
      {
        internalType: "address",
        name: "ancon",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "chain",
        type: "uint256",
      },
    ],
    stateMutability: "nonpayable",
    type: "constructor",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "uint256",
        name: "id",
        type: "uint256",
      },
      {
        indexed: true,
        internalType: "bytes32",
        name: "category",
        type: "bytes32",
      },
      {
        indexed: false,
        internalType: "string",
        name: "metadata",
        type: "string",
      },
    ],
    name: "IssuerAdded",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "previousOwner",
        type: "address",
      },
      {
        indexed: true,
        internalType: "address",
        name: "newOwner",
        type: "address",
      },
    ],
    name: "OwnershipTransferred",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "payee",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "weiAmount",
        type: "uint256",
      },
    ],
    name: "Withdrawn",
    type: "event",
  },
  {
    inputs: [],
    name: "anconprotocol",
    outputs: [
      {
        internalType: "contract IAnconProtocol",
        name: "",
        type: "address",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "fee",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "bytes32",
        name: "",
        type: "bytes32",
      },
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    name: "issuers",
    outputs: [
      {
        internalType: "uint256",
        name: "id",
        type: "uint256",
      },
      {
        internalType: "bytes32",
        name: "category",
        type: "bytes32",
      },
      {
        internalType: "string",
        name: "metadata",
        type: "string",
      },
      {
        internalType: "uint256",
        name: "reputation",
        type: "uint256",
      },
      {
        internalType: "bool",
        name: "enabled",
        type: "bool",
      },
      {
        internalType: "address",
        name: "creator",
        type: "address",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "bytes32",
        name: "",
        type: "bytes32",
      },
    ],
    name: "issuersCount",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "owner",
    outputs: [
      {
        internalType: "address",
        name: "",
        type: "address",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [],
    name: "renounceOwnership",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "stablecoin",
    outputs: [
      {
        internalType: "contract IERC20",
        name: "",
        type: "address",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "newOwner",
        type: "address",
      },
    ],
    name: "transferOwnership",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address payable",
        name: "payee",
        type: "address",
      },
    ],
    name: "withdraw",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address payable",
        name: "payee",
        type: "address",
      },
      {
        internalType: "address",
        name: "erc20token",
        type: "address",
      },
    ],
    name: "withdrawToken",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "uint256",
        name: "_fee",
        type: "uint256",
      },
    ],
    name: "setFee",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "getFee",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "bytes32",
        name: "category",
        type: "bytes32",
      },
    ],
    name: "getIssuerLength",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "bytes32",
        name: "category",
        type: "bytes32",
      },
      {
        internalType: "uint256",
        name: "id",
        type: "uint256",
      },
    ],
    name: "getIssuer",
    outputs: [
      {
        components: [
          {
            internalType: "uint256",
            name: "id",
            type: "uint256",
          },
          {
            internalType: "bytes32",
            name: "category",
            type: "bytes32",
          },
          {
            internalType: "string",
            name: "metadata",
            type: "string",
          },
          {
            internalType: "uint256",
            name: "reputation",
            type: "uint256",
          },
          {
            internalType: "bool",
            name: "enabled",
            type: "bool",
          },
          {
            internalType: "address",
            name: "creator",
            type: "address",
          },
        ],
        internalType: "struct KYX.Issuer",
        name: "",
        type: "tuple",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "bytes32",
        name: "moniker",
        type: "bytes32",
      },
      {
        internalType: "bytes",
        name: "packet",
        type: "bytes",
      },
      {
        components: [
          {
            internalType: "bool",
            name: "valid",
            type: "bool",
          },
          {
            internalType: "bytes",
            name: "key",
            type: "bytes",
          },
          {
            internalType: "bytes",
            name: "value",
            type: "bytes",
          },
          {
            components: [
              {
                internalType: "bool",
                name: "valid",
                type: "bool",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "hash",
                type: "uint8",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "prehash_key",
                type: "uint8",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "prehash_value",
                type: "uint8",
              },
              {
                internalType: "enum Ics23Helper.LengthOp",
                name: "len",
                type: "uint8",
              },
              {
                internalType: "bytes",
                name: "prefix",
                type: "bytes",
              },
            ],
            internalType: "struct Ics23Helper.LeafOp",
            name: "leaf",
            type: "tuple",
          },
          {
            components: [
              {
                internalType: "bool",
                name: "valid",
                type: "bool",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "hash",
                type: "uint8",
              },
              {
                internalType: "bytes",
                name: "prefix",
                type: "bytes",
              },
              {
                internalType: "bytes",
                name: "suffix",
                type: "bytes",
              },
            ],
            internalType: "struct Ics23Helper.InnerOp[]",
            name: "path",
            type: "tuple[]",
          },
        ],
        internalType: "struct Ics23Helper.ExistenceProof",
        name: "userProof",
        type: "tuple",
      },
      {
        components: [
          {
            internalType: "bool",
            name: "valid",
            type: "bool",
          },
          {
            internalType: "bytes",
            name: "key",
            type: "bytes",
          },
          {
            internalType: "bytes",
            name: "value",
            type: "bytes",
          },
          {
            components: [
              {
                internalType: "bool",
                name: "valid",
                type: "bool",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "hash",
                type: "uint8",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "prehash_key",
                type: "uint8",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "prehash_value",
                type: "uint8",
              },
              {
                internalType: "enum Ics23Helper.LengthOp",
                name: "len",
                type: "uint8",
              },
              {
                internalType: "bytes",
                name: "prefix",
                type: "bytes",
              },
            ],
            internalType: "struct Ics23Helper.LeafOp",
            name: "leaf",
            type: "tuple",
          },
          {
            components: [
              {
                internalType: "bool",
                name: "valid",
                type: "bool",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "hash",
                type: "uint8",
              },
              {
                internalType: "bytes",
                name: "prefix",
                type: "bytes",
              },
              {
                internalType: "bytes",
                name: "suffix",
                type: "bytes",
              },
            ],
            internalType: "struct Ics23Helper.InnerOp[]",
            name: "path",
            type: "tuple[]",
          },
        ],
        internalType: "struct Ics23Helper.ExistenceProof",
        name: "packetProof",
        type: "tuple",
      },
    ],
    name: "enrollIssuerWithProof",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "bytes32",
        name: "category",
        type: "bytes32",
      },
      {
        internalType: "uint256",
        name: "issuerID",
        type: "uint256",
      },
      {
        internalType: "string",
        name: "metadataUri",
        type: "string",
      },
    ],
    name: "setIssuerWithProof",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "bytes32",
        name: "category",
        type: "bytes32",
      },
      {
        internalType: "uint256",
        name: "issuerID",
        type: "uint256",
      },
      {
        internalType: "string",
        name: "metadataUri",
        type: "string",
      },
    ],
    name: "setIssuerRatingWithProof",
    outputs: [],
    stateMutability: "nonpayable",
    type: "function",
  },
];

const _bytecode =
  "0x608060405260006006553480156200001657600080fd5b5060405162001b0238038062001b028339810160408190526200003991620000ea565b62000044336200007d565b600480546001600160a01b039485166001600160a01b03199182161790915560058054939094169216919091179091556006556200012b565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b80516001600160a01b0381168114620000e557600080fd5b919050565b6000806000606084860312156200010057600080fd5b6200010b84620000cd565b92506200011b60208501620000cd565b9150604084015190509250925092565b6119c7806200013b6000396000f3fe608060405234801561001057600080fd5b506004361061011b5760003560e01c80638da5cb5b116100b2578063b543f44b11610081578063ddca3f4311610066578063ddca3f4314610258578063e9cbd82214610261578063f2fde38b1461027457600080fd5b8063b543f44b14610230578063ced72f871461025057600080fd5b80638da5cb5b146101ec57806392bf6c6c146101fd578063a4660320146101fd578063adb9b61e1461021057600080fd5b8063585ca9a3116100ee578063585ca9a31461019957806366f53645146101be57806369fe0e2d146101d1578063715018a6146101e457600080fd5b80631479f99e146101205780633aeac4e1146101465780633cd559a41461015b57806351cff8d914610186575b600080fd5b61013361012e366004610f51565b610287565b6040519081526020015b60405180910390f35b610159610154366004610f7f565b6102fa565b005b60055461016e906001600160a01b031681565b6040516001600160a01b03909116815260200161013d565b610159610194366004610fb8565b6104fc565b6101ac6101a7366004610fdc565b610645565b60405161013d9695949392919061105a565b6101336101cc366004611499565b610719565b6101596101df366004610f51565b610b14565b610159610b73565b6000546001600160a01b031661016e565b61015961020b36600461152b565b505050565b61013361021e366004610f51565b60016020526000908152604090205481565b61024361023e366004610fdc565b610bd9565b60405161013d919061158f565b600354610133565b61013360035481565b60045461016e906001600160a01b031681565b610159610282366004610fb8565b610d6e565b6000818152600160205260408120546102e75760405162461bcd60e51b815260206004820152601060248201527f6e6f206973737565727320666f756e640000000000000000000000000000000060448201526064015b60405180910390fd5b5060009081526001602052604090205490565b6000546001600160a01b031633146103545760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102de565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000906001600160a01b038316906370a0823190602401602060405180830381865afa1580156103b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103d891906115f4565b6040517fa9059cbb0000000000000000000000000000000000000000000000000000000081526001600160a01b038581166004830152602482018390529192509083169063a9059cbb906044016020604051808303816000875af1158015610444573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610468919061160d565b6104b45760405162461bcd60e51b815260206004820152600f60248201527f7472616e73666572206661696c6564000000000000000000000000000000000060448201526064016102de565b826001600160a01b03167f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5826040516104ef91815260200190565b60405180910390a2505050565b6000546001600160a01b031633146105565760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102de565b604051479060009081906001600160a01b0385169084908381818185875af1925050503d80600081146105a5576040519150601f19603f3d011682016040523d82523d6000602084013e6105aa565b606091505b5091509150816105fc5760405162461bcd60e51b815260206004820152601460248201527f4661696c656420746f2073656e6420457468657200000000000000000000000060448201526064016102de565b836001600160a01b03167f7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d58460405161063791815260200190565b60405180910390a250505050565b60026020818152600093845260408085209091529183529120805460018201549282018054919392916106779061162a565b80601f01602080910402602001604051908101604052809291908181526020018280546106a39061162a565b80156106f05780601f106106c5576101008083540402835291602001916106f0565b820191906000526020600020905b8154815290600101906020018083116106d357829003601f168201915b50505050600383015460049093015491929160ff8116915061010090046001600160a01b031686565b6005546040517f5eccc371000000000000000000000000000000000000000000000000000000008152600481018690526000917fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470916001600160a01b0390911690635eccc37190602401600060405180830381865afa1580156107a0573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526107c89190810190611695565b80519060200120141561081d5760405162461bcd60e51b815260206004820152600f60248201527f496e76616c6964206d6f6e696b6572000000000000000000000000000000000060448201526064016102de565b60055460208301516040517f97554e8f0000000000000000000000000000000000000000000000000000000081526001600160a01b03909216916397554e8f9161087491899133918991908b908a9060040161188f565b6020604051808303816000875af1158015610893573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108b7919061160d565b6109035760405162461bcd60e51b815260206004820152601460248201527f696e76616c6964207061636b65742070726f6f6600000000000000000000000060448201526064016102de565b60008060008680602001905181019061091c91906118fe565b6000838152600260209081526040808320858452909152902060040154929550909350915060ff1615156001146109bb5760405162461bcd60e51b815260206004820152602160248201527f69737375657220616c72656164792065786973747320616e6420656e61626c6560448201527f640000000000000000000000000000000000000000000000000000000000000060648201526084016102de565b6000838152600160208190526040909120546109d691611958565b600084815260016020818152604080842094909455835160c08101855286815280820188815281860187815260608301869052608083018590523360a084015289865260028085528787208a8852855296909520825181559051938101939093559251805193949293610a50938501929190910190610eb8565b506060820151600382015560808201516004909101805460a0909301516001600160a01b0316610100027fffffffffffffffffffffff0000000000000000000000000000000000000000ff921515929092167fffffffffffffffffffffff00000000000000000000000000000000000000000090931692909217179055604051839083907f8e59ef813ee9e64c487b9abfcbae81bf1389ec5fcf6bfb3b67cf6a3f892aebf190610b0190859061197e565b60405180910390a3509695505050505050565b6000546001600160a01b03163314610b6e5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102de565b600355565b6000546001600160a01b03163314610bcd5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102de565b610bd76000610e50565b565b6040805160c08101825260008082526020808301829052606083850181905283018290526080830182905260a08301829052858252600281528382208583529052919091206004015460ff1615610c725760405162461bcd60e51b815260206004820152601060248201527f6e6f206973737565727320666f756e640000000000000000000000000000000060448201526064016102de565b6000838152600260208181526040808420868552825292839020835160c0810185528154815260018201549281019290925291820180549193840191610cb79061162a565b80601f0160208091040260200160405190810160405280929190818152602001828054610ce39061162a565b8015610d305780601f10610d0557610100808354040283529160200191610d30565b820191906000526020600020905b815481529060010190602001808311610d1357829003601f168201915b50505091835250506003820154602082015260049091015460ff81161515604083015261010090046001600160a01b03166060909101529392505050565b6000546001600160a01b03163314610dc85760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016102de565b6001600160a01b038116610e445760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016102de565b610e4d81610e50565b50565b600080546001600160a01b038381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b828054610ec49061162a565b90600052602060002090601f016020900481019282610ee65760008555610f2c565b82601f10610eff57805160ff1916838001178555610f2c565b82800160010185558215610f2c579182015b82811115610f2c578251825591602001919060010190610f11565b50610f38929150610f3c565b5090565b5b80821115610f385760008155600101610f3d565b600060208284031215610f6357600080fd5b5035919050565b6001600160a01b0381168114610e4d57600080fd5b60008060408385031215610f9257600080fd5b8235610f9d81610f6a565b91506020830135610fad81610f6a565b809150509250929050565b600060208284031215610fca57600080fd5b8135610fd581610f6a565b9392505050565b60008060408385031215610fef57600080fd5b50508035926020909101359150565b60005b83811015611019578181015183820152602001611001565b83811115611028576000848401525b50505050565b60008151808452611046816020860160208601610ffe565b601f01601f19169290920160200192915050565b86815285602082015260c06040820152600061107960c083018761102e565b60608301959095525091151560808301526001600160a01b031660a0909101529392505050565b634e487b7160e01b600052604160045260246000fd5b60405160c0810167ffffffffffffffff811182821017156110d9576110d96110a0565b60405290565b6040516080810167ffffffffffffffff811182821017156110d9576110d96110a0565b60405160a0810167ffffffffffffffff811182821017156110d9576110d96110a0565b604051601f8201601f1916810167ffffffffffffffff8111828210171561114e5761114e6110a0565b604052919050565b600067ffffffffffffffff821115611170576111706110a0565b50601f01601f191660200190565b600061119161118c84611156565b611125565b90508281528383830111156111a557600080fd5b828260208301376000602084830101529392505050565b600082601f8301126111cd57600080fd5b610fd58383356020850161117e565b8015158114610e4d57600080fd5b80356111f5816111dc565b919050565b8035600681106111f557600080fd5b600060c0828403121561121b57600080fd5b6112236110b6565b90508135611230816111dc565b815261123e602083016111fa565b602082015261124f604083016111fa565b6040820152611260606083016111fa565b606082015260808201356009811061127757600080fd5b608082015260a082013567ffffffffffffffff81111561129657600080fd5b6112a2848285016111bc565b60a08301525092915050565b600082601f8301126112bf57600080fd5b8135602067ffffffffffffffff808311156112dc576112dc6110a0565b8260051b6112eb838201611125565b938452858101830193838101908886111561130557600080fd5b84880192505b858310156113c3578235848111156113235760008081fd5b88016080818b03601f190181131561133b5760008081fd5b6113436110df565b87830135611350816111dc565b8152604061135f8482016111fa565b89830152606080850135898111156113775760008081fd5b6113858f8c838901016111bc565b848401525092840135928884111561139f57600091508182fd5b6113ad8e8b868801016111bc565b908301525084525050918401919084019061130b565b98975050505050505050565b600060a082840312156113e157600080fd5b6113e9611102565b90506113f4826111ea565b8152602082013567ffffffffffffffff8082111561141157600080fd5b61141d858386016111bc565b6020840152604084013591508082111561143657600080fd5b611442858386016111bc565b6040840152606084013591508082111561145b57600080fd5b61146785838601611209565b6060840152608084013591508082111561148057600080fd5b5061148d848285016112ae565b60808301525092915050565b600080600080608085870312156114af57600080fd5b84359350602085013567ffffffffffffffff808211156114ce57600080fd5b6114da888389016111bc565b945060408701359150808211156114f057600080fd5b6114fc888389016113cf565b9350606087013591508082111561151257600080fd5b5061151f878288016113cf565b91505092959194509250565b60008060006060848603121561154057600080fd5b8335925060208401359150604084013567ffffffffffffffff81111561156557600080fd5b8401601f8101861361157657600080fd5b6115858682356020840161117e565b9150509250925092565b6020815281516020820152602082015160408201526000604083015160c060608401526115bf60e084018261102e565b9050606084015160808401526080840151151560a08401526001600160a01b0360a08501511660c08401528091505092915050565b60006020828403121561160657600080fd5b5051919050565b60006020828403121561161f57600080fd5b8151610fd5816111dc565b600181811c9082168061163e57607f821691505b6020821081141561165f57634e487b7160e01b600052602260045260246000fd5b50919050565b600061167361118c84611156565b905082815283838301111561168757600080fd5b610fd5836020830184610ffe565b6000602082840312156116a757600080fd5b815167ffffffffffffffff8111156116be57600080fd5b8201601f810184136116cf57600080fd5b6116de84825160208401611665565b949350505050565b634e487b7160e01b600052602160045260246000fd5b6006811061170c5761170c6116e6565b9052565b600081518084526020808501808196508360051b8101915082860160005b8581101561179f5782840389528151608081511515865286820151611755888801826116fc565b50604080830151828289015261176d8389018261102e565b925050506060808301519250868203818801525061178b818361102e565b9a87019a955050509084019060010161172e565b5091979650505050505050565b8051151582526000602082015160a060208501526117cd60a085018261102e565b9050604083015184820360408601526117e6828261102e565b91505060608301518482036060860152805115158252602081015161180e60208401826116fc565b50604081015161182160408401826116fc565b50606081015161183460608401826116fc565b5060808101516009811061184a5761184a6116e6565b8060808401525060a0810151905060c060a083015261186c60c083018261102e565b915050608083015184820360808601526118868282611710565b95945050505050565b8681526001600160a01b038616602082015260c0604082015260006118b760c08301876117ac565b82810360608401526118c9818761102e565b905082810360808401526118dd818661102e565b905082810360a08401526118f181856117ac565b9998505050505050505050565b60008060006060848603121561191357600080fd5b8351925060208401519150604084015167ffffffffffffffff81111561193857600080fd5b8401601f8101861361194957600080fd5b61158586825160208401611665565b6000821982111561197957634e487b7160e01b600052601160045260246000fd5b500190565b602081526000610fd5602083018461102e56fea2646970667358221220a1b823cbdc3af475d52edc7253ed678c6aeea8aba76ce4a197b391ff4c6e9d7064736f6c634300080b0033";

type KYXConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: KYXConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class KYX__factory extends ContractFactory {
  constructor(...args: KYXConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
  }

  deploy(
    tokenERC20: string,
    ancon: string,
    chain: BigNumberish,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<KYX> {
    return super.deploy(
      tokenERC20,
      ancon,
      chain,
      overrides || {}
    ) as Promise<KYX>;
  }
  getDeployTransaction(
    tokenERC20: string,
    ancon: string,
    chain: BigNumberish,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(
      tokenERC20,
      ancon,
      chain,
      overrides || {}
    );
  }
  attach(address: string): KYX {
    return super.attach(address) as KYX;
  }
  connect(signer: Signer): KYX__factory {
    return super.connect(signer) as KYX__factory;
  }
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): KYXInterface {
    return new utils.Interface(_abi) as KYXInterface;
  }
  static connect(address: string, signerOrProvider: Signer | Provider): KYX {
    return new Contract(address, _abi, signerOrProvider) as KYX;
  }
}