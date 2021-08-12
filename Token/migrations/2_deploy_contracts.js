const CAToken = artifacts.require("CAToken");

module.exports = function (deployer) {
  deployer.deploy(CAToken, "CAToken", "CAT", "0x8833515f216Aff780d9ECAA2D01098fDbfc2520e");
};
