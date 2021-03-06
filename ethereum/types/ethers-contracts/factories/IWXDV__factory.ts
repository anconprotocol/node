/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */

import { Contract, Signer, utils } from "ethers";
import { Provider } from "@ethersproject/providers";
import type { IWXDV, IWXDVInterface } from "../IWXDV";

const _abi = [
  {
    inputs: [
      {
        internalType: "address",
        name: "sender",
        type: "address",
      },
      {
        internalType: "uint256",
        name: "newItemId",
        type: "uint256",
      },
      {
        internalType: "bytes32",
        name: "moniker",
        type: "bytes32",
      },
      {
        internalType: "bytes",
        name: "key",
        type: "bytes",
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
        name: "proof",
        type: "tuple",
      },
      {
        internalType: "bytes32",
        name: "hash",
        type: "bytes32",
      },
    ],
    name: "submitMintWithProof",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "payable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "sender",
        type: "address",
      },
      {
        internalType: "bytes32",
        name: "moniker",
        type: "bytes32",
      },
      {
        internalType: "bytes",
        name: "key",
        type: "bytes",
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
        name: "proof",
        type: "tuple",
      },
      {
        internalType: "bytes32",
        name: "hash",
        type: "bytes32",
      },
    ],
    name: "lockWithProof",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "payable",
    type: "function",
  },
  {
    inputs: [
      {
        internalType: "address",
        name: "sender",
        type: "address",
      },
      {
        internalType: "bytes32",
        name: "moniker",
        type: "bytes32",
      },
      {
        internalType: "bytes",
        name: "key",
        type: "bytes",
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
        name: "proof",
        type: "tuple",
      },
      {
        internalType: "bytes32",
        name: "hash",
        type: "bytes32",
      },
    ],
    name: "releaseWithProof",
    outputs: [
      {
        internalType: "uint256",
        name: "",
        type: "uint256",
      },
    ],
    stateMutability: "payable",
    type: "function",
  },
];

export class IWXDV__factory {
  static readonly abi = _abi;
  static createInterface(): IWXDVInterface {
    return new utils.Interface(_abi) as IWXDVInterface;
  }
  static connect(address: string, signerOrProvider: Signer | Provider): IWXDV {
    return new Contract(address, _abi, signerOrProvider) as IWXDV;
  }
}
