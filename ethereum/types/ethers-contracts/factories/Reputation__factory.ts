/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import { Signer, utils, Contract, ContractFactory, Overrides } from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type { Reputation, ReputationInterface } from "../Reputation";

const _abi = [
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "_from",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "_amount",
        type: "uint256",
      },
    ],
    name: "Burn",
    type: "event",
  },
  {
    anonymous: false,
    inputs: [
      {
        indexed: true,
        internalType: "address",
        name: "_to",
        type: "address",
      },
      {
        indexed: false,
        internalType: "uint256",
        name: "_amount",
        type: "uint256",
      },
    ],
    name: "Mint",
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
    inputs: [],
    name: "decimals",
    outputs: [
      {
        internalType: "uint8",
        name: "",
        type: "uint8",
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
        internalType: "address",
        name: "_user",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "_amount",
        type: "uint256",
      },
    ],
    name: "mint",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "_user",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "_amount",
        type: "uint256",
      },
    ],
    name: "burn",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "nonpayable",
    type: "function",
  },
  {
    inputs: [],
    name: "totalSupply",
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
        internalType: "address",
        name: "_owner",
        type: "address",
      },
    ],
    name: "balanceOf",
    outputs: [
      {
        internalType: "uint256",
        name: "balance",
        type: "uint256",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "uint256",
        name: "_blockNumber",
        type: "uint256",
      },
    ],
    name: "totalSupplyAt",
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
        internalType: "address",
        name: "_owner",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "_blockNumber",
        type: "uint256",
      },
    ],
    name: "balanceOfAt",
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
];

const _bytecode =
  "0x60806040526000805460ff60a01b1916600960a11b17905534801561002357600080fd5b5061002d33610032565b610082565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b610b28806100916000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c8063715018a611610076578063981b24d01161005b578063981b24d0146101835780639dc29fac14610196578063f2fde38b146101a957600080fd5b8063715018a61461015e5780638da5cb5b1461016857600080fd5b806340c10f19116100a757806340c10f19146101155780634ee2cd7e1461013857806370a082311461014b57600080fd5b806318160ddd146100c3578063313ce567146100de575b600080fd5b6100cb6101bc565b6040519081526020015b60405180910390f35b6000546101039074010000000000000000000000000000000000000000900460ff1681565b60405160ff90911681526020016100d5565b610128610123366004610a10565b6101cc565b60405190151581526020016100d5565b6100cb610146366004610a10565b610381565b6100cb610159366004610a3a565b610417565b610166610423565b005b6000546040516001600160a01b0390911681526020016100d5565b6100cb610191366004610a5c565b610489565b6101286101a4366004610a10565b6104df565b6101666101b7366004610a3a565b6105e4565b60006101c743610489565b905090565b600080546001600160a01b0316331461022c5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b60006102366101bc565b9050806102438482610a8b565b10156102915760405162461bcd60e51b815260206004820152601560248201527f746f74616c20737570706c79206f766572666c6f7700000000000000000000006044820152606401610223565b600061029c85610417565b9050806102a98582610a8b565b10156102f75760405162461bcd60e51b815260206004820152600f60248201527f62616c616365206f766572666c6f7700000000000000000000000000000000006044820152606401610223565b61030b60026103068685610a8b565b6106c6565b6001600160a01b0385166000908152600160205260409020610331906103068684610a8b565b846001600160a01b03167f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d41213968858560405161036c91815260200190565b60405180910390a26001925050505b92915050565b6001600160a01b03821660009081526001602052604081205415806103e157506001600160a01b038316600090815260016020526040812080548492906103ca576103ca610aa3565b6000918252602090912001546001600160801b0316115b156103ee5750600061037b565b6001600160a01b03831660009081526001602052604090206104109083610808565b905061037b565b600061037b8243610381565b6000546001600160a01b0316331461047d5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610223565b6104876000610991565b565b60025460009015806104c257508160026000815481106104ab576104ab610aa3565b6000918252602090912001546001600160801b0316115b156104cf57506000919050565b61037b600283610808565b919050565b600080546001600160a01b0316331461053a5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610223565b60006105446101bc565b905082600061055286610417565b905081811015610560578091505b61056f60026103068486610ab9565b6001600160a01b0386166000908152600160205260409020610595906103068484610ab9565b856001600160a01b03167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5836040516105d091815260200190565b60405180910390a250600195945050505050565b6000546001600160a01b0316331461063e5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610223565b6001600160a01b0381166106ba5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610223565b6106c381610991565b50565b80816001600160801b03161461071e5760405162461bcd60e51b815260206004820152601360248201527f72657075746174696f6e206f766572666c6f77000000000000000000000000006044820152606401610223565b81541580610760575081544390839061073990600190610ab9565b8154811061074957610749610aa3565b6000918252602090912001546001600160801b0316105b156107b7578154600090610775906001610a8b565b9050600083828154811061078b5761078b610aa3565b600091825260209091206001600160801b03858116600160801b02439190911617910155506108049050565b81546000906107c890600190610ab9565b905060008382815481106107de576107de610aa3565b600091825260209091200180546001600160801b03808616600160801b02911617905550505b5050565b81546000906108195750600061037b565b8254839061082990600190610ab9565b8154811061083957610839610aa3565b6000918252602090912001546001600160801b03168210610897578254839061086490600190610ab9565b8154811061087457610874610aa3565b600091825260209091200154600160801b90046001600160801b0316905061037b565b826000815481106108aa576108aa610aa3565b6000918252602090912001546001600160801b03168210156108ce5750600061037b565b825460009081906108e190600190610ab9565b90505b8181111561095a57600060026108fa8484610a8b565b610905906001610a8b565b61090f9190610ad0565b90508486828154811061092457610924610aa3565b6000918252602090912001546001600160801b03161161094657809250610954565b610951600182610ab9565b91505b506108e4565b84828154811061096c5761096c610aa3565b600091825260209091200154600160801b90046001600160801b031695945050505050565b600080546001600160a01b038381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b80356001600160a01b03811681146104da57600080fd5b60008060408385031215610a2357600080fd5b610a2c836109f9565b946020939093013593505050565b600060208284031215610a4c57600080fd5b610a55826109f9565b9392505050565b600060208284031215610a6e57600080fd5b5035919050565b634e487b7160e01b600052601160045260246000fd5b60008219821115610a9e57610a9e610a75565b500190565b634e487b7160e01b600052603260045260246000fd5b600082821015610acb57610acb610a75565b500390565b600082610aed57634e487b7160e01b600052601260045260246000fd5b50049056fea2646970667358221220b0a7b5e4608812218b3d7593504cba602143fa6a5cc12fd3a84813c19691263c64736f6c634300080c0033";

type ReputationConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: ReputationConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class Reputation__factory extends ContractFactory {
  constructor(...args: ReputationConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
  }

  deploy(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<Reputation> {
    return super.deploy(overrides || {}) as Promise<Reputation>;
  }
  getDeployTransaction(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(overrides || {});
  }
  attach(address: string): Reputation {
    return super.attach(address) as Reputation;
  }
  connect(signer: Signer): Reputation__factory {
    return super.connect(signer) as Reputation__factory;
  }
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): ReputationInterface {
    return new utils.Interface(_abi) as ReputationInterface;
  }
  static connect(
    address: string,
    signerOrProvider: Signer | Provider
  ): Reputation {
    return new Contract(address, _abi, signerOrProvider) as Reputation;
  }
}
