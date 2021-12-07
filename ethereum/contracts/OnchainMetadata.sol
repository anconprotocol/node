// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

contract OnchainMetadata {

  event AddOnchainMetadata(
    string name, 
    string description, 
    string indexed image, 
    string owner, 
    string parent, 
    bytes sources
    );

  event EncodeDagJson(
    string types,
    string values
  );

  constructor() {

  }

  function setOnchainMetadata(
    string memory name, 
    string memory description, 
    string memory image, 
    string memory owner, 
    string memory parent, 
    bytes memory sources
  ) public{



    emit AddOnchainMetadata(name, description, image, owner, parent, sources);

  }

  function encodeDagjsonBlock(
    string memory types,
    string memory values
  ) public returns (bool) {

    emit EncodeDagJson(types, values);

    return true;
  }
  //emit AddOnchainMetadata()
}

