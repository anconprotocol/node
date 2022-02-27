/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import { Signer, utils, Contract, ContractFactory, Overrides } from "ethers";
import { Provider, TransactionRequest } from "@ethersproject/providers";
import type { ICS23, ICS23Interface } from "../ICS23";

const _abi = [
  {
    inputs: [],
    name: "getIavlSpec",
    outputs: [
      {
        components: [
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
            name: "leafSpec",
            type: "tuple",
          },
          {
            components: [
              {
                internalType: "uint256[]",
                name: "childOrder",
                type: "uint256[]",
              },
              {
                internalType: "uint256",
                name: "childSize",
                type: "uint256",
              },
              {
                internalType: "uint256",
                name: "minPrefixLength",
                type: "uint256",
              },
              {
                internalType: "uint256",
                name: "maxPrefixLength",
                type: "uint256",
              },
              {
                internalType: "bytes",
                name: "emptyChild",
                type: "bytes",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "hash",
                type: "uint8",
              },
            ],
            internalType: "struct Ics23Helper.InnerSpec",
            name: "innerSpec",
            type: "tuple",
          },
          {
            internalType: "uint256",
            name: "maxDepth",
            type: "uint256",
          },
          {
            internalType: "uint256",
            name: "minDepth",
            type: "uint256",
          },
        ],
        internalType: "struct Ics23Helper.ProofSpec",
        name: "",
        type: "tuple",
      },
    ],
    stateMutability: "pure",
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
        name: "proof",
        type: "tuple",
      },
      {
        components: [
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
            name: "leafSpec",
            type: "tuple",
          },
          {
            components: [
              {
                internalType: "uint256[]",
                name: "childOrder",
                type: "uint256[]",
              },
              {
                internalType: "uint256",
                name: "childSize",
                type: "uint256",
              },
              {
                internalType: "uint256",
                name: "minPrefixLength",
                type: "uint256",
              },
              {
                internalType: "uint256",
                name: "maxPrefixLength",
                type: "uint256",
              },
              {
                internalType: "bytes",
                name: "emptyChild",
                type: "bytes",
              },
              {
                internalType: "enum Ics23Helper.HashOp",
                name: "hash",
                type: "uint8",
              },
            ],
            internalType: "struct Ics23Helper.InnerSpec",
            name: "innerSpec",
            type: "tuple",
          },
          {
            internalType: "uint256",
            name: "maxDepth",
            type: "uint256",
          },
          {
            internalType: "uint256",
            name: "minDepth",
            type: "uint256",
          },
        ],
        internalType: "struct Ics23Helper.ProofSpec",
        name: "spec",
        type: "tuple",
      },
      {
        internalType: "bytes",
        name: "root",
        type: "bytes",
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
    ],
    name: "verify",
    outputs: [],
    stateMutability: "pure",
    type: "function",
  },
];

