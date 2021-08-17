import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import {Menu, Input} from 'semantic-ui-react';
import Web3 from 'web3'
import pkg from 'semantic-ui-react/package.json'
class NavBar extends Component {
  constructor(props){
    super(props);
    this.state = {
      isLoggedIn: false,
      web3: {},
      accounts: []
    };

  };

  //function to login
    async loginHandler() {
    //run the login comp
    await this.loadWeb3();
    await this.loadBlockchainData();
    this.state.web3.eth.personal.sign("hello", this.state.accounts[0]).then(
      signature =>{
      console.log("after login instantiation")
      const data = {signature, this.state.accounts[0]}
      const options = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      };

    fetch('/login', options).then( response => {
      console.log(response)
      //change isLoggedIn variable

    });

    })
    
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
  
  async  loadBlockchainData() {
    const web3 = window.web3;
    console.log(web3)
    //load account
    const accounts = await web3.eth.getAccounts()
    this.setState({
      web3: web3,
      accounts: accounts
    })  
  }
  
  handleItemClick = (e, { name }) => {
    this.setState({ activeItem: name })
    console.log(name)
    this.setState({ activeItem: "" })
    //case statement for individual button functions
    switch(name){
      case "login":
        this.loginHandler()
        break
      default:
    }
  }
  render() {
    const isLoggedIn = this.state.isLoggedIn;
    const { activeItem } = this.state
    
    let nav;
    if(isLoggedIn){
     nav =( 
     <Menu secondary>
        <Menu.Item
          name='buy'
          active={activeItem === 'home'}
          onClick={this.handleItemClick}
        />
        <Menu.Item
          name='owned'
          active={activeItem === 'messages'}
          onClick={this.handleItemClick}
        />
        <Menu.Menu position='right'>
          <Menu.Item
            name='logout'
            active={activeItem === 'logout'}
            onClick={this.handleItemClick}
          />
        </Menu.Menu>
      </Menu>
     )
    } else{
      nav = (
      <Menu widths={1}>
        <Menu.Menu position ='middle'>
        <Menu.Item 
          name='login'
          active={activeItem === 'login'}
          onClick={this.handleItemClick}
        />
        </Menu.Menu>
        
      </Menu>
      )
    }
    return(
      <div>
      {nav}
      </div>
    );
  }
}
export default NavBar; 