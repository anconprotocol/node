/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import {
  BaseContract,
  BigNumber,
  BigNumberish,
  BytesLike,
  CallOverrides,
  ContractTransaction,
  Overrides,
  PayableOverrides,
  PopulatedTransaction,
  Signer,
  utils,
} from "ethers";
import { FunctionFragment, Result, EventFragment } from "@ethersproject/abi";
import { Listener, Provider } from "@ethersproject/providers";
import { TypedEventFilter, TypedEvent, TypedListener, OnEvent } from "./common";

export type LeafOpStruct = {
  valid: boolean;
  hash: BigNumberish;
  prehash_key: BigNumberish;
  prehash_value: BigNumberish;
  len: BigNumberish;
  prefix: BytesLike;
};

export type LeafOpStructOutput = [
  boolean,
  number,
  number,
  number,
  number,
  string
] & {
  valid: boolean;
  hash: number;
  prehash_key: number;
  prehash_value: number;
  len: number;
  prefix: string;
};

export type InnerSpecStruct = {
  childOrder: BigNumberish[];
  childSize: BigNumberish;
  minPrefixLength: BigNumberish;
  maxPrefixLength: BigNumberish;
  emptyChild: BytesLike;
  hash: BigNumberish;
};

export type InnerSpecStructOutput = [
  BigNumber[],
  BigNumber,
  BigNumber,
  BigNumber,
  string,
  number
] & {
  childOrder: BigNumber[];
  childSize: BigNumber;
  minPrefixLength: BigNumber;
  maxPrefixLength: BigNumber;
  emptyChild: string;
  hash: number;
};

export type ProofSpecStruct = {
  leafSpec: LeafOpStruct;
  innerSpec: InnerSpecStruct;
  maxDepth: BigNumberish;
  minDepth: BigNumberish;
};

export type ProofSpecStructOutput = [
  LeafOpStructOutput,
  InnerSpecStructOutput,
  BigNumber,
  BigNumber
] & {
  leafSpec: LeafOpStructOutput;
  innerSpec: InnerSpecStructOutput;
  maxDepth: BigNumber;
  minDepth: BigNumber;
};

export type InnerOpStruct = {
  valid: boolean;
  hash: BigNumberish;
  prefix: BytesLike;
  suffix: BytesLike;
};

export type InnerOpStructOutput = [boolean, number, string, string] & {
  valid: boolean;
  hash: number;
  prefix: string;
  suffix: string;
};

export type ExistenceProofStruct = {
  valid: boolean;
  key: BytesLike;
  value: BytesLike;
  leaf: LeafOpStruct;
  path: InnerOpStruct[];
};

export type ExistenceProofStructOutput = [
  boolean,
  string,
  string,
  LeafOpStructOutput,
  InnerOpStructOutput[]
] & {
  valid: boolean;
  key: string;
  value: string;
  leaf: LeafOpStructOutput;
  path: InnerOpStructOutput[];
};

export interface AnconProtocolInterface extends utils.Interface {
  functions: {
    "ENROLL_PAYMENT()": FunctionFragment;
    "SUBMIT_PAYMENT()": FunctionFragment;
    "accountByAddrProofs(address)": FunctionFragment;
    "accountProofs(bytes)": FunctionFragment;
    "accountRegistrationFee()": FunctionFragment;
    "getIavlSpec()": FunctionFragment;
    "owner()": FunctionFragment;
    "proofs(bytes)": FunctionFragment;
    "protocolFee()": FunctionFragment;
    "relayNetworkHash()": FunctionFragment;
    "relayer()": FunctionFragment;
    "renounceOwnership()": FunctionFragment;
    "stablecoin()": FunctionFragment;
    "transferOwnership(address)": FunctionFragment;
    "verify((bool,bytes,bytes,(bool,uint8,uint8,uint8,uint8,bytes),(bool,uint8,bytes,bytes)[]),((bool,uint8,uint8,uint8,uint8,bytes),(uint256[],uint256,uint256,uint256,bytes,uint8),uint256,uint256),bytes,bytes,bytes)": FunctionFragment;
    "setPaymentToken(address)": FunctionFragment;
    "withdraw(address)": FunctionFragment;
    "withdrawToken(address,address)": FunctionFragment;
    "setProtocolFee(uint256)": FunctionFragment;
    "setAccountRegistrationFee(uint256)": FunctionFragment;
    "getProtocolHeader()": FunctionFragment;
    "getProof(bytes)": FunctionFragment;
    "hasProof(bytes)": FunctionFragment;
    "enrollL2Account(bytes,bytes,(bool,bytes,bytes,(bool,uint8,uint8,uint8,uint8,bytes),(bool,uint8,bytes,bytes)[]))": FunctionFragment;
    "updateProtocolHeader(bytes)": FunctionFragment;
    "submitPacketWithProof(bytes,bytes,(bool,bytes,bytes,(bool,uint8,uint8,uint8,uint8,bytes),(bool,uint8,bytes,bytes)[]))": FunctionFragment;
    "verifyProofWithKV(bytes,bytes,(bool,bytes,bytes,(bool,uint8,uint8,uint8,uint8,bytes),(bool,uint8,bytes,bytes)[]))": FunctionFragment;
  };