const _bytecode =
  "0x608060405234801561001057600080fd5b50611b88806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806327dcd78c1461003b578063b0d264e714610059575b600080fd5b61004361006e565b60405161005091906112fe565b60405180910390f35b61006c610067366004611913565b610182565b005b610076611144565b61007e611144565b6040805160028082526060820183526000926020830190803683370190505090506000816000815181106100b4576100b46119e5565b6020026020010181815250506001816001815181106100d5576100d56119e5565b6020026020010181815250506040518060c0016040528060011515815260200160016005811115610108576101086111e7565b815260200160008152602001600181526020016001815260408051808201825260018082526000602083810182905294850192909252938652815160c081018352858152602181850152600481840152600c60608201528251938401909252825260808101919091529060a0820152602083015250919050565b61018c85856102f6565b61019a85602001518361041e565b6101eb5760405162461bcd60e51b815260206004820181905260248201527f50726f7669646564206b657920646f65736e2774206d617463682070726f6f6660448201526064015b60405180910390fd5b6101f985604001518261041e565b61026b5760405162461bcd60e51b815260206004820152602260248201527f50726f76696465642076616c756520646f65736e2774206d617463682070726f60448201527f6f6600000000000000000000000000000000000000000000000000000000000060648201526084016101e2565b61027d6102778661044b565b8461041e565b6102ef5760405162461bcd60e51b815260206004820152602c60248201527f43616c63756c636174656420726f6f7420646f65736e2774206d61746368207060448201527f726f766964656420726f6f74000000000000000000000000000000000000000060648201526084016101e2565b5050505050565b6103048260600151826104ba565b6060810151158061031e5750806060015182608001515110155b61036a5760405162461bcd60e51b815260206004820152601860248201527f496e6e65724f707320646570746820746f6f2073686f7274000000000000000060448201526064016101e2565b604081015115806103845750806040015182608001515110155b6103d05760405162461bcd60e51b815260206004820152601860248201527f496e6e65724f707320646570746820746f6f2073686f7274000000000000000060448201526064016101e2565b60005b82608001515181101561041957610407836080015182815181106103f9576103f96119e5565b602002602001015183610708565b8061041181611a11565b9150506103d3565b505050565b6000815183511461043157506000610445565b508151602082810182902090840191909120145b92915050565b606060006104668360600151846020015185604001516108d3565b905060005b8360800151518110156104b35761049f84608001518281518110610491576104916119e5565b6020026020010151836109fd565b9150806104ab81611a11565b91505061046b565b5092915050565b80516020015160058111156104d1576104d16111e7565b826020015160058111156104e7576104e76111e7565b146105345760405162461bcd60e51b815260206004820152601160248201527f556e657870656374656420486173684f7000000000000000000000000000000060448201526064016101e2565b805160400151600581111561054b5761054b6111e7565b82604001516005811115610561576105616111e7565b146105ae5760405162461bcd60e51b815260206004820152601560248201527f556e657870656374656420507265686173684b6579000000000000000000000060448201526064016101e2565b80516060015160058111156105c5576105c56111e7565b826060015160058111156105db576105db6111e7565b146106285760405162461bcd60e51b815260206004820152601560248201527f556e657870656374656420507265686173684b6579000000000000000000000060448201526064016101e2565b805160800151600881111561063f5761063f6111e7565b82608001516008811115610655576106556111e7565b146106a25760405162461bcd60e51b815260206004820152601a60248201527f556e657870656374656c65616653706563204c656e6774684f7000000000000060448201526064016101e2565b6106b88260a00151826000015160a00151610aac565b6107045760405162461bcd60e51b815260206004820152601760248201527f4c6561664f704c69623a2077726f6e672070726566697800000000000000000060448201526064016101e2565b5050565b805160200151600581111561071f5761071f6111e7565b82602001516005811115610735576107356111e7565b146107825760405162461bcd60e51b815260206004820152601160248201527f556e657870656374656420486173684f7000000000000000000000000000000060448201526064016101e2565b6107988260400151826000015160a00151610aac565b156107e55760405162461bcd60e51b815260206004820152601860248201527f496e6e65724f704c69623a2077726f6e6720707265666978000000000000000060448201526064016101e2565b80602001516040015182604001515110156108425760405162461bcd60e51b815260206004820152601860248201527f496e6e65724f702070726566697820746f6f2073686f7274000000000000000060448201526064016101e2565b602080820151908101519051516000919061085f90600190611a2c565b6108699190611a43565b90508082602001516060015161087f9190611a62565b83604001515111156104195760405162461bcd60e51b815260206004820152601860248201527f496e6e65724f702070726566697820746f6f2073686f7274000000000000000060448201526064016101e2565b606060008351116109265760405162461bcd60e51b815260206004820152601160248201527f4c656166206f70206e65656473206b657900000000000000000000000000000060448201526064016101e2565b60008251116109775760405162461bcd60e51b815260206004820152601360248201527f4c656166206f70206e656564732076616c75650000000000000000000000000060448201526064016101e2565b60008460a001516109918660400151876080015187610b6f565b6040516020016109a2929190611a7a565b6040516020818303038152906040526109c48660600151876080015186610b6f565b6040516020016109d5929190611a7a565b60405160208183030381529060405290506109f4856020015182610b8c565b95945050505050565b60606000825111610a505760405162461bcd60e51b815260206004820152601a60248201527f496e6e6572206f70206e65656473206368696c642076616c756500000000000060448201526064016101e2565b610aa58360200151846040015184604051602001610a6f929190611a7a565b60408051601f19818403018152908290526060870151610a9192602001611a7a565b604051602081830303815290604052610b8c565b9392505050565b6000815183511015610ac057506000610445565b60005b8251811015610b6557828181518110610ade57610ade6119e5565b602001015160f81c60f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916848281518110610b1d57610b1d6119e5565b01602001517fff000000000000000000000000000000000000000000000000000000000000001614610b53576000915050610445565b80610b5d81611a11565b915050610ac3565b5060019392505050565b6060610b8483610b7f8685610daa565b610dd7565b949350505050565b60606001836005811115610ba257610ba26111e7565b1415610c0757610c00600283604051610bbb9190611aa9565b602060405180830381855afa158015610bd8573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190610bfb9190611ac5565b610f6b565b9050610445565b6002836005811115610c1b57610c1b6111e7565b1415610c695760405162461bcd60e51b815260206004820152601660248201527f534841353132206e6f7420696d706c656d656e7465640000000000000000000060448201526064016101e2565b6004836005811115610c7d57610c7d6111e7565b1415610cd157610c00600383604051610c969190611aa9565b602060405180830381855afa158015610cb3573d6000803e3d6000fd5b5050506040515160601b6bffffffffffffffffffffffff1916610f6b565b6005836005811115610ce557610ce56111e7565b1415610d62576000600283604051610cfd9190611aa9565b602060405180830381855afa158015610d1a573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190610d3d9190611ac5565b9050610d5a6003610d4d83610f6b565b604051610c969190611aa9565b915050610445565b60405162461bcd60e51b815260206004820152601260248201527f556e737570706f7274656420686173686f70000000000000000000000000000060448201526064016101e2565b60606000836005811115610dc057610dc06111e7565b1415610dcd575080610445565b610aa58383610b8c565b60606000836008811115610ded57610ded6111e7565b1415610dfa575080610445565b6001836008811115610e0e57610e0e6111e7565b1415610e4657610e1e8251610f95565b82604051602001610e30929190611a7a565b6040516020818303038152906040529050610445565b6007836008811115610e5a57610e5a6111e7565b1415610eb8578151602014610eb15760405162461bcd60e51b815260206004820152601160248201527f457870656374656420333220627974657300000000000000000000000000000060448201526064016101e2565b5080610445565b6008836008811115610ecc57610ecc6111e7565b1415610f23578151604014610eb15760405162461bcd60e51b815260206004820152601160248201527f457870656374656420363420627974657300000000000000000000000000000060448201526064016101e2565b60405162461bcd60e51b815260206004820152601460248201527f556e737570706f72746564206c656e6774686f7000000000000000000000000060448201526064016101e2565b60408051602080825281830190925260609160208201818036833750505060208101929092525090565b60608160015b607f8267ffffffffffffffff161115610fd35760078267ffffffffffffffff16901c9150600181610fcc9190611ade565b9050610f9b565b60008167ffffffffffffffff1667ffffffffffffffff811115610ff857610ff86113c3565b6040519080825280601f01601f191660200182016040528015611022576020820181803683370190505b50905084925060005b8267ffffffffffffffff168167ffffffffffffffff1610156110c05783607f1660801760f81b828267ffffffffffffffff168151811061106d5761106d6119e5565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535060078467ffffffffffffffff16901c935080806110b890611b01565b91505061102b565b507f7f00000000000000000000000000000000000000000000000000000000000000816110ee600185611b29565b67ffffffffffffffff1681518110611108576111086119e5565b0160200180519091167fff000000000000000000000000000000000000000000000000000000000000001690600082901a905350949350505050565b60408051610140810190915260006080820181815260a0830182905260c0830182905260e08301829052610100830191909152606061012083015281526020810161118d6111a1565b815260200160008152602001600081525090565b6040518060c001604052806060815260200160008152602001600081526020016000815260200160608152602001600060058111156111e2576111e26111e7565b905290565b634e487b7160e01b600052602160045260246000fd5b6006811061120d5761120d6111e7565b9052565b60005b8381101561122c578181015183820152602001611214565b8381111561123b576000848401525b50505050565b60008151808452611259816020860160208601611211565b601f01601f19169290920160200192915050565b805160c080845281519084018190526000916020919082019060e0860190845b818110156112a95783518352928401929184019160010161128d565b5050828501518387015260408501516040870152606085015160608701526080850151925085810360808701526112e08184611241565b9250505060a08301516112f660a08601826111fd565b509392505050565b6020815260008251608060208401528051151560a0840152602081015161132860c08501826111fd565b50604081015161133b60e08501826111fd565b50606081015161134f6101008501826111fd565b50608081015160098110611365576113656111e7565b61012084015260a0015160c0610140840152611385610160840182611241565b90506020840151601f198483030160408501526113a2828261126d565b91505060408401516060840152606084015160808401528091505092915050565b634e487b7160e01b600052604160045260246000fd5b60405160c0810167ffffffffffffffff811182821017156113fc576113fc6113c3565b60405290565b6040516080810167ffffffffffffffff811182821017156113fc576113fc6113c3565b60405160a0810167ffffffffffffffff811182821017156113fc576113fc6113c3565b604051601f8201601f1916810167ffffffffffffffff81118282101715611471576114716113c3565b604052919050565b8035801515811461148957600080fd5b919050565b600082601f83011261149f57600080fd5b813567ffffffffffffffff8111156114b9576114b96113c3565b6114cc601f8201601f1916602001611448565b8181528460208386010111156114e157600080fd5b816020850160208301376000918101602001919091529392505050565b80356006811061148957600080fd5b600060c0828403121561151f57600080fd5b6115276113d9565b905061153282611479565b8152611540602083016114fe565b6020820152611551604083016114fe565b6040820152611562606083016114fe565b606082015260808201356009811061157957600080fd5b608082015260a082013567ffffffffffffffff81111561159857600080fd5b6115a48482850161148e565b60a08301525092915050565b600067ffffffffffffffff8211156115ca576115ca6113c3565b5060051b60200190565b600082601f8301126115e557600080fd5b813560206115fa6115f5836115b0565b611448565b82815260059290921b8401810191818101908684111561161957600080fd5b8286015b848110156116dd57803567ffffffffffffffff8082111561163e5760008081fd5b908801906080828b03601f19018113156116585760008081fd5b611660611402565b61166b888501611479565b8152604061167a8186016114fe565b89830152606080860135858111156116925760008081fd5b6116a08f8c838a010161148e565b84840152509285013592848411156116ba57600091508182fd5b6116c88e8b8689010161148e565b9083015250865250505091830191830161161d565b509695505050505050565b600060a082840312156116fa57600080fd5b611702611425565b905061170d82611479565b8152602082013567ffffffffffffffff8082111561172a57600080fd5b6117368583860161148e565b6020840152604084013591508082111561174f57600080fd5b61175b8583860161148e565b6040840152606084013591508082111561177457600080fd5b6117808583860161150d565b6060840152608084013591508082111561179957600080fd5b506117a6848285016115d4565b60808301525092915050565b600082601f8301126117c357600080fd5b813560206117d36115f5836115b0565b82815260059290921b840181019181810190868411156117f257600080fd5b8286015b848110156116dd57803583529183019183016117f6565b60006080828403121561181f57600080fd5b611827611402565b9050813567ffffffffffffffff8082111561184157600080fd5b61184d8583860161150d565b8352602084013591508082111561186357600080fd5b9083019060c0828603121561187757600080fd5b61187f6113d9565b82358281111561188e57600080fd5b61189a878286016117b2565b8252506020830135602082015260408301356040820152606083013560608201526080830135828111156118cd57600080fd5b6118d98782860161148e565b6080830152506118eb60a084016114fe565b60a0820152806020850152505050604082013560408201526060820135606082015292915050565b600080600080600060a0868803121561192b57600080fd5b853567ffffffffffffffff8082111561194357600080fd5b61194f89838a016116e8565b9650602088013591508082111561196557600080fd5b61197189838a0161180d565b9550604088013591508082111561198757600080fd5b61199389838a0161148e565b945060608801359150808211156119a957600080fd5b6119b589838a0161148e565b935060808801359150808211156119cb57600080fd5b506119d88882890161148e565b9150509295509295909350565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600019821415611a2557611a256119fb565b5060010190565b600082821015611a3e57611a3e6119fb565b500390565b6000816000190483118215151615611a5d57611a5d6119fb565b500290565b60008219821115611a7557611a756119fb565b500190565b60008351611a8c818460208801611211565b835190830190611aa0818360208801611211565b01949350505050565b60008251611abb818460208701611211565b9190910192915050565b600060208284031215611ad757600080fd5b5051919050565b600067ffffffffffffffff808316818516808303821115611aa057611aa06119fb565b600067ffffffffffffffff80831681811415611b1f57611b1f6119fb565b6001019392505050565b600067ffffffffffffffff83811690831681811015611b4a57611b4a6119fb565b03939250505056fea2646970667358221220257fb1ddf16a3767754cfad06ce2b9f327b3f4726284ad7949757c5522c9b8b264736f6c634300080c0033";

type ICS23ConstructorParams =
  | [signer?: Signer]
  | ConstructorParameters<typeof ContractFactory>;

const isSuperArgs = (
  xs: ICS23ConstructorParams
): xs is ConstructorParameters<typeof ContractFactory> => xs.length > 1;

export class ICS23__factory extends ContractFactory {
  constructor(...args: ICS23ConstructorParams) {
    if (isSuperArgs(args)) {
      super(...args);
    } else {
      super(_abi, _bytecode, args[0]);
    }
  }

  deploy(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ICS23> {
    return super.deploy(overrides || {}) as Promise<ICS23>;
  }
  getDeployTransaction(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): TransactionRequest {
    return super.getDeployTransaction(overrides || {});
  }
  attach(address: string): ICS23 {
    return super.attach(address) as ICS23;
  }
  connect(signer: Signer): ICS23__factory {
    return super.connect(signer) as ICS23__factory;
  }
  static readonly bytecode = _bytecode;
  static readonly abi = _abi;
  static createInterface(): ICS23Interface {
    return new utils.Interface(_abi) as ICS23Interface;
  }
  static connect(address: string, signerOrProvider: Signer | Provider): ICS23 {
    return new Contract(address, _abi, signerOrProvider) as ICS23;
  }
}
