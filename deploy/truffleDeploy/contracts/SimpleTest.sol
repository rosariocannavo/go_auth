// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleTest {
  constructor() public {
  }

   uint256 private _storedData;

    // Function to set the stored value
    function set(uint256 x) public {
        _storedData = x;
    }

    // Function to get the stored value
    function get() public view returns (uint256) {
        return _storedData;
    }
}
