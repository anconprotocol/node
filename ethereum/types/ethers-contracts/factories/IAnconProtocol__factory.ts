/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */

import { Contract, Signer, utils } from "ethers";
import { Provider } from "@ethersproject/providers";
import type {
  IAnconProtocol,
  IAnconProtocolInterface,
} from "../IAnconProtocol";

const _abi = [
  {
    inputs: [
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
        name: "proof",
        type: "tuple",
      },
    ],
    name: "submitPacketWithProof",
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
        name: "exProof",
        type: "tuple",
      },
    ],
    name: "verifyProof",
    outputs: [
      {
        internalType: "bool",
        name: "",
        type: "bool",
      },
    ],
    stateMutability: "view",
    type: "function",
  },
];

export class IAnconProtocol__factory {
  static readonly abi = _abi;
  static createInterface(): IAnconProtocolInterface {
    return new utils.Interface(_abi) as IAnconProtocolInterface;
  }
  static connect(
    address: string,
    signerOrProvider: Signer | Provider
  ): IAnconProtocol {
    return new Contract(address, _abi, signerOrProvider) as IAnconProtocol;
  }
}