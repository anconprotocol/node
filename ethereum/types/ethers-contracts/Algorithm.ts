/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import {
  BaseContract,
  BigNumber,
  BytesLike,
  CallOverrides,
  PopulatedTransaction,
  Signer,
  utils,
} from "ethers";
import { FunctionFragment, Result } from "@ethersproject/abi";
import { Listener, Provider } from "@ethersproject/providers";
import { TypedEventFilter, TypedEvent, TypedListener, OnEvent } from "./common";

export interface AlgorithmInterface extends utils.Interface {
  functions: {
    "verify(bytes,bytes,bytes)": FunctionFragment;
  };

  encodeFunctionData(
    functionFragment: "verify",
    values: [BytesLike, BytesLike, BytesLike]
  ): string;

  decodeFunctionResult(functionFragment: "verify", data: BytesLike): Result;

  events: {};
}

export interface Algorithm extends BaseContract {
  connect(signerOrProvider: Signer | Provider | string): this;
  attach(addressOrName: string): this;
  deployed(): Promise<this>;

  interface: AlgorithmInterface;

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
    /**
     * Verifies a signature.
     * @param data The signed data to verify.
     * @param key The public key to verify with.
     * @param signature The signature to verify.
     */
    verify(
      key: BytesLike,
      data: BytesLike,
      signature: BytesLike,
      overrides?: CallOverrides
    ): Promise<[boolean]>;
  };

  /**
   * Verifies a signature.
   * @param data The signed data to verify.
   * @param key The public key to verify with.
   * @param signature The signature to verify.
   */
  verify(
    key: BytesLike,
    data: BytesLike,
    signature: BytesLike,
    overrides?: CallOverrides
  ): Promise<boolean>;

  callStatic: {
    /**
     * Verifies a signature.
     * @param data The signed data to verify.
     * @param key The public key to verify with.
     * @param signature The signature to verify.
     */
    verify(
      key: BytesLike,
      data: BytesLike,
      signature: BytesLike,
      overrides?: CallOverrides
    ): Promise<boolean>;
  };

  filters: {};

  estimateGas: {
    /**
     * Verifies a signature.
     * @param data The signed data to verify.
     * @param key The public key to verify with.
     * @param signature The signature to verify.
     */
    verify(
      key: BytesLike,
      data: BytesLike,
      signature: BytesLike,
      overrides?: CallOverrides
    ): Promise<BigNumber>;
  };

  populateTransaction: {
    /**
     * Verifies a signature.
     * @param data The signed data to verify.
     * @param key The public key to verify with.
     * @param signature The signature to verify.
     */
    verify(
      key: BytesLike,
      data: BytesLike,
      signature: BytesLike,
      overrides?: CallOverrides
    ): Promise<PopulatedTransaction>;
  };
}