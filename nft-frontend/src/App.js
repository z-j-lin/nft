import React, {Component} from 'react';
import Web3 from 'web3'
import './App.css';
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom';
import Login from './pages/Login';
import NavBar from './NavBar';
//import detectEthereumProvier from "@metamask/detect-provider"

class App extends Component {
  
  constructor(props){
    super(props);
    this.state = {
      account: ""
    }
  }

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

  render() {
    return (
      <Router>
        <div className="App">
          <NavBar account = {this.state.account}/>
          <div id = "page-body">
            {/*homepage */}
            <Route path = "/Login" component={Login} exact/>
          </div>
        </div>
      </Router>
      
    );
  }
}

export default App;