  encodeFunctionData(
    functionFragment: "ENROLL_PAYMENT",
    values?: undefined
  ): string;
  encodeFunctionData(
    functionFragment: "SUBMIT_PAYMENT",
    values?: undefined
  ): string;
  encodeFunctionData(
    functionFragment: "accountByAddrProofs",
    values: [string]
  ): string;
  encodeFunctionData(
    functionFragment: "accountProofs",
    values: [BytesLike]
  ): string;
  encodeFunctionData(
    functionFragment: "accountRegistrationFee",
    values?: undefined
  ): string;
  encodeFunctionData(
    functionFragment: "getIavlSpec",
    values?: undefined
  ): string;
  encodeFunctionData(functionFragment: "owner", values?: undefined): string;
  encodeFunctionData(functionFragment: "proofs", values: [BytesLike]): string;
  encodeFunctionData(
    functionFragment: "protocolFee",
    values?: undefined
  ): string;
  encodeFunctionData(
    functionFragment: "relayNetworkHash",
    values?: undefined
  ): string;
  encodeFunctionData(functionFragment: "relayer", values?: undefined): string;
  encodeFunctionData(
    functionFragment: "renounceOwnership",
    values?: undefined
  ): string;
  encodeFunctionData(
    functionFragment: "stablecoin",
    values?: undefined
  ): string;
  encodeFunctionData(
    functionFragment: "transferOwnership",
    values: [string]
  ): string;
  encodeFunctionData(
    functionFragment: "verify",
    values: [
      ExistenceProofStruct,
      ProofSpecStruct,
      BytesLike,
      BytesLike,
      BytesLike
    ]
  ): string;
  encodeFunctionData(
    functionFragment: "setPaymentToken",
    values: [string]
  ): string;
  encodeFunctionData(functionFragment: "withdraw", values: [string]): string;
  encodeFunctionData(
    functionFragment: "withdrawToken",
    values: [string, string]
  ): string;
  encodeFunctionData(
    functionFragment: "setProtocolFee",
    values: [BigNumberish]
  ): string;
  encodeFunctionData(
    functionFragment: "setAccountRegistrationFee",
    values: [BigNumberish]
  ): string;
  encodeFunctionData(
    functionFragment: "getProtocolHeader",
    values?: undefined
  ): string;
  encodeFunctionData(functionFragment: "getProof", values: [BytesLike]): string;
  encodeFunctionData(functionFragment: "hasProof", values: [BytesLike]): string;
  encodeFunctionData(
    functionFragment: "enrollL2Account",
    values: [BytesLike, BytesLike, ExistenceProofStruct]
  ): string;
  encodeFunctionData(
    functionFragment: "updateProtocolHeader",
    values: [BytesLike]
  ): string;
  encodeFunctionData(
    functionFragment: "submitPacketWithProof",
    values: [BytesLike, BytesLike, ExistenceProofStruct]
  ): string;
  encodeFunctionData(
    functionFragment: "verifyProofWithKV",
    values: [BytesLike, BytesLike, ExistenceProofStruct]
  ): string;

