const assert = require('http-assert')

const CAT = artifacts.require("./CAToken.sol")

require('chai')
  .use(require('chai-as-promised'))
  .should()

contract('CAToken', (accounts) => {
  let contract
  describe('deployment', async () => {
    it('deploys successfully', async () => {
      contract = await CAT.deployed()
      const address = contract.address
      console.log(address)
      assert.notEqual(address, '')
    })
  })

})