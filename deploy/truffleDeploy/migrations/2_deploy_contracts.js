const SimpleTest = artifacts.require("./SimpleTest.sol");

module.exports = function(deployer) {
  deployer.deploy(SimpleTest);
};

//truffle migrate --network development