  decodeFunctionResult(
    functionFragment: "ENROLL_PAYMENT",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "SUBMIT_PAYMENT",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "accountByAddrProofs",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "accountProofs",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "accountRegistrationFee",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "getIavlSpec",
    data: BytesLike
  ): Result;
  decodeFunctionResult(functionFragment: "owner", data: BytesLike): Result;
  decodeFunctionResult(functionFragment: "proofs", data: BytesLike): Result;
  decodeFunctionResult(
    functionFragment: "protocolFee",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "relayNetworkHash",
    data: BytesLike
  ): Result;
  decodeFunctionResult(functionFragment: "relayer", data: BytesLike): Result;
  decodeFunctionResult(
    functionFragment: "renounceOwnership",
    data: BytesLike
  ): Result;
  decodeFunctionResult(functionFragment: "stablecoin", data: BytesLike): Result;
  decodeFunctionResult(
    functionFragment: "transferOwnership",
    data: BytesLike
  ): Result;
  decodeFunctionResult(functionFragment: "verify", data: BytesLike): Result;
  decodeFunctionResult(
    functionFragment: "setPaymentToken",
    data: BytesLike
  ): Result;
  decodeFunctionResult(functionFragment: "withdraw", data: BytesLike): Result;
  decodeFunctionResult(
    functionFragment: "withdrawToken",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "setProtocolFee",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "setAccountRegistrationFee",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "getProtocolHeader",
    data: BytesLike
  ): Result;
  decodeFunctionResult(functionFragment: "getProof", data: BytesLike): Result;
  decodeFunctionResult(functionFragment: "hasProof", data: BytesLike): Result;
  decodeFunctionResult(
    functionFragment: "enrollL2Account",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "updateProtocolHeader",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "submitPacketWithProof",
    data: BytesLike
  ): Result;
  decodeFunctionResult(
    functionFragment: "verifyProofWithKV",
    data: BytesLike
  ): Result;

  events: {
    "AccountRegistered(bool,bytes,bytes)": EventFragment;
    "HeaderUpdated(bytes)": EventFragment;
    "OwnershipTransferred(address,address)": EventFragment;
    "ProofPacketSubmitted(bytes,bytes)": EventFragment;
    "ServiceFeePaid(address,uint256)": EventFragment;
    "Withdrawn(address,uint256)": EventFragment;
  };

  getEvent(nameOrSignatureOrTopic: "AccountRegistered"): EventFragment;
  getEvent(nameOrSignatureOrTopic: "HeaderUpdated"): EventFragment;
  getEvent(nameOrSignatureOrTopic: "OwnershipTransferred"): EventFragment;
  getEvent(nameOrSignatureOrTopic: "ProofPacketSubmitted"): EventFragment;
  getEvent(nameOrSignatureOrTopic: "ServiceFeePaid"): EventFragment;
  getEvent(nameOrSignatureOrTopic: "Withdrawn"): EventFragment;
}

export type AccountRegisteredEvent = TypedEvent<
  [boolean, string, string],
  { enrolledStatus: boolean; key: string; value: string }
>;

export type AccountRegisteredEventFilter =
  TypedEventFilter<AccountRegisteredEvent>;

export type HeaderUpdatedEvent = TypedEvent<[string], { hash: string }>;

export type HeaderUpdatedEventFilter = TypedEventFilter<HeaderUpdatedEvent>;

export type OwnershipTransferredEvent = TypedEvent<
  [string, string],
  { previousOwner: string; newOwner: string }
>;

export type OwnershipTransferredEventFilter =
  TypedEventFilter<OwnershipTransferredEvent>;

export type ProofPacketSubmittedEvent = TypedEvent<
  [string, string],
  { key: string; packet: string }
>;

export type ProofPacketSubmittedEventFilter =
  TypedEventFilter<ProofPacketSubmittedEvent>;

export type ServiceFeePaidEvent = TypedEvent<
  [string, BigNumber],
  { from: string; fee: BigNumber }
>;

export type ServiceFeePaidEventFilter = TypedEventFilter<ServiceFeePaidEvent>;

export type WithdrawnEvent = TypedEvent<
  [string, BigNumber],
  { paymentAddress: string; amount: BigNumber }
