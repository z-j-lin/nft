import React from 'react';
import Web3 from 'web3'
import { Container, Form, Grid, Header, Message, Button } from 'semantic-ui-react';

class Login extends React.Component{
 
  constructor(props){
    super(props);
    this.state = {
      account: "",
      web3: {}
    }
  }
  //hook to trigger functions before page renders
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
    console.log(web3)
    //load account
    const accounts = await web3.eth.getAccounts()
    console.log(accounts)
    await web3.eth.personal.sign("hello", accounts[0]).then(console.log)
    //this is the first account in the wallet
    this.setState({account: accounts[0],
      web3: web3    
    })
  }
  
}

export default Login;