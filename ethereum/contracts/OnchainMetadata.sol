// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

contract OnchainMetadata {

  event AddOnchainMetadata(
    string name, 
    string description, 
    string image, 
    string owner, 
    string parent, 
    bytes sources
    );

  event EncodeDagJson(
    string path, 
    string hexdata
  );

  event EncodeDagCbor(
    string path, 
    string hexdata
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

function sum(uint x, uint y) public pure returns (uint){
  return 0;
}
  function encodeDagjsonBlock(
    string memory path,
    string memory hexdata
  ) public returns (bool) {

    emit EncodeDagJson(path, hexdata);

    return true;
  }
  //emit AddOnchainMetadata()
}