>;

export type WithdrawnEventFilter = TypedEventFilter<WithdrawnEvent>;

export interface AnconProtocol extends BaseContract {
  connect(signerOrProvider: Signer | Provider | string): this;
  attach(addressOrName: string): this;
  deployed(): Promise<this>;

  interface: AnconProtocolInterface;

  queryFilter<TEvent extends TypedEvent>(
    event: TypedEventFilter<TEvent>,
    fromBlockOrBlockhash?: string | number | undefined,
    toBlock?: string | number | undefined
  ): Promise<Array<TEvent>>;

  listeners<TEvent extends TypedEvent>(
    eventFilter?: TypedEventFilter<TEvent>
  ): Array<TypedListener<TEvent>>;
  listeners(eventName?: string): Array<Listener>;
  removeAllListeners<TEvent extends TypedEvent>(
    eventFilter: TypedEventFilter<TEvent>
  ): this;
  removeAllListeners(eventName?: string): this;
  off: OnEvent<this>;
  on: OnEvent<this>;
  once: OnEvent<this>;
  removeListener: OnEvent<this>;

  functions: {
    ENROLL_PAYMENT(overrides?: CallOverrides): Promise<[string]>;

    SUBMIT_PAYMENT(overrides?: CallOverrides): Promise<[string]>;

    accountByAddrProofs(
      arg0: string,
      overrides?: CallOverrides
    ): Promise<[string]>;

    accountProofs(
      arg0: BytesLike,
      overrides?: CallOverrides
    ): Promise<[string]>;

    accountRegistrationFee(overrides?: CallOverrides): Promise<[BigNumber]>;

    getIavlSpec(overrides?: CallOverrides): Promise<[ProofSpecStructOutput]>;

    /**
     * Returns the address of the current owner.
     */
    owner(overrides?: CallOverrides): Promise<[string]>;

    proofs(arg0: BytesLike, overrides?: CallOverrides): Promise<[boolean]>;

    protocolFee(overrides?: CallOverrides): Promise<[BigNumber]>;

    relayNetworkHash(overrides?: CallOverrides): Promise<[string]>;

    relayer(overrides?: CallOverrides): Promise<[string]>;

    /**
     * Leaves the contract without owner. It will not be possible to call `onlyOwner` functions anymore. Can only be called by the current owner. NOTE: Renouncing ownership will leave the contract without an owner, thereby removing any functionality that is only available to the owner.
     */
    renounceOwnership(
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    stablecoin(overrides?: CallOverrides): Promise<[string]>;

    /**
     * Transfers ownership of the contract to a new account (`newOwner`). Can only be called by the current owner.
     */
    transferOwnership(
      newOwner: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    verify(
      proof: ExistenceProofStruct,
      spec: ProofSpecStruct,
      root: BytesLike,
      key: BytesLike,
      value: BytesLike,
      overrides?: CallOverrides
    ): Promise<[void]>;

    setPaymentToken(
      tokenAddress: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    withdraw(
      payee: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    withdrawToken(
      payee: string,
      erc20token: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    setProtocolFee(
      _fee: BigNumberish,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    setAccountRegistrationFee(
      _fee: BigNumberish,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    getProtocolHeader(overrides?: CallOverrides): Promise<[string]>;

    getProof(did: BytesLike, overrides?: CallOverrides): Promise<[string]>;

    hasProof(key: BytesLike, overrides?: CallOverrides): Promise<[boolean]>;

    enrollL2Account(
      key: BytesLike,
      did: BytesLike,
      proof: ExistenceProofStruct,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    updateProtocolHeader(
      rootHash: BytesLike,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    submitPacketWithProof(
      key: BytesLike,
      packet: BytesLike,
      proof: ExistenceProofStruct,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<ContractTransaction>;

    verifyProofWithKV(
      key: BytesLike,
      value: BytesLike,
      exProof: ExistenceProofStruct,
      overrides?: CallOverrides
    ): Promise<[boolean]>;
  };

  ENROLL_PAYMENT(overrides?: CallOverrides): Promise<string>;

  SUBMIT_PAYMENT(overrides?: CallOverrides): Promise<string>;

  accountByAddrProofs(arg0: string, overrides?: CallOverrides): Promise<string>;

  accountProofs(arg0: BytesLike, overrides?: CallOverrides): Promise<string>;

  accountRegistrationFee(overrides?: CallOverrides): Promise<BigNumber>;

  getIavlSpec(overrides?: CallOverrides): Promise<ProofSpecStructOutput>;

  /**
   * Returns the address of the current owner.
   */
  owner(overrides?: CallOverrides): Promise<string>;

  proofs(arg0: BytesLike, overrides?: CallOverrides): Promise<boolean>;

  protocolFee(overrides?: CallOverrides): Promise<BigNumber>;

  relayNetworkHash(overrides?: CallOverrides): Promise<string>;

  relayer(overrides?: CallOverrides): Promise<string>;

  /**
   * Leaves the contract without owner. It will not be possible to call `onlyOwner` functions anymore. Can only be called by the current owner. NOTE: Renouncing ownership will leave the contract without an owner, thereby removing any functionality that is only available to the owner.
   */
  renounceOwnership(
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  stablecoin(overrides?: CallOverrides): Promise<string>;

  /**
   * Transfers ownership of the contract to a new account (`newOwner`). Can only be called by the current owner.
   */
  transferOwnership(
    newOwner: string,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  verify(
    proof: ExistenceProofStruct,
    spec: ProofSpecStruct,
    root: BytesLike,
    key: BytesLike,
    value: BytesLike,
    overrides?: CallOverrides
  ): Promise<void>;

  setPaymentToken(
    tokenAddress: string,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  withdraw(
    payee: string,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  withdrawToken(
    payee: string,
    erc20token: string,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  setProtocolFee(
    _fee: BigNumberish,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  setAccountRegistrationFee(
    _fee: BigNumberish,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  getProtocolHeader(overrides?: CallOverrides): Promise<string>;

  getProof(did: BytesLike, overrides?: CallOverrides): Promise<string>;

  hasProof(key: BytesLike, overrides?: CallOverrides): Promise<boolean>;

  enrollL2Account(
    key: BytesLike,
    did: BytesLike,
    proof: ExistenceProofStruct,
    overrides?: PayableOverrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  updateProtocolHeader(
    rootHash: BytesLike,
    overrides?: Overrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  submitPacketWithProof(
    key: BytesLike,
    packet: BytesLike,
    proof: ExistenceProofStruct,
    overrides?: PayableOverrides & { from?: string | Promise<string> }
  ): Promise<ContractTransaction>;

  verifyProofWithKV(
    key: BytesLike,
    value: BytesLike,
    exProof: ExistenceProofStruct,
    overrides?: CallOverrides
  ): Promise<boolean>;

  callStatic: {
    ENROLL_PAYMENT(overrides?: CallOverrides): Promise<string>;

    SUBMIT_PAYMENT(overrides?: CallOverrides): Promise<string>;

    accountByAddrProofs(
      arg0: string,
      overrides?: CallOverrides
    ): Promise<string>;

    accountProofs(arg0: BytesLike, overrides?: CallOverrides): Promise<string>;

    accountRegistrationFee(overrides?: CallOverrides): Promise<BigNumber>;

    getIavlSpec(overrides?: CallOverrides): Promise<ProofSpecStructOutput>;

    /**
     * Returns the address of the current owner.
     */
    owner(overrides?: CallOverrides): Promise<string>;

    proofs(arg0: BytesLike, overrides?: CallOverrides): Promise<boolean>;

    protocolFee(overrides?: CallOverrides): Promise<BigNumber>;

    relayNetworkHash(overrides?: CallOverrides): Promise<string>;

    relayer(overrides?: CallOverrides): Promise<string>;

    /**
     * Leaves the contract without owner. It will not be possible to call `onlyOwner` functions anymore. Can only be called by the current owner. NOTE: Renouncing ownership will leave the contract without an owner, thereby removing any functionality that is only available to the owner.
     */
    renounceOwnership(overrides?: CallOverrides): Promise<void>;

    stablecoin(overrides?: CallOverrides): Promise<string>;

    /**
     * Transfers ownership of the contract to a new account (`newOwner`). Can only be called by the current owner.
     */
    transferOwnership(
      newOwner: string,
      overrides?: CallOverrides
    ): Promise<void>;

    verify(
      proof: ExistenceProofStruct,
      spec: ProofSpecStruct,
      root: BytesLike,
      key: BytesLike,
      value: BytesLike,
      overrides?: CallOverrides
    ): Promise<void>;

    setPaymentToken(
      tokenAddress: string,
      overrides?: CallOverrides
    ): Promise<void>;

    withdraw(payee: string, overrides?: CallOverrides): Promise<void>;

    withdrawToken(
      payee: string,
      erc20token: string,
      overrides?: CallOverrides
    ): Promise<void>;

    setProtocolFee(
      _fee: BigNumberish,
      overrides?: CallOverrides
    ): Promise<void>;

    setAccountRegistrationFee(
      _fee: BigNumberish,
      overrides?: CallOverrides
    ): Promise<void>;

    getProtocolHeader(overrides?: CallOverrides): Promise<string>;

    getProof(did: BytesLike, overrides?: CallOverrides): Promise<string>;

    hasProof(key: BytesLike, overrides?: CallOverrides): Promise<boolean>;

    enrollL2Account(
      key: BytesLike,
      did: BytesLike,
      proof: ExistenceProofStruct,
      overrides?: CallOverrides
    ): Promise<boolean>;

    updateProtocolHeader(
      rootHash: BytesLike,
      overrides?: CallOverrides
    ): Promise<void>;

    submitPacketWithProof(
      key: BytesLike,
      packet: BytesLike,
      proof: ExistenceProofStruct,
      overrides?: CallOverrides
    ): Promise<boolean>;

    verifyProofWithKV(
      key: BytesLike,
      value: BytesLike,
      exProof: ExistenceProofStruct,
      overrides?: CallOverrides
    ): Promise<boolean>;
  };

  filters: {
    "AccountRegistered(bool,bytes,bytes)"(
      enrolledStatus?: null,
      key?: null,
      value?: null
    ): AccountRegisteredEventFilter;
    AccountRegistered(
      enrolledStatus?: null,
      key?: null,
      value?: null
    ): AccountRegisteredEventFilter;

    "HeaderUpdated(bytes)"(hash?: null): HeaderUpdatedEventFilter;
    HeaderUpdated(hash?: null): HeaderUpdatedEventFilter;

    "OwnershipTransferred(address,address)"(
      previousOwner?: string | null,
      newOwner?: string | null
    ): OwnershipTransferredEventFilter;
    OwnershipTransferred(
      previousOwner?: string | null,
      newOwner?: string | null
    ): OwnershipTransferredEventFilter;

    "ProofPacketSubmitted(bytes,bytes)"(
      key?: null,
      packet?: null
    ): ProofPacketSubmittedEventFilter;
    ProofPacketSubmitted(
      key?: null,
      packet?: null
    ): ProofPacketSubmittedEventFilter;

    "ServiceFeePaid(address,uint256)"(
      from?: string | null,
      fee?: null
    ): ServiceFeePaidEventFilter;
    ServiceFeePaid(from?: string | null, fee?: null): ServiceFeePaidEventFilter;

    "Withdrawn(address,uint256)"(
      paymentAddress?: string | null,
      amount?: null
    ): WithdrawnEventFilter;
    Withdrawn(
      paymentAddress?: string | null,
      amount?: null
    ): WithdrawnEventFilter;
  };

  estimateGas: {
    ENROLL_PAYMENT(overrides?: CallOverrides): Promise<BigNumber>;

    SUBMIT_PAYMENT(overrides?: CallOverrides): Promise<BigNumber>;

    accountByAddrProofs(
      arg0: string,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    accountProofs(
      arg0: BytesLike,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    accountRegistrationFee(overrides?: CallOverrides): Promise<BigNumber>;

    getIavlSpec(overrides?: CallOverrides): Promise<BigNumber>;

    /**
     * Returns the address of the current owner.
     */
    owner(overrides?: CallOverrides): Promise<BigNumber>;

    proofs(arg0: BytesLike, overrides?: CallOverrides): Promise<BigNumber>;

    protocolFee(overrides?: CallOverrides): Promise<BigNumber>;

    relayNetworkHash(overrides?: CallOverrides): Promise<BigNumber>;

    relayer(overrides?: CallOverrides): Promise<BigNumber>;

    /**
     * Leaves the contract without owner. It will not be possible to call `onlyOwner` functions anymore. Can only be called by the current owner. NOTE: Renouncing ownership will leave the contract without an owner, thereby removing any functionality that is only available to the owner.
     */
    renounceOwnership(
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    stablecoin(overrides?: CallOverrides): Promise<BigNumber>;

    /**
     * Transfers ownership of the contract to a new account (`newOwner`). Can only be called by the current owner.
     */
    transferOwnership(
      newOwner: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    verify(
      proof: ExistenceProofStruct,
      spec: ProofSpecStruct,
      root: BytesLike,
      key: BytesLike,
      value: BytesLike,
      overrides?: CallOverrides
    ): Promise<BigNumber>;

    setPaymentToken(
      tokenAddress: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    withdraw(
      payee: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    withdrawToken(
      payee: string,
      erc20token: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    setProtocolFee(
      _fee: BigNumberish,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    setAccountRegistrationFee(
      _fee: BigNumberish,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    getProtocolHeader(overrides?: CallOverrides): Promise<BigNumber>;

    getProof(did: BytesLike, overrides?: CallOverrides): Promise<BigNumber>;

    hasProof(key: BytesLike, overrides?: CallOverrides): Promise<BigNumber>;

    enrollL2Account(
      key: BytesLike,
      did: BytesLike,
      proof: ExistenceProofStruct,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    updateProtocolHeader(
      rootHash: BytesLike,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    submitPacketWithProof(
      key: BytesLike,
      packet: BytesLike,
      proof: ExistenceProofStruct,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<BigNumber>;

    verifyProofWithKV(
      key: BytesLike,
      value: BytesLike,
      exProof: ExistenceProofStruct,
      overrides?: CallOverrides
    ): Promise<BigNumber>;
  };

  populateTransaction: {
    ENROLL_PAYMENT(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    SUBMIT_PAYMENT(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    accountByAddrProofs(
      arg0: string,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    accountProofs(
      arg0: BytesLike,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    accountRegistrationFee(
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    getIavlSpec(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    /**
     * Returns the address of the current owner.
     */
    owner(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    proofs(
      arg0: BytesLike,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    protocolFee(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    relayNetworkHash(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    relayer(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    /**
     * Leaves the contract without owner. It will not be possible to call `onlyOwner` functions anymore. Can only be called by the current owner. NOTE: Renouncing ownership will leave the contract without an owner, thereby removing any functionality that is only available to the owner.
     */
    renounceOwnership(
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    stablecoin(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    /**
     * Transfers ownership of the contract to a new account (`newOwner`). Can only be called by the current owner.
     */
    transferOwnership(
      newOwner: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    verify(
      proof: ExistenceProofStruct,
      spec: ProofSpecStruct,
      root: BytesLike,
      key: BytesLike,
      value: BytesLike,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    setPaymentToken(
      tokenAddress: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    withdraw(
      payee: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    withdrawToken(
      payee: string,
      erc20token: string,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    setProtocolFee(
      _fee: BigNumberish,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    setAccountRegistrationFee(
      _fee: BigNumberish,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    getProtocolHeader(overrides?: CallOverrides): Promise<PopulatedTransaction>;

    getProof(
      did: BytesLike,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    hasProof(
      key: BytesLike,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;

    enrollL2Account(
      key: BytesLike,
      did: BytesLike,
      proof: ExistenceProofStruct,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    updateProtocolHeader(
      rootHash: BytesLike,
      overrides?: Overrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    submitPacketWithProof(
      key: BytesLike,
      packet: BytesLike,
      proof: ExistenceProofStruct,
      overrides?: PayableOverrides & { from?: string | Promise<string> }
    ): Promise<PopulatedTransaction>;

    verifyProofWithKV(
      key: BytesLike,
      value: BytesLike,
      exProof: ExistenceProofStruct,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;
  };
}
