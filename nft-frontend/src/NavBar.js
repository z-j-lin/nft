import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import {Menu, Input} from 'semantic-ui-react';
import Web3 from 'web3'
import pkg from 'semantic-ui-react/package.json'
const backendurl = 'http://127.0.0.1:8081/';
class NavBar extends Component {
  constructor(props){
    super(props);
    this.state = {
      isLoggedIn: "0",
      web3: {},
      accounts: []
    };
  };
  /*
  getJwt = async () => {
    const { data } = await axios.get(`/jwt`);
    setJwt(data.token);
  }
*/
  buyHandler(){

  }

  logoutHandler(){

  }
  inventoryHandler(){

  }

  //function to login
  async loginHandler() {
    //run the login comp
    await this.loadWeb3();
    await this.loadBlockchainData();
    const account = this.state.accounts[0]
    var data
    this.state.web3.eth.personal.sign("hello", this.state.accounts[0]).then(
      signature =>{
        data = {"signature": signature, "account": account}
        console.log(data)
        const options = {
          method: 'POST',
          mode: 'cors',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
          //credentials: 'include',
          body: JSON.stringify(data)
        };
        console.log(options)
        fetch(backendurl+'login', options).then( response => {
          console.log(response.headers)
          //change isLoggedIn variable
          return response.json() 
        }).then(JsonResp => {
          console.log(JsonResp)
          if (JsonResp.Isloggedin === "1"){
            this.props.func(this.state.web3, this.state.accounts)
          }
          //sets the state and rerenders all pages with navbar
          this.setState({isLoggedIn: JsonResp.Isloggedin})
        })
      });
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
      case "logout":
        this.logoutHandler()
        break
      case "buy":
        this.buyHandler()
        break
      case "owned":
        this.inventoryHandler()
        break
      default:
    }
  }
  render() {
    const { activeItem } = this.state
    let nav;
    if(this.state.isLoggedIn == "1"){
      nav =( 
      <Menu >
        <Menu.Item as={Link} 
        to="/store"
        name='Store'
        active={activeItem === 'Store'}
        onClick={this.handleItemClick}
        >
        </Menu.Item>   
        <Menu.Item
          name='owned'
          active={activeItem === 'messages'}
          onClick={this.handleItemClick}
        />
        <Link to="/store"/>
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
        <Menu.Menu >
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
      <>
      {nav}
      </>
    );
  }
}
export default NavBar; 