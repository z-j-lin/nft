import React from 'react';
import Web3 from 'web3'

class Login {

  constructor(props){
    super(props)
  }
  
  client = {
    web3Provider: null,
    contracts: {},
    account: '0x0',
    loading: false,
    contractInstance: null,
    msg: '0x0',
    signature: '0x0',
  

  async componentWillMount() {
    await this.loadWeb3();
    await this.loadBlockchainData();
  }


  async loadWeb3() {
    if(window.ethereum){
      window.web3 = new Web3(window.ethereum);
      await window.ethereum.enable();
    }
    else if (window.web3){
      window.web3 = new Web3(window.web3.currentProvider)
    }
    else{
      window.alert('no ethereum browser detect, try installing metamask')
    }
  }

  async loadBlockchainData() {
    const web3 = window.web3;
    //load account
    const accounts = await web3.eth.getAccounts()
    //this is the first account in the wallet
    this.setState({account: accounts[0] })
  }

  render(){
    return(

    )
  }
  
  
}